package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/watercompany/multiplex/job"
	"github.com/watercompany/multiplex/worker"
	"github.com/watercompany/multiplex/worker/client"
)

const (
	BasePortNumber = 9090
)

var (
	numberOfWorkers int
)

var wg = sync.WaitGroup{}

func init() {
	flag.IntVar(&numberOfWorkers, "workers", 1, "number of available workers")
}

func main() {
	flag.Parse()
	RunDispatcher()
}

func RunDispatcher() {
	// Make worker ports channel
	var availPorts []int
	availWorkers := GetAvailableWorkers(numberOfWorkers)
	availPorts = append(availPorts, availWorkers...)
	availPortsCh := make(chan int, len(availPorts))
	go func() {
		for _, val := range availPorts {
			availPortsCh <- val
		}
	}()

	for {
		time.Sleep(1 * time.Second)
		// wg.Add(1)
		go func() {
			// defer wg.Done()

			jobs, err := job.GetNumberOfCurrentJobs()
			if err != redis.Nil && err != nil {
				panic(err)
			}

			if jobs != 0 {
				clientCfg, err := job.GetJob()
				if err != nil {
					panic(err)
				}

				// Deploy a job to worker
				// Get worker port number from
				// avail ports channel.
				currPortNum := <-availPortsCh

				var res *worker.Result
				client.CallWorker(*clientCfg, fmt.Sprintf(":%v", currPortNum), res)

				// Append back the worker port number used so that
				// it can be used by another go routine.
				availPortsCh <- currPortNum
			}

		}()
	}
	// wg.Wait()
}

func GetAvailableWorkers(numberOfAvailableWorkers int) []int {
	var workersAddr []int
	for i := 0; i < numberOfAvailableWorkers; i++ {
		workersAddr = append(workersAddr, BasePortNumber+i)
	}
	return workersAddr
}
