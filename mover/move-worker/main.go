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

	fmt.Printf("Running mover now...\n")
	log.Printf("Running mover now...\n")

	err = RunMover()
	if err != nil {
		log.Printf("error mover: %v", err)
		panic(err)
	}
}

func RunMover() error {
	localFinalDirs := []string{
		// "/mnt/ssd1/plotfiles/final",
		// "/mnt/ssd2/plotfiles/final",
		// "/mnt/ssd3/plotfiles/final",
		// "/mnt/ssd4/plotfiles/final",
		// "/mnt/ssd5/plotfiles/final",
		// "/mnt/ssd6/plotfiles/final",
		// "/mnt/ssd7/plotfiles/final",
		// "/mnt/ssd8/plotfiles/final",
		"/mnt/md0/final",
		"/mnt/md1/final",
		"/mnt/md2/final",
		"/mnt/md3/final",
	}

	finalDirs := []string{
		// "/mnt/skynas-1",
		// "/mnt/skynas-2",
		// "/mnt/skynas-3",
		// "/mnt/skynas-4",
		// "/mnt/skynas-5",
		// "/mnt/skynas-6",
		// "/mnt/skynas-7",
		// "/mnt/skynas-8",
		// "/mnt/skynas-9",
		// "/mnt/skynas-10",
		// "/mnt/skynas-11",
		// "/mnt/skynas-12",
		// "/mnt/skynas-13",
		// "/mnt/skynas-14",
		// "/mnt/skynas-15",
		// "/mnt/skynas-16",
		"/mnt/ct-01/hdd1",
		"/mnt/ct-01/hdd2",
		"/mnt/ct-01/hdd3",
		"/mnt/ct-01/hdd4",
		"/mnt/ct-01/hdd5",
		"/mnt/ct-01/hdd6",
		"/mnt/ct-01/hdd7",
		"/mnt/ct-01/hdd8",
		"/mnt/ct-01/hdd9",
		"/mnt/ct-01/hdd10",
		"/mnt/ct-01/hdd11",
		"/mnt/ct-01/hdd12",
		"/mnt/ct-01/hdd13",
		"/mnt/ct-01/hdd14",
		"/mnt/ct-01/hdd15",
		"/mnt/ct-01/hdd16",
		"/mnt/ct-01/hdd17",
		"/mnt/ct-01/hdd18",
		"/mnt/ct-01/hdd19",
		"/mnt/ct-01/hdd20",
		"/mnt/ct-01/hdd21",
		"/mnt/ct-01/hdd22",
		"/mnt/ct-01/hdd23",
		"/mnt/ct-01/hdd24",
		"/mnt/ct-02/hdd1",
		"/mnt/ct-02/hdd2",
		"/mnt/ct-02/hdd3",
		"/mnt/ct-02/hdd4",
		"/mnt/ct-02/hdd5",
		"/mnt/ct-02/hdd6",
		"/mnt/ct-02/hdd7",
		"/mnt/ct-02/hdd8",
		"/mnt/ct-02/hdd9",
		"/mnt/ct-02/hdd10",
		"/mnt/ct-02/hdd11",
		"/mnt/ct-02/hdd12",
		"/mnt/ct-02/hdd13",
		"/mnt/ct-02/hdd14",
		"/mnt/ct-02/hdd15",
		"/mnt/ct-02/hdd16",
		"/mnt/ct-02/hdd17",
		"/mnt/ct-02/hdd18",
		"/mnt/ct-02/hdd19",
		"/mnt/ct-02/hdd20",
		"/mnt/ct-02/hdd21",
		"/mnt/ct-02/hdd22",
		"/mnt/ct-02/hdd23",
		"/mnt/ct-02/hdd24",
		"/mnt/ct-03/hdd1",
		"/mnt/ct-03/hdd2",
		"/mnt/ct-03/hdd3",
		"/mnt/ct-03/hdd4",
		"/mnt/ct-03/hdd5",
		"/mnt/ct-03/hdd6",
		"/mnt/ct-03/hdd7",
		"/mnt/ct-03/hdd8",
		"/mnt/ct-03/hdd9",
		"/mnt/ct-03/hdd10",
		"/mnt/ct-03/hdd11",
		"/mnt/ct-03/hdd12",
		"/mnt/ct-03/hdd13",
		"/mnt/ct-03/hdd14",
		"/mnt/ct-03/hdd15",
		"/mnt/ct-03/hdd16",
		"/mnt/ct-03/hdd17",
		"/mnt/ct-03/hdd18",
		"/mnt/ct-03/hdd19",
		"/mnt/ct-03/hdd20",
		"/mnt/ct-03/hdd21",
		"/mnt/ct-03/hdd22",
		"/mnt/ct-03/hdd23",
		"/mnt/ct-03/hdd24",
		"/mnt/ct-04/hdd1",
		"/mnt/ct-04/hdd2",
		"/mnt/ct-04/hdd3",
		"/mnt/ct-04/hdd4",
		"/mnt/ct-04/hdd5",
		"/mnt/ct-04/hdd6",
		"/mnt/ct-04/hdd7",
		"/mnt/ct-04/hdd8",
		"/mnt/ct-04/hdd9",
		"/mnt/ct-04/hdd10",
		"/mnt/ct-04/hdd11",
		"/mnt/ct-04/hdd12",
		"/mnt/ct-04/hdd13",
		"/mnt/ct-04/hdd14",
		"/mnt/ct-04/hdd15",
		"/mnt/ct-04/hdd16",
		"/mnt/ct-04/hdd17",
		"/mnt/ct-04/hdd18",
		"/mnt/ct-04/hdd19",
		"/mnt/ct-04/hdd20",
		"/mnt/ct-04/hdd21",
		"/mnt/ct-04/hdd22",
		"/mnt/ct-04/hdd23",
		"/mnt/ct-04/hdd24",
		"/mnt/ct-05/hdd1",
		"/mnt/ct-05/hdd2",
		"/mnt/ct-05/hdd3",
		"/mnt/ct-05/hdd4",
		"/mnt/ct-05/hdd5",
		"/mnt/ct-05/hdd6",
		"/mnt/ct-05/hdd7",
		"/mnt/ct-05/hdd8",
		"/mnt/ct-05/hdd9",
		"/mnt/ct-05/hdd10",
		"/mnt/ct-05/hdd11",
		"/mnt/ct-05/hdd12",
		"/mnt/ct-05/hdd13",
		"/mnt/ct-05/hdd14",
		"/mnt/ct-05/hdd15",
		"/mnt/ct-05/hdd16",
		"/mnt/ct-05/hdd17",
		"/mnt/ct-05/hdd18",
		"/mnt/ct-05/hdd19",
		"/mnt/ct-05/hdd20",
		"/mnt/ct-05/hdd21",
		"/mnt/ct-05/hdd22",
		"/mnt/ct-05/hdd23",
		"/mnt/ct-05/hdd24",
		"/mnt/ct-06/hdd1",
		"/mnt/ct-06/hdd2",
		"/mnt/ct-06/hdd3",
		"/mnt/ct-06/hdd4",
		"/mnt/ct-06/hdd5",
		"/mnt/ct-06/hdd6",
		"/mnt/ct-06/hdd7",
		"/mnt/ct-06/hdd8",
		"/mnt/ct-06/hdd9",
		"/mnt/ct-06/hdd10",
		"/mnt/ct-06/hdd11",
		"/mnt/ct-06/hdd12",
		"/mnt/ct-06/hdd13",
		"/mnt/ct-06/hdd14",
		"/mnt/ct-06/hdd15",
		"/mnt/ct-06/hdd16",
		"/mnt/ct-06/hdd17",
		"/mnt/ct-06/hdd18",
		"/mnt/ct-06/hdd19",
		"/mnt/ct-06/hdd20",
		"/mnt/ct-06/hdd21",
		"/mnt/ct-06/hdd22",
		"/mnt/ct-06/hdd23",
		"/mnt/ct-06/hdd24",
		"/mnt/ct-07/hdd1",
		"/mnt/ct-07/hdd2",
		"/mnt/ct-07/hdd3",
		"/mnt/ct-07/hdd4",
		"/mnt/ct-07/hdd5",
		"/mnt/ct-07/hdd6",
		"/mnt/ct-07/hdd7",
		"/mnt/ct-07/hdd8",
		"/mnt/ct-07/hdd9",
		"/mnt/ct-07/hdd10",
		"/mnt/ct-07/hdd11",
		"/mnt/ct-07/hdd12",
		"/mnt/ct-07/hdd13",
		"/mnt/ct-07/hdd14",
		"/mnt/ct-07/hdd15",
		"/mnt/ct-07/hdd16",
		"/mnt/ct-07/hdd17",
		"/mnt/ct-07/hdd18",
		"/mnt/ct-07/hdd19",
		"/mnt/ct-07/hdd20",
		"/mnt/ct-07/hdd21",
		"/mnt/ct-07/hdd22",
		"/mnt/ct-07/hdd23",
		"/mnt/ct-07/hdd24",
		"/mnt/ct-08/hdd1",
		"/mnt/ct-08/hdd2",
		"/mnt/ct-08/hdd3",
		"/mnt/ct-08/hdd4",
		"/mnt/ct-08/hdd5",
		"/mnt/ct-08/hdd6",
		"/mnt/ct-08/hdd7",
		"/mnt/ct-08/hdd8",
		"/mnt/ct-08/hdd9",
		"/mnt/ct-08/hdd10",
		"/mnt/ct-08/hdd11",
		"/mnt/ct-08/hdd12",
		"/mnt/ct-08/hdd13",
		"/mnt/ct-08/hdd14",
		"/mnt/ct-08/hdd15",
		"/mnt/ct-08/hdd16",
		"/mnt/ct-08/hdd17",
		"/mnt/ct-08/hdd18",
		"/mnt/ct-08/hdd19",
		"/mnt/ct-08/hdd20",
		"/mnt/ct-08/hdd21",
		"/mnt/ct-08/hdd22",
		"/mnt/ct-08/hdd23",
		"/mnt/ct-08/hdd24",
		"/mnt/ct-09/hdd1",
		"/mnt/ct-09/hdd2",
		"/mnt/ct-09/hdd3",
		"/mnt/ct-09/hdd4",
		"/mnt/ct-09/hdd5",
		"/mnt/ct-09/hdd6",
		"/mnt/ct-09/hdd7",
		"/mnt/ct-09/hdd8",
		"/mnt/ct-09/hdd9",
		"/mnt/ct-09/hdd10",
		"/mnt/ct-09/hdd11",
		"/mnt/ct-09/hdd12",
		"/mnt/ct-09/hdd13",
		"/mnt/ct-09/hdd14",
		"/mnt/ct-09/hdd15",
		"/mnt/ct-09/hdd16",
		"/mnt/ct-09/hdd17",
		"/mnt/ct-09/hdd18",
		"/mnt/ct-09/hdd19",
		"/mnt/ct-09/hdd20",
		"/mnt/ct-09/hdd21",
		"/mnt/ct-09/hdd22",
		"/mnt/ct-09/hdd23",
		"/mnt/ct-09/hdd24",
		"/mnt/ct-10/hdd1",
		"/mnt/ct-10/hdd2",
		"/mnt/ct-10/hdd3",
		"/mnt/ct-10/hdd4",
		"/mnt/ct-10/hdd5",
		"/mnt/ct-10/hdd6",
		"/mnt/ct-10/hdd7",
		"/mnt/ct-10/hdd8",
		"/mnt/ct-10/hdd9",
		"/mnt/ct-10/hdd10",
		"/mnt/ct-10/hdd11",
		"/mnt/ct-10/hdd12",
		"/mnt/ct-10/hdd13",
		"/mnt/ct-10/hdd14",
		"/mnt/ct-10/hdd15",
		"/mnt/ct-10/hdd16",
		"/mnt/ct-10/hdd17",
		"/mnt/ct-10/hdd18",
		"/mnt/ct-10/hdd19",
		"/mnt/ct-10/hdd20",
		"/mnt/ct-10/hdd21",
		"/mnt/ct-10/hdd22",
		"/mnt/ct-10/hdd23",
		"/mnt/ct-10/hdd24",
	}

	maxLockFiles := 1
	maxAgeOfLockFilesInHours := 2.0
	movingFiles := 0
	maxMovingFiles := 1
	for {
		time.Sleep(5 * time.Second)

		// Delete old transfer locks
		for _, finalDir := range finalDirs {
			err := mover.DeleteIfNoActivity(finalDir, mover.TransferLockName, maxAgeOfLockFilesInHours)
			if err != nil {
				// Log only, dont return error
				log.Printf("error in cleaning transfer locks: %v", err)
			}
		}

		for _, localFinalDir := range localFinalDirs {
			dirs, err := mover.GetDirs(localFinalDir)
			if err != nil {
				return err
			}

			for _, dir := range dirs {
				fileName := dir.Name()
				if dir.IsDir() || !strings.Contains(fileName, "plot") || strings.Contains(fileName, "tmp") || movingFiles > maxMovingFiles {
					continue
				}

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

				movingFiles++
				go func(localFinalDir, finalDir, fileName string) {
					log.Printf("Moving file %v from %v to %v", fileName, localFinalDir, finalDir)
					// Moves and Deletes
					err = mover.MoveFile(localFinalDir, finalDir, fileName)
					if err != nil {
						movingFiles--
						log.Printf("err moving file %v from %v to %v: %v", fileName, localFinalDir, finalDir, err)
						return
					}
					movingFiles--
					log.Printf("Finished moving file %v from %v to %v\n", fileName, localFinalDir, finalDir)
					log.Printf("Moving Files: %v\n", movingFiles)
				}(localFinalDir, finalDir, fileName)

				// sleep 15 seconds
				// let MoveFile make transfer lock first
				time.Sleep(15 * time.Second)
			}
		}
	}
}
