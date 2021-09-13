package mover

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

func FileCount(path string) (int, error) {
	var dirs []fs.DirEntry
	var err error
	// TODO: find better fix for readdirent error
	// with using Readdirnames on cifs mounts
	readDirPass := true
	for readDirPass {
		dirs, err = os.ReadDir(path)
		if err == nil {
			readDirPass = false
			continue
		}

		if !strings.Contains(err.Error(), "readdirent") {
			return 0, fmt.Errorf("error counting files: %v", err)
		}
		time.Sleep(5 * time.Second)
	}

	return len(dirs), nil
}

func FileCountSubString(path string, subStr string) (int, error) {
	count := 0
	var dirs []fs.DirEntry
	var err error
	// TODO: find better fix for readdirent error
	// with using Readdirnames on cifs mounts
	readDirPass := true
	for readDirPass {
		dirs, err = os.ReadDir(path)
		if err == nil {
			readDirPass = false
			continue
		}

		if !strings.Contains(err.Error(), "readdirent") {
			return 0, fmt.Errorf("error counting files: %v", err)
		}
		time.Sleep(5 * time.Second)
	}

	for _, dir := range dirs {
		if strings.Contains(dir.Name(), subStr) {
			count++
		}
	}

	return count, nil
}

func GetDirs(path string) ([]fs.DirEntry, error) {
	var dirs []fs.DirEntry
	var err error
	// TODO: find better fix for readdirent error
	// with using Readdirnames on cifs mounts
	readDirPass := true
	for readDirPass {
		dirs, err = os.ReadDir(path)
		if err == nil {
			readDirPass = false
			continue
		}

		if !strings.Contains(err.Error(), "readdirent") {
			return []fs.DirEntry{}, fmt.Errorf("error counting files: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
	return dirs, nil
}

func GetFreeDiskSpaceInMB(path string) uint64 {
	var stat unix.Statfs_t
	unix.Statfs(path, &stat)

	// Available blocks * size per block = available space in bytes
	// divided by 1e6 = megabytes
	return (stat.Bavail * uint64(stat.Bsize)) / 1e6
}

func DirSizeInMB(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size / 1e6, err
}

func FileSizeInMB(path string) (int64, error) {
	var size int64
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	size = info.Size()
	return size / 1e6, err
}

// GetFinalDir returns final dir that
// can accomodate the size of the file.
func GetFinalDir(fileSize int64, finalDirs []string, maxLockFiles int) (string, error) {
	// randomize finalDirs
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(finalDirs), func(i, j int) { finalDirs[i], finalDirs[j] = finalDirs[j], finalDirs[i] })

	// no priority
	for _, finalDir := range finalDirs {
		// Check if directory has row-1, row-2, or row-3
		row1Count, err := FileCountSubString(finalDir, "row-1")
		if err != nil {
			return "", err
		}

		row2Count, err := FileCountSubString(finalDir, "row-2")
		if err != nil {
			return "", err
		}

		row3Count, err := FileCountSubString(finalDir, "row-3")
		if err != nil {
			return "", err
		}

		if row1Count == 0 && row2Count == 0 && row3Count == 0 {
			continue
		}

		// Check if exceeded max lock files
		count, err := FileCountSubString(finalDir, TransferLockName)
		if err != nil {
			return "", err
		}
		hasExceededMaxLocks := count >= maxLockFiles

		// Check if destpath have available size for that
		destFreeSpace := GetFreeDiskSpaceInMB(finalDir)
		// has available space if filesize+30GB+(current transfer count*102GB) is less than free space
		hasAvailableSpace := (fileSize + 30000 + (int64(count) * 102000)) < int64(destFreeSpace)

		if hasAvailableSpace && !hasExceededMaxLocks {
			return finalDir, nil
		}
	}

	// // row-1 group priority
	// for _, finalDir := range finalDirs {
	// 	// Check if directory is row-1
	// 	row1Count, err := FileCountSubString(finalDir, "row-1")
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	if row1Count == 0 {
	// 		continue
	// 	}

	// 	// Check if exceeded max lock files
	// 	count, err := FileCountSubString(finalDir, TransferLockName)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	hasExceededMaxLocks := count >= maxLockFiles

	// 	// Check if destpath have available size for that
	// 	destFreeSpace := GetFreeDiskSpaceInMB(finalDir)
	// 	// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
	// 	hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

	// 	if hasAvailableSpace && !hasExceededMaxLocks {
	// 		return finalDir, nil
	// 	}
	// }

	// // row-2 group priority
	// for _, finalDir := range finalDirs {
	// 	// Check if directory is row-2
	// 	row2Count, err := FileCountSubString(finalDir, "row-2")
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	if row2Count == 0 {
	// 		continue
	// 	}

	// 	// Check if exceeded max lock files
	// 	count, err := FileCountSubString(finalDir, TransferLockName)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	hasExceededMaxLocks := count >= maxLockFiles

	// 	// Check if destpath have available size for that
	// 	destFreeSpace := GetFreeDiskSpaceInMB(finalDir)
	// 	// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
	// 	hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

	// 	if hasAvailableSpace && !hasExceededMaxLocks {
	// 		return finalDir, nil
	// 	}
	// }

	// // row-3 group priority
	// for _, finalDir := range finalDirs {
	// 	// Check if directory is row-3
	// 	row3Count, err := FileCountSubString(finalDir, "row-3")
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	if row3Count == 0 {
	// 		continue
	// 	}

	// 	// Check if exceeded max lock files
	// 	count, err := FileCountSubString(finalDir, TransferLockName)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	hasExceededMaxLocks := count >= maxLockFiles

	// 	// Check if destpath have available size for that
	// 	destFreeSpace := GetFreeDiskSpaceInMB(finalDir)
	// 	// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
	// 	hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

	// 	if hasAvailableSpace && !hasExceededMaxLocks {
	// 		return finalDir, nil
	// 	}
	// }

	// no row group priority
	// for _, finalDir := range finalDirs {
	// 	// Check if exceeded max lock files
	// 	count, err := FileCountSubString(finalDir, TransferLockName)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	hasExceededMaxLocks := count >= maxLockFiles

	// 	// Check if destpath have available size for that
	// 	destFreeSpace := GetFreeDiskSpaceInMB(finalDir)
	// 	// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
	// 	hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

	// 	if hasAvailableSpace && !hasExceededMaxLocks {
	// 		return finalDir, nil
	// 	}
	// }

	// No available dirs
	return "", nil
}

// Delete files with substring subStr that
// have no activity for more than AgeLimitInHours.
func DeleteIfNoActivity(path string, subStr string, AgeLimitInHours float64) error {
	var dirs []fs.DirEntry
	var err error
	// TODO: find better fix for readdirent error
	// with using Readdirnames on cifs mounts
	readDirPass := true
	for readDirPass {
		dirs, err = os.ReadDir(path)
		if err == nil {
			readDirPass = false
			continue
		}

		if !strings.Contains(err.Error(), "readdirent") {
			return fmt.Errorf("error counting files: %v", err)
		}
		time.Sleep(5 * time.Second)
	}

	for _, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			return fmt.Errorf("error getting file info: %v", err)
		}

		if strings.Contains(dir.Name(), subStr) &&
			time.Since(info.ModTime()).Hours() > AgeLimitInHours {
			// Delete file
			err = os.Remove(path + "/" + dir.Name())
			if err != nil {
				return fmt.Errorf("failed deleting file: %s", err)
			}
		}
	}

	return nil
}

func GetRsyncDestPath(destPath string) (string, error) {
	rsyncDestPath := ""
	if !strings.HasSuffix(destPath, "/") {
		destPath = destPath + "/"
	}

	switch destPath {
	case "/mnt/ct-01/hdd1/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd1/tmp"
	case "/mnt/ct-01/hdd2/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd2/tmp"
	case "/mnt/ct-01/hdd3/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd3/tmp"
	case "/mnt/ct-01/hdd4/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd4/tmp"
	case "/mnt/ct-01/hdd5/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd5/tmp"
	case "/mnt/ct-01/hdd6/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd6/tmp"
	case "/mnt/ct-01/hdd7/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd7/tmp"
	case "/mnt/ct-01/hdd8/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd8/tmp"
	case "/mnt/ct-01/hdd9/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd9/tmp"
	case "/mnt/ct-01/hdd10/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd10/tmp"
	case "/mnt/ct-01/hdd11/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd11/tmp"
	case "/mnt/ct-01/hdd12/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd12/tmp"
	case "/mnt/ct-01/hdd13/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd13/tmp"
	case "/mnt/ct-01/hdd14/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd14/tmp"
	case "/mnt/ct-01/hdd15/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd15/tmp"
	case "/mnt/ct-01/hdd16/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd16/tmp"
	case "/mnt/ct-01/hdd17/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd17/tmp"
	case "/mnt/ct-01/hdd18/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd18/tmp"
	case "/mnt/ct-01/hdd19/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd19/tmp"
	case "/mnt/ct-01/hdd20/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd20/tmp"
	case "/mnt/ct-01/hdd21/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd21/tmp"
	case "/mnt/ct-01/hdd22/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd22/tmp"
	case "/mnt/ct-01/hdd23/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd23/tmp"
	case "/mnt/ct-01/hdd24/":
		rsyncDestPath = "ct@192.168.4.1:/mnt/hdd24/tmp"
	case "/mnt/ct-02/hdd1/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd1/tmp"
	case "/mnt/ct-02/hdd2/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd2/tmp"
	case "/mnt/ct-02/hdd3/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd3/tmp"
	case "/mnt/ct-02/hdd4/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd4/tmp"
	case "/mnt/ct-02/hdd5/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd5/tmp"
	case "/mnt/ct-02/hdd6/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd6/tmp"
	case "/mnt/ct-02/hdd7/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd7/tmp"
	case "/mnt/ct-02/hdd8/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd8/tmp"
	case "/mnt/ct-02/hdd9/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd9/tmp"
	case "/mnt/ct-02/hdd10/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd10/tmp"
	case "/mnt/ct-02/hdd11/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd11/tmp"
	case "/mnt/ct-02/hdd12/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd12/tmp"
	case "/mnt/ct-02/hdd13/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd13/tmp"
	case "/mnt/ct-02/hdd14/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd14/tmp"
	case "/mnt/ct-02/hdd15/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd15/tmp"
	case "/mnt/ct-02/hdd16/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd16/tmp"
	case "/mnt/ct-02/hdd17/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd17/tmp"
	case "/mnt/ct-02/hdd18/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd18/tmp"
	case "/mnt/ct-02/hdd19/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd19/tmp"
	case "/mnt/ct-02/hdd20/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd20/tmp"
	case "/mnt/ct-02/hdd21/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd21/tmp"
	case "/mnt/ct-02/hdd22/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd22/tmp"
	case "/mnt/ct-02/hdd23/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd23/tmp"
	case "/mnt/ct-02/hdd24/":
		rsyncDestPath = "ct@192.168.4.2:/mnt/hdd24/tmp"
	case "/mnt/ct-03/hdd1/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd1/tmp"
	case "/mnt/ct-03/hdd2/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd2/tmp"
	case "/mnt/ct-03/hdd3/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd3/tmp"
	case "/mnt/ct-03/hdd4/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd4/tmp"
	case "/mnt/ct-03/hdd5/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd5/tmp"
	case "/mnt/ct-03/hdd6/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd6/tmp"
	case "/mnt/ct-03/hdd7/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd7/tmp"
	case "/mnt/ct-03/hdd8/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd8/tmp"
	case "/mnt/ct-03/hdd9/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd9/tmp"
	case "/mnt/ct-03/hdd10/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd10/tmp"
	case "/mnt/ct-03/hdd11/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd11/tmp"
	case "/mnt/ct-03/hdd12/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd12/tmp"
	case "/mnt/ct-03/hdd13/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd13/tmp"
	case "/mnt/ct-03/hdd14/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd14/tmp"
	case "/mnt/ct-03/hdd15/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd15/tmp"
	case "/mnt/ct-03/hdd16/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd16/tmp"
	case "/mnt/ct-03/hdd17/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd17/tmp"
	case "/mnt/ct-03/hdd18/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd18/tmp"
	case "/mnt/ct-03/hdd19/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd19/tmp"
	case "/mnt/ct-03/hdd20/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd20/tmp"
	case "/mnt/ct-03/hdd21/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd21/tmp"
	case "/mnt/ct-03/hdd22/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd22/tmp"
	case "/mnt/ct-03/hdd23/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd23/tmp"
	case "/mnt/ct-03/hdd24/":
		rsyncDestPath = "ct@192.168.4.3:/mnt/hdd24/tmp"
	case "/mnt/ct-04/hdd1/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd1/tmp"
	case "/mnt/ct-04/hdd2/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd2/tmp"
	case "/mnt/ct-04/hdd3/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd3/tmp"
	case "/mnt/ct-04/hdd4/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd4/tmp"
	case "/mnt/ct-04/hdd5/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd5/tmp"
	case "/mnt/ct-04/hdd6/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd6/tmp"
	case "/mnt/ct-04/hdd7/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd7/tmp"
	case "/mnt/ct-04/hdd8/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd8/tmp"
	case "/mnt/ct-04/hdd9/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd9/tmp"
	case "/mnt/ct-04/hdd10/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd10/tmp"
	case "/mnt/ct-04/hdd11/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd11/tmp"
	case "/mnt/ct-04/hdd12/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd12/tmp"
	case "/mnt/ct-04/hdd13/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd13/tmp"
	case "/mnt/ct-04/hdd14/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd14/tmp"
	case "/mnt/ct-04/hdd15/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd15/tmp"
	case "/mnt/ct-04/hdd16/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd16/tmp"
	case "/mnt/ct-04/hdd17/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd17/tmp"
	case "/mnt/ct-04/hdd18/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd18/tmp"
	case "/mnt/ct-04/hdd19/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd19/tmp"
	case "/mnt/ct-04/hdd20/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd20/tmp"
	case "/mnt/ct-04/hdd21/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd21/tmp"
	case "/mnt/ct-04/hdd22/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd22/tmp"
	case "/mnt/ct-04/hdd23/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd23/tmp"
	case "/mnt/ct-04/hdd24/":
		rsyncDestPath = "ct@192.168.4.4:/mnt/hdd24/tmp"
	case "/mnt/ct-05/hdd1/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd1/tmp"
	case "/mnt/ct-05/hdd2/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd2/tmp"
	case "/mnt/ct-05/hdd3/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd3/tmp"
	case "/mnt/ct-05/hdd4/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd4/tmp"
	case "/mnt/ct-05/hdd5/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd5/tmp"
	case "/mnt/ct-05/hdd6/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd6/tmp"
	case "/mnt/ct-05/hdd7/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd7/tmp"
	case "/mnt/ct-05/hdd8/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd8/tmp"
	case "/mnt/ct-05/hdd9/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd9/tmp"
	case "/mnt/ct-05/hdd10/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd10/tmp"
	case "/mnt/ct-05/hdd11/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd11/tmp"
	case "/mnt/ct-05/hdd12/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd12/tmp"
	case "/mnt/ct-05/hdd13/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd13/tmp"
	case "/mnt/ct-05/hdd14/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd14/tmp"
	case "/mnt/ct-05/hdd15/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd15/tmp"
	case "/mnt/ct-05/hdd16/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd16/tmp"
	case "/mnt/ct-05/hdd17/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd17/tmp"
	case "/mnt/ct-05/hdd18/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd18/tmp"
	case "/mnt/ct-05/hdd19/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd19/tmp"
	case "/mnt/ct-05/hdd20/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd20/tmp"
	case "/mnt/ct-05/hdd21/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd21/tmp"
	case "/mnt/ct-05/hdd22/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd22/tmp"
	case "/mnt/ct-05/hdd23/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd23/tmp"
	case "/mnt/ct-05/hdd24/":
		rsyncDestPath = "ct@192.168.4.5:/mnt/hdd24/tmp"
	case "/mnt/ct-06/hdd1/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd1/tmp"
	case "/mnt/ct-06/hdd2/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd2/tmp"
	case "/mnt/ct-06/hdd3/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd3/tmp"
	case "/mnt/ct-06/hdd4/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd4/tmp"
	case "/mnt/ct-06/hdd5/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd5/tmp"
	case "/mnt/ct-06/hdd6/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd6/tmp"
	case "/mnt/ct-06/hdd7/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd7/tmp"
	case "/mnt/ct-06/hdd8/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd8/tmp"
	case "/mnt/ct-06/hdd9/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd9/tmp"
	case "/mnt/ct-06/hdd10/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd10/tmp"
	case "/mnt/ct-06/hdd11/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd11/tmp"
	case "/mnt/ct-06/hdd12/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd12/tmp"
	case "/mnt/ct-06/hdd13/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd13/tmp"
	case "/mnt/ct-06/hdd14/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd14/tmp"
	case "/mnt/ct-06/hdd15/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd15/tmp"
	case "/mnt/ct-06/hdd16/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd16/tmp"
	case "/mnt/ct-06/hdd17/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd17/tmp"
	case "/mnt/ct-06/hdd18/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd18/tmp"
	case "/mnt/ct-06/hdd19/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd19/tmp"
	case "/mnt/ct-06/hdd20/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd20/tmp"
	case "/mnt/ct-06/hdd21/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd21/tmp"
	case "/mnt/ct-06/hdd22/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd22/tmp"
	case "/mnt/ct-06/hdd23/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd23/tmp"
	case "/mnt/ct-06/hdd24/":
		rsyncDestPath = "ct@192.168.4.6:/mnt/hdd24/tmp"
	case "/mnt/ct-07/hdd1/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd1/tmp"
	case "/mnt/ct-07/hdd2/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd2/tmp"
	case "/mnt/ct-07/hdd3/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd3/tmp"
	case "/mnt/ct-07/hdd4/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd4/tmp"
	case "/mnt/ct-07/hdd5/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd5/tmp"
	case "/mnt/ct-07/hdd6/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd6/tmp"
	case "/mnt/ct-07/hdd7/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd7/tmp"
	case "/mnt/ct-07/hdd8/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd8/tmp"
	case "/mnt/ct-07/hdd9/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd9/tmp"
	case "/mnt/ct-07/hdd10/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd10/tmp"
	case "/mnt/ct-07/hdd11/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd11/tmp"
	case "/mnt/ct-07/hdd12/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd12/tmp"
	case "/mnt/ct-07/hdd13/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd13/tmp"
	case "/mnt/ct-07/hdd14/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd14/tmp"
	case "/mnt/ct-07/hdd15/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd15/tmp"
	case "/mnt/ct-07/hdd16/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd16/tmp"
	case "/mnt/ct-07/hdd17/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd17/tmp"
	case "/mnt/ct-07/hdd18/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd18/tmp"
	case "/mnt/ct-07/hdd19/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd19/tmp"
	case "/mnt/ct-07/hdd20/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd20/tmp"
	case "/mnt/ct-07/hdd21/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd21/tmp"
	case "/mnt/ct-07/hdd22/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd22/tmp"
	case "/mnt/ct-07/hdd23/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd23/tmp"
	case "/mnt/ct-07/hdd24/":
		rsyncDestPath = "ct@192.168.4.7:/mnt/hdd24/tmp"
	case "/mnt/ct-08/hdd1/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd1/tmp"
	case "/mnt/ct-08/hdd2/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd2/tmp"
	case "/mnt/ct-08/hdd3/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd3/tmp"
	case "/mnt/ct-08/hdd4/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd4/tmp"
	case "/mnt/ct-08/hdd5/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd5/tmp"
	case "/mnt/ct-08/hdd6/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd6/tmp"
	case "/mnt/ct-08/hdd7/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd7/tmp"
	case "/mnt/ct-08/hdd8/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd8/tmp"
	case "/mnt/ct-08/hdd9/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd9/tmp"
	case "/mnt/ct-08/hdd10/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd10/tmp"
	case "/mnt/ct-08/hdd11/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd11/tmp"
	case "/mnt/ct-08/hdd12/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd12/tmp"
	case "/mnt/ct-08/hdd13/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd13/tmp"
	case "/mnt/ct-08/hdd14/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd14/tmp"
	case "/mnt/ct-08/hdd15/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd15/tmp"
	case "/mnt/ct-08/hdd16/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd16/tmp"
	case "/mnt/ct-08/hdd17/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd17/tmp"
	case "/mnt/ct-08/hdd18/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd18/tmp"
	case "/mnt/ct-08/hdd19/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd19/tmp"
	case "/mnt/ct-08/hdd20/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd20/tmp"
	case "/mnt/ct-08/hdd21/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd21/tmp"
	case "/mnt/ct-08/hdd22/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd22/tmp"
	case "/mnt/ct-08/hdd23/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd23/tmp"
	case "/mnt/ct-08/hdd24/":
		rsyncDestPath = "ct@192.168.4.8:/mnt/hdd24/tmp"
	case "/mnt/ct-09/hdd1/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd1/tmp"
	case "/mnt/ct-09/hdd2/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd2/tmp"
	case "/mnt/ct-09/hdd3/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd3/tmp"
	case "/mnt/ct-09/hdd4/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd4/tmp"
	case "/mnt/ct-09/hdd5/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd5/tmp"
	case "/mnt/ct-09/hdd6/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd6/tmp"
	case "/mnt/ct-09/hdd7/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd7/tmp"
	case "/mnt/ct-09/hdd8/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd8/tmp"
	case "/mnt/ct-09/hdd9/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd9/tmp"
	case "/mnt/ct-09/hdd10/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd10/tmp"
	case "/mnt/ct-09/hdd11/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd11/tmp"
	case "/mnt/ct-09/hdd12/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd12/tmp"
	case "/mnt/ct-09/hdd13/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd13/tmp"
	case "/mnt/ct-09/hdd14/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd14/tmp"
	case "/mnt/ct-09/hdd15/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd15/tmp"
	case "/mnt/ct-09/hdd16/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd16/tmp"
	case "/mnt/ct-09/hdd17/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd17/tmp"
	case "/mnt/ct-09/hdd18/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd18/tmp"
	case "/mnt/ct-09/hdd19/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd19/tmp"
	case "/mnt/ct-09/hdd20/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd20/tmp"
	case "/mnt/ct-09/hdd21/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd21/tmp"
	case "/mnt/ct-09/hdd22/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd22/tmp"
	case "/mnt/ct-09/hdd23/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd23/tmp"
	case "/mnt/ct-09/hdd24/":
		rsyncDestPath = "ct@192.168.4.9:/mnt/hdd24/tmp"
	case "/mnt/ct-10/hdd1/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd1/tmp"
	case "/mnt/ct-10/hdd2/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd2/tmp"
	case "/mnt/ct-10/hdd3/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd3/tmp"
	case "/mnt/ct-10/hdd4/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd4/tmp"
	case "/mnt/ct-10/hdd5/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd5/tmp"
	case "/mnt/ct-10/hdd6/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd6/tmp"
	case "/mnt/ct-10/hdd7/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd7/tmp"
	case "/mnt/ct-10/hdd8/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd8/tmp"
	case "/mnt/ct-10/hdd9/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd9/tmp"
	case "/mnt/ct-10/hdd10/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd10/tmp"
	case "/mnt/ct-10/hdd11/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd11/tmp"
	case "/mnt/ct-10/hdd12/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd12/tmp"
	case "/mnt/ct-10/hdd13/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd13/tmp"
	case "/mnt/ct-10/hdd14/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd14/tmp"
	case "/mnt/ct-10/hdd15/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd15/tmp"
	case "/mnt/ct-10/hdd16/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd16/tmp"
	case "/mnt/ct-10/hdd17/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd17/tmp"
	case "/mnt/ct-10/hdd18/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd18/tmp"
	case "/mnt/ct-10/hdd19/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd19/tmp"
	case "/mnt/ct-10/hdd20/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd20/tmp"
	case "/mnt/ct-10/hdd21/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd21/tmp"
	case "/mnt/ct-10/hdd22/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd22/tmp"
	case "/mnt/ct-10/hdd23/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd23/tmp"
	case "/mnt/ct-10/hdd24/":
		rsyncDestPath = "ct@192.168.4.10:/mnt/hdd24/tmp"
	default:
		return "", errors.New("unregistered destination")
	}

	return rsyncDestPath, nil
}
