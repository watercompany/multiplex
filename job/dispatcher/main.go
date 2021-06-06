package main

import (
	"flag"
	"fmt"
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

// var wg = sync.WaitGroup{}

func init() {
	flag.IntVar(&numberOfWorkers, "workers", 1, "number of available workers")
}

func main() {
	flag.Parse()

	// Clean out temp dirs
	err := CleanOutTempDirs()
	if err != nil {
		panic(err)
	}

	// Clean out final temp dirs
	err = CleanFinalTempdDir()
	if err != nil {
		panic(err)
	}

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
		time.Sleep(5 * time.Second)
		jobs, err := job.GetNumberOfQueuedJobs("")
		if err != redis.Nil && err != nil {
			panic(err)
		}

		if jobs != 0 {
			// Deploy a job to worker
			// Get worker port number from
			// avail ports channel.
			currPortNum := <-availPortsCh

			clientCfg, err := job.GetJob()
			if err != nil {
				panic(err)
			}

			go func(clientCfg *client.CallWorkerConfig, currPortNum int) {

				err := job.IncrActiveJobs()
				if err != nil {
					panic(err)
				}
				var res *worker.Result
				client.CallWorker(*clientCfg, fmt.Sprintf(":%v", currPortNum), res)

				err = job.DecrActiveJobs()
				if err != nil {
					panic(err)
				}
				// Append back the worker port number used so that
				// it can be used by another go routine.
				availPortsCh <- currPortNum
			}(clientCfg, currPortNum)
		}
	}
}

func GetAvailableWorkers(numberOfAvailableWorkers int) []int {
	var workersAddr []int
	for i := 0; i < numberOfAvailableWorkers; i++ {
		workersAddr = append(workersAddr, BasePortNumber+i)
	}
	return workersAddr
}
