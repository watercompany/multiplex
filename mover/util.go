package mover

import (
	"fmt"
	"io/fs"
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
	// row-1 group priority
	for _, finalDir := range finalDirs {
		// Check if directory is row-1
		row1Count, err := FileCountSubString(finalDir, "row-1")
		if err != nil {
			return "", err
		}

		if row1Count == 0 {
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
		// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
		hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

		if hasAvailableSpace && !hasExceededMaxLocks {
			return finalDir, nil
		}
	}

	// row-2 group priority
	for _, finalDir := range finalDirs {
		// Check if directory is row-2
		row2Count, err := FileCountSubString(finalDir, "row-2")
		if err != nil {
			return "", err
		}

		if row2Count == 0 {
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
		// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
		hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

		if hasAvailableSpace && !hasExceededMaxLocks {
			return finalDir, nil
		}
	}

	// row-3 group priority
	for _, finalDir := range finalDirs {
		// Check if directory is row-3
		row3Count, err := FileCountSubString(finalDir, "row-3")
		if err != nil {
			return "", err
		}

		if row3Count == 0 {
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
		// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
		hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

		if hasAvailableSpace && !hasExceededMaxLocks {
			return finalDir, nil
		}
	}

	// no row group priority
	for _, finalDir := range finalDirs {
		// Check if exceeded max lock files
		count, err := FileCountSubString(finalDir, TransferLockName)
		if err != nil {
			return "", err
		}
		hasExceededMaxLocks := count >= maxLockFiles

		// Check if destpath have available size for that
		destFreeSpace := GetFreeDiskSpaceInMB(finalDir)
		// has available space if filesize+110GB+(current transfer count*102GB) is less than free space
		hasAvailableSpace := (fileSize + 110000 + (int64(count) * 102000)) < int64(destFreeSpace)

		if hasAvailableSpace && !hasExceededMaxLocks {
			return finalDir, nil
		}
	}

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
