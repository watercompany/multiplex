package main

import (
	"flag"
	"fmt"
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

	// Clean out temp dirs
	err := CleanOutTempDirs()
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

func CleanOutTempDirs() error {
	ssd1a := "/mnt/ssd1/plotfiles/temp"
	ssd1b := "/mnt/ssd1/plotfiles/temp2"
	ssd2a := "/mnt/ssd2/plotfiles/temp"
	ssd2b := "/mnt/ssd2/plotfiles/temp2"
	ssd3a := "/mnt/ssd3/plotfiles/temp"
	ssd3b := "/mnt/ssd3/plotfiles/temp2"
	ssd4a := "/mnt/ssd4/plotfiles/temp"
	ssd4b := "/mnt/ssd4/plotfiles/temp2"
	ssd5a := "/mnt/ssd5/plotfiles/temp"
	ssd5b := "/mnt/ssd5/plotfiles/temp2"
	ssd6a := "/mnt/ssd6/plotfiles/temp"
	ssd6b := "/mnt/ssd6/plotfiles/temp2"
	ssd7a := "/mnt/ssd7/plotfiles/temp"
	ssd7b := "/mnt/ssd7/plotfiles/temp2"
	ssd8a := "/mnt/ssd8/plotfiles/temp"
	ssd8b := "/mnt/ssd8/plotfiles/temp2"

	// SSD1
	err := DeleteTempFiles(ssd1a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd1b)
	if err != nil {
		return err
	}

	// SSD2
	err = DeleteTempFiles(ssd2a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd2b)
	if err != nil {
		return err
	}

	// SSD3
	err = DeleteTempFiles(ssd3a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd3b)
	if err != nil {
		return err
	}

	// SSD4
	err = DeleteTempFiles(ssd4a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd4b)
	if err != nil {
		return err
	}

	// SSD5
	err = DeleteTempFiles(ssd5a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd5b)
	if err != nil {
		return err
	}

	// SSD6
	err = DeleteTempFiles(ssd6a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd6b)
	if err != nil {
		return err
	}

	// SSD7
	err = DeleteTempFiles(ssd7a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd7b)
	if err != nil {
		return err
	}

	// SSD8
	err = DeleteTempFiles(ssd8a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(ssd8b)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTempFiles(dir string) error {
	// Delete all temp files
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	// make empty tmp dir
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	return nil
}
