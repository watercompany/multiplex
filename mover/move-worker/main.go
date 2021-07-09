package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	// Clean out temp dirs
	err = CleanOutTempDirs()
	if err != nil {
		log.Printf("error cleaning out temp dirs: %v", err)
		panic(err)
	}
	fmt.Printf("Finished cleaning local temp files.\n")
	log.Printf("Finished cleaning local temp files.\n")

	// // Clean out final temp dirs
	// err = CleanFinalTempdDir()
	// if err != nil {
	// 	log.Printf("error cleaning out final temp dirs: %v", err)
	// 	panic(err)
	// }
	// fmt.Printf("Finished cleaning final destination temp files.\n")
	// log.Printf("Finished cleaning final destination temp files.\n")

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
				if dir.IsDir() || !strings.Contains(fileName, "plot") || strings.Contains(fileName, "tmp") {
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

				// sleep 15 seconds
				// let MoveFile make transfer lock first
				time.Sleep(15 * time.Second)
			}
		}
	}
}

func CleanOutTempDirs() error {
	md0a := "/mnt/md0/temp"
	md0b := "/mnt/md0/temp2"
	md1a := "/mnt/md1/temp"
	md1b := "/mnt/md1/temp2"
	md2a := "/mnt/md2/temp"
	md2b := "/mnt/md2/temp2"
	md3a := "/mnt/md3/temp"
	md3b := "/mnt/md3/temp2"

	// MD0
	err := DeleteTempFiles(md0a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(md0b)
	if err != nil {
		return err
	}

	// MD1
	err = DeleteTempFiles(md1a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(md1b)
	if err != nil {
		return err
	}

	// MD2
	err = DeleteTempFiles(md2a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(md2b)
	if err != nil {
		return err
	}

	// MD3
	err = DeleteTempFiles(md3a)
	if err != nil {
		return err
	}
	err = DeleteTempFiles(md3b)
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
	_, err = cleanerLog.Write([]byte("Deleting files that were last modified 4 hours or older...\n"))
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
		time.Sleep(5 * time.Second)
	}

	for _, name := range names {
		path := filepath.Join(dir, name)
		file, err := os.Stat(path)
		if err != nil {
			return err
		}

		modifiedtime := file.ModTime()
		timeLastModifiedInHours := time.Since(modifiedtime).Hours()
		if timeLastModifiedInHours > 4.0 && strings.Contains(name, "plot") {
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
