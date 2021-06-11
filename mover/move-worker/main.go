package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/watercompany/multiplex/mover"
)

func main() {
	// Setup log
	timeNow := time.Now()
	timeNowFormatted := timeNow.Format(time.RFC3339)
	outputName := fmt.Sprintf("%v-mover-log", timeNowFormatted)
	f, err := os.OpenFile("./output/"+outputName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	err = RunMover()
	if err != nil {
		log.Printf("error mover: %v", err)
		panic(err)
	}
}

func RunMover() error {
	localFinalDirs := []string{
		"/mnt/ssd1/plotfiles/final",
		"/mnt/ssd2/plotfiles/final",
		"/mnt/ssd3/plotfiles/final",
		"/mnt/ssd4/plotfiles/final",
		"/mnt/ssd5/plotfiles/final",
		"/mnt/ssd6/plotfiles/final",
		"/mnt/ssd7/plotfiles/final",
		"/mnt/ssd8/plotfiles/final",
	}

	finalDirs := []string{
		"/mnt/skynas-1",
		"/mnt/skynas-2",
		"/mnt/skynas-3",
		"/mnt/skynas-4",
		"/mnt/skynas-5",
		"/mnt/skynas-6",
		"/mnt/skynas-7",
		"/mnt/skynas-8",
		"/mnt/skynas-9",
		"/mnt/skynas-10",
		"/mnt/skynas-11",
		"/mnt/skynas-12",
		"/mnt/skynas-13",
		"/mnt/skynas-14",
		"/mnt/skynas-15",
		"/mnt/skynas-16",
	}

	maxLockFiles := 8

	for {
		time.Sleep(5 * time.Second)

		for _, localFinalDir := range localFinalDirs {
			dirs, err := mover.GetDirs(localFinalDir)
			if err != nil {
				return err
			}

			for _, dir := range dirs {
				if dir.IsDir() || !strings.Contains(dir.Name(), "plot") {
					continue
				}

				fileName := dir.Name()

				// Get src file size
				fileSize, err := mover.FileSizeInMB(localFinalDir + "/" + fileName)
				if err != nil {
					return err
				}

				// Choose random destination that
				// does not exceed max parallel transfer
				finalDir, err := mover.GetFinalDir(fileSize, finalDirs, maxLockFiles)
				if err != nil {
					return err
				}

				if finalDir == "" {
					continue
				}

				go func(localFinalDir, finalDir, fileName string) {
					log.Printf("Moving file %v from %v to %v", fileName, localFinalDir, finalDir)
					// Moves and Deletes
					err = mover.MoveFile(localFinalDir, finalDir, fileName)
					if err != nil {
						log.Printf("err moving file %v from %v to %v: %v", fileName, localFinalDir, finalDir, err)
						return
					}
					log.Printf("Finished moving file %v from %v to %v", fileName, localFinalDir, finalDir)
				}(localFinalDir, finalDir, fileName)

				// sleep 5 seconds
				// let MoveFile make transfer lock first
				time.Sleep(5 * time.Second)
			}
		}
	}
}
