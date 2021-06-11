package mover

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

func FileCount(path string) (int, error) {
	d, e := os.ReadDir(path)
	if e != nil {
		return 0, fmt.Errorf("error counting files: %v", e)
	}
	return len(d), nil
}

func FileCountSubString(path string, subStr string) (int, error) {
	count := 0
	d, e := os.ReadDir(path)
	if e != nil {
		return 0, fmt.Errorf("error counting files: %v", e)
	}

	for _, dir := range d {
		if strings.Contains(dir.Name(), subStr) {
			count++
		}
	}

	return count, nil
}

func GetDirs(path string) ([]fs.DirEntry, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return []fs.DirEntry{}, fmt.Errorf("error counting files: %v", err)
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
	for _, finalDir := range finalDirs {
		// Check if destpath have available size for that
		destFreeSpace := GetFreeDiskSpaceInMB(finalDir)
		hasAvailableSpace := fileSize < int64(destFreeSpace)

		// Check if exceeded max lock files
		count, err := FileCountSubString(finalDir, transferLockName)
		if err != nil {
			return "", err
		}
		hasExceededMaxLocks := count >= maxLockFiles

		if hasAvailableSpace && !hasExceededMaxLocks {
			return finalDir, nil
		}
	}

	// No available dirs
	return "", nil
}
