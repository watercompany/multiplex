package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// Clean out temp dirs
	err = CleanOutTempDirs()
	if err != nil {
		log.Printf("error cleaning out temp dirs: %v", err)
		panic(err)
	}
	fmt.Printf("Finished cleaning local temp files.\n")

	// Clean out final temp dirs
	err = CleanFinalTempdDir()
	if err != nil {
		log.Printf("error cleaning out final temp dirs: %v", err)
		panic(err)
	}
	fmt.Printf("Finished cleaning final destination temp files.\n")
	fmt.Printf("Running dispatcher now...\n")

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

func CleanFinalTempdDir() error {
	nas1 := "/mnt/skynas-1/tmp"
	nas2 := "/mnt/skynas-2/tmp"
	nas3 := "/mnt/skynas-3/tmp"
	nas4 := "/mnt/skynas-4/tmp"
	nas5 := "/mnt/skynas-5/tmp"
	nas6 := "/mnt/skynas-6/tmp"
	nas7 := "/mnt/skynas-7/tmp"
	nas8 := "/mnt/skynas-8/tmp"
	nas9 := "/mnt/skynas-9/tmp"
	nas10 := "/mnt/skynas-10/tmp"
	nas11 := "/mnt/skynas-11/tmp"
	nas12 := "/mnt/skynas-12/tmp"
	nas13 := "/mnt/skynas-13/tmp"
	nas14 := "/mnt/skynas-14/tmp"
	nas15 := "/mnt/skynas-15/tmp"
	nas16 := "/mnt/skynas-16/tmp"

	// NAS1
	err := RemoveStagnantTempFiles(nas1)
	if err != nil {
		return err
	}

	// NAS2
	err = RemoveStagnantTempFiles(nas2)
	if err != nil {
		return err
	}

	// NAS3
	err = RemoveStagnantTempFiles(nas3)
	if err != nil {
		return err
	}

	// NAS4
	err = RemoveStagnantTempFiles(nas4)
	if err != nil {
		return err
	}

	// NAS5
	err = RemoveStagnantTempFiles(nas5)
	if err != nil {
		return err
	}

	// NAS6
	err = RemoveStagnantTempFiles(nas6)
	if err != nil {
		return err
	}

	// NAS7
	err = RemoveStagnantTempFiles(nas7)
	if err != nil {
		return err
	}

	// NAS8
	err = RemoveStagnantTempFiles(nas8)
	if err != nil {
		return err
	}

	// NAS9
	err = RemoveStagnantTempFiles(nas9)
	if err != nil {
		return err
	}

	// NAS10
	err = RemoveStagnantTempFiles(nas10)
	if err != nil {
		return err
	}

	// NAS11
	err = RemoveStagnantTempFiles(nas11)
	if err != nil {
		return err
	}

	// NAS12
	err = RemoveStagnantTempFiles(nas12)
	if err != nil {
		return err
	}

	// NAS13
	err = RemoveStagnantTempFiles(nas13)
	if err != nil {
		return err
	}

	// NAS14
	err = RemoveStagnantTempFiles(nas14)
	if err != nil {
		return err
	}

	// NAS15
	err = RemoveStagnantTempFiles(nas15)
	if err != nil {
		return err
	}

	// NAS16
	err = RemoveStagnantTempFiles(nas16)
	if err != nil {
		return err
	}

	return nil
}

func RemoveStagnantTempFiles(dir string) error {
	// make empty tmp dir
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	// make log file
	timeNow := time.Now()
	timeNowFormatted := timeNow.Format(time.RFC3339)
	cleanerLogName := fmt.Sprintf("/%v-cleaner-log", timeNowFormatted)
	cleanerLog, err := os.OpenFile(dir+cleanerLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer cleanerLog.Close()

	// Log header
	_, err = cleanerLog.Write([]byte("Deleting files that were last modified 2 hours or older...\n"))
	if err != nil {
		return err
	}

	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	var names []string

	// TODO: find better fix for readdirent error
	// with using Readdirnames on cifs mounts
	readDirPass := true
	for readDirPass {
		names, err = d.Readdirnames(-1)
		if err == nil {
			readDirPass = false
			continue
		}

		if !strings.Contains(err.Error(), "readdirent") {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	for _, name := range names {
		path := filepath.Join(dir, name)
		file, err := os.Stat(path)
		if err != nil {
			return err
		}

		modifiedtime := file.ModTime()
		timeLastModifiedInHours := time.Since(modifiedtime).Hours()
		if timeLastModifiedInHours > 2.0 {
			err = os.RemoveAll(path)
			if err != nil {
				return err
			}

			// Log deleted file
			_, err = cleanerLog.Write([]byte(fmt.Sprintf("Deleted File: %v\nLast Modified: %v hours\n", name, timeLastModifiedInHours)))
			if err != nil {
				return err
			}
			_, err = cleanerLog.Write([]byte(fmt.Sprintf("Deleted File Info: %+v\n\n", file)))
			if err != nil {
				return err
			}

		}
	}
	return nil
}
