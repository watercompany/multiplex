package mover

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func FileCount(path string) (int, error) {
	d, e := os.ReadDir(path)
	if e != nil {
		return 0, fmt.Errorf("error counting files: %v", e)
	}

	println(len(d))
	return len(d), nil
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
	return size / 1024 / 1024, err
}
