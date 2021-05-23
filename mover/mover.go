package mover

import (
	"fmt"
	"io"
	"os"
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

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
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
	return nil
}
