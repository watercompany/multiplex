package mover

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// func Mover() {
// 	srcPath := "/Users/kenjehofilena/Desktop"
// 	destPath := "/Users/kenjehofilena/Documents"
// 	for {
// 		fc, err := fileCount(srcPath)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if fc > 0 {
// 			// Check src path size
// 			dirSize, err := DirSizeInMB(srcPath)
// 			if err != nil {
// 				panic(err)
// 			}

// 			// Check if destpath have available size for that
// 			destFreeSpace := getFreeDiskSpaceInMB(destPath)
// 			if dirSize > int64(destFreeSpace) {
// 				panic("file size greater than destination free space")
// 			}

// 			// Move
// 			MoveFile(srcPath, destPath)
// 		}
// 	}
// }

func MoveFile(sourcePath, destPath, filename string) error {
	if !strings.HasSuffix(sourcePath, "/") {
		sourcePath = sourcePath + "/"
	}

	if !strings.HasSuffix(destPath, "/") {
		destPath = destPath + "/"
	}

	// Make tmp dir if does not exist
	// for destination path
	MakeTempDir(destPath)

	sourcePath = sourcePath + filename
	tmpDestPath := destPath + "tmp/" + filename
	destPath = destPath + filename

	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(tmpDestPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}

	// Rename dest path from tmp to root
	err = os.Rename(tmpDestPath, destPath)
	if err != nil {
		return fmt.Errorf("failed renaming final file: %s", err)
	}
	return nil
}

func MakeTempDir(dir string) error {
	dir = dir + "tmp"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0777)
		return fmt.Errorf("failed to make tmp folder: %s", err)
	}
	return nil
}
