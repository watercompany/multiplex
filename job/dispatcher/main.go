package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	// Setup log
	timeNow := time.Now()
	timeNowFormatted := timeNow.Format(time.RFC3339)
	outputName := fmt.Sprintf("%v-dispatcher-log", timeNowFormatted)
	f, err := os.OpenFile("./output/"+outputName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	fmt.Printf("Running dispatcher now...\n")
	log.Printf("Running dispatcher now...\n")

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
			log.Printf("error getting number of queued jobs: %v", err)
			panic(err)
		}

		if jobs != 0 {
			// Deploy a job to worker
			// Get worker port number from
			// avail ports channel.
			currPortNum := <-availPortsCh

			clientCfg, err := job.GetJob()
			if err != nil {
				log.Printf("error getting job: %v", err)
				panic(err)
			}

			go func(clientCfg *client.CallWorkerConfig, currPortNum int) {
				err := job.IncrActiveJobs()
				if err != nil {
					log.Printf("error incrementing active jobs: %v", err)
					panic(err)
				}
				var res *worker.Result
				client.CallWorker(*clientCfg, fmt.Sprintf(":%v", currPortNum), res)

				err = job.DecrActiveJobs()
				if err != nil {
					log.Printf("error decrementing active jobs: %v", err)
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
