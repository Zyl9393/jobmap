package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Zyl9393/jobmap"
)

var m *sync.Map
var x int32

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()
	m = &sync.Map{}
	jm := jobmap.New()
	go pumpJobs(jm, cancelFunc)
	jm.Run(ctx)
}

func pumpJobs(jm *jobmap.JobMap, cancelFunc context.CancelFunc) {
	defer cancelFunc()
	jm.TakeJob("berry", func(ctx context.Context) { intSetSlow(time.Second, "berry", 1) })
	jm.TakeJob("cheese", func(ctx context.Context) { intSetSlow(time.Second, "cheese", 1) })
	time.Sleep(time.Millisecond * 500)
	jm.TakeJob("berry", func(ctx context.Context) { intSetSlow(time.Second, "berry", 2) })
	jm.TakeJob("berry", func(ctx context.Context) { intSetSlow(time.Second, "berry", 3) })
	jm.TakeJob("berry", func(ctx context.Context) { intSetSlow(time.Second, "berry", 4) })
	jm.TakeJob("berry", func(ctx context.Context) { intSetSlow(time.Second, "berry", 5) })
	jm.TakeJob("cheese", func(ctx context.Context) { intSetSlow(time.Second, "cheese", 2) })
	jm.TakeJob("cheese", func(ctx context.Context) { intSetSlow(time.Second, "cheese", 3) })
	jm.TakeJob("cheese", func(ctx context.Context) { intSetSlow(time.Second, "cheese", 4) })
	jm.TakeJob("cheese", func(ctx context.Context) { intSetSlow(time.Second, "cheese", 5) })
	time.Sleep(time.Millisecond * 1000)
	if v, ok := m.Load("berry"); !ok || v.(int) != 1 {
		log.Fatalf("berry was %d when expected 1.", v.(int))
	}
	if v, ok := m.Load("cheese"); !ok || v.(int) != 1 {
		log.Fatalf("cheese was %d when expected 1.", v.(int))
	}
	time.Sleep(time.Millisecond * 1000)
	if v, ok := m.Load("berry"); !ok || v.(int) != 5 {
		log.Fatalf("berry was %d when expected 5.", v.(int))
	}
	if v, ok := m.Load("cheese"); !ok || v.(int) != 5 {
		log.Fatalf("cheese was %d when expected 5.", v.(int))
	}
	time.Sleep(time.Millisecond * 1000)
	xNow := atomic.LoadInt32(&x)
	if xNow != 4 {
		log.Fatalf("x was %d when expected 4.", x)
	}
}

func intSetSlow(waitTime time.Duration, intName string, value int) {
	time.Sleep(time.Second)
	m.Store(intName, value)
	atomic.AddInt32(&x, 1)
}
