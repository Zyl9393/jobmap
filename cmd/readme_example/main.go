package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Zyl9393/jobmap"
)

func main() {
	jm := jobmap.New()
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()
	go pumpJobs(jm, cancelFunc)
	jm.Run(ctx)
}

func pumpJobs(jm *jobmap.JobMap, cancelFunc context.CancelFunc) {
	defer cancelFunc()
	jm.TakeJob("Apple", func(ctx context.Context) { fmt.Println("Apple has been bought."); time.Sleep(time.Second) })
	jm.TakeJob("Cheese", func(ctx context.Context) { fmt.Println("Cheese has been bought."); time.Sleep(time.Second) })
	jm.TakeJob("Apple", func(ctx context.Context) { fmt.Println("Apple has been peeled.") })
	jm.TakeJob("Cheese", func(ctx context.Context) { fmt.Println("Cheese has been cut.") })
	jm.TakeJob("Apple", func(ctx context.Context) { fmt.Println("Apple has been eaten.") })
	jm.TakeJob("Cheese", func(ctx context.Context) { fmt.Println("Cheese has been eaten.") })
	time.Sleep(time.Second * 2)
}
