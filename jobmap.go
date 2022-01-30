package jobmap

import (
	"context"
	"sync"
)

// JobFunc is a function which represents a job.
type JobFunc func(context.Context)

type jobDescription struct {
	key interface{}
	job JobFunc
}

// JobMap associates function invocations ("jobs") with keys (interface{}),
// where no two jobs with the same key may run at the same time, and only
// the latest known job for a given key may be pending to run.
type JobMap struct {
	jobRequests    chan jobDescription
	jobCompletions chan interface{}
	runningJobs    map[interface{}]JobFunc
	pendingJobs    map[interface{}]JobFunc
	wgRunningJobs  *sync.WaitGroup
}

// New constructs a new JobMap and returns a pointer to it.
func New() *JobMap {
	return &JobMap{
		jobRequests:    make(chan jobDescription, 64),
		jobCompletions: make(chan interface{}, 64),
		runningJobs:    make(map[interface{}]JobFunc),
		pendingJobs:    make(map[interface{}]JobFunc),
		wgRunningJobs:  &sync.WaitGroup{},
	}
}

// TakeJob thread-safely queues a job to run on a new go-routine for given key.
func (jm *JobMap) TakeJob(key interface{}, job JobFunc) {
	defer func() {
		recover()
	}()
	jm.jobRequests <- jobDescription{key: key, job: job}
}

// Run starts jobs queued with TakeJob() on their own go-routines.
// If no job for a given key is running, the next job with that key will run immediately.
// If a job for a given key is running, a following job with that key will be pending to run after the running job finishes.
// If a job for a given key is pending to run, a following job with the same key will take its place.
// When ctx is cancelled, Run will wait for running jobs to finish and return.
func (jm *JobMap) Run(ctx context.Context) {
	for {
		select {
		case request := <-jm.jobRequests:
			if _, ok := jm.runningJobs[request.key]; !ok {
				jm.runJob(ctx, request.key, request.job)
			} else {
				jm.pendingJobs[request.key] = request.job
			}
		case key := <-jm.jobCompletions:
			delete(jm.runningJobs, key)
			job, ok := jm.pendingJobs[key]
			if ok {
				delete(jm.pendingJobs, key)
				jm.runJob(ctx, key, job)
			}
		case <-ctx.Done():
			close(jm.jobRequests)
			close(jm.jobCompletions)
			jm.wgRunningJobs.Wait()
			return
		}
	}
}

func (jm *JobMap) runJob(ctx context.Context, key interface{}, job JobFunc) {
	if ctx.Err() != nil {
		return
	}
	jm.wgRunningJobs.Add(1)
	jm.runningJobs[key] = job
	go func() {
		job(ctx)
		jm.wgRunningJobs.Done()
		defer func() {
			recover()
		}()
		jm.jobCompletions <- key
	}()
}
