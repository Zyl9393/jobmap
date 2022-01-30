## jobmap

Job-processor which allows pending jobs to be superceded by newer jobs based on an arbitrary key associated with each job.

```go
func main() {
	js := jobmap.New()
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()
	go pumpJobs(js, cancelFunc)
	js.Run(ctx)
}

func pumpJobs(js *jobmap.JobMap, cancelFunc context.CancelFunc) {
	defer cancelFunc()
	js.TakeJob("Apple", func(ctx context.Context) { fmt.Println("Apple has been bought."); time.Sleep(time.Second) })
	js.TakeJob("Cheese", func(ctx context.Context) { fmt.Println("Cheese has been bought."); time.Sleep(time.Second) })
	js.TakeJob("Apple", func(ctx context.Context) { fmt.Println("Apple has been peeled.") })
	js.TakeJob("Cheese", func(ctx context.Context) { fmt.Println("Cheese has been cut.") })
	js.TakeJob("Apple", func(ctx context.Context) { fmt.Println("Apple has been eaten.") })
	js.TakeJob("Cheese", func(ctx context.Context) { fmt.Println("Cheese has been eaten.") })
	time.Sleep(time.Second * 2)
}
```

Output:

> Apple has been bought.  
> Cheese has been bought.  
> Apple has been eaten.  
> Cheese has been eaten.  

Two jobs with information about `Apple` and `Cheese` having been eaten do arrive before the first job for either finishes, so the two jobs containing information about their intermediate states (*peeled* and *cut*) are effectively skipped.
