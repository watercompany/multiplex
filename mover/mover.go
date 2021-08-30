package mover

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	TransferLockName = "transfer-lock"
)

func MoveFile(sourcePath, destPath, filename string) error {
	if !strings.HasSuffix(sourcePath, "/") {
		sourcePath = sourcePath + "/"
	}

	if !strings.HasSuffix(destPath, "/") {
		destPath = destPath + "/"
	}

	// Make moving dir if does not exist
	// for source path
	MakeTempDir(sourcePath, "moving")

	// Make tmp dir if does not exist
	// for destination path
	MakeTempDir(destPath, "tmp")

	tmpSrcPath := sourcePath + "moving/" + filename
	sourcePath = sourcePath + filename
	destDir := destPath
	tmpDestPath := destDir + "tmp/" + filename
	destPath = destPath + filename

	// Rename source path from root to moving
	err := os.Rename(sourcePath, tmpSrcPath)
	if err != nil {
		return fmt.Errorf("failed renaming final file: %s", err)
	}

	// Create lock files
	// transfer-lock-[unix-nano-time]
	transferLock := fmt.Sprintf("%s-%v", TransferLockName, time.Now().UnixNano())
	tfl, err := os.Create(destDir + transferLock)
	if err != nil {
		// Rename source path from temp src to src
		err1 := os.Rename(tmpSrcPath, sourcePath)
		if err1 != nil {
			return fmt.Errorf("failed creating transfer lock file: %v: failed renaming final file back to source: %s", err, err1)
		}

		return fmt.Errorf("failed creating transfer lock file: %s", err)
	}
	tfl.Chmod(0777)
	tfl.Close()

	// OLD COPY
	// inputFile, err := os.Open(tmpSrcPath)
	// if err != nil {
	// 	// Rename source path from temp src to src
	// 	err1 := os.Rename(tmpSrcPath, sourcePath)
	// 	if err1 != nil {
	// 		return fmt.Errorf("couldn't open source file: %v: failed renaming final file back to source: %s", err, err1)
	// 	}

	// 	return fmt.Errorf("couldn't open source file: %s", err)
	// }
	// outputFile, err := os.Create(tmpDestPath)
	// if err != nil {
	// 	inputFile.Close()
	// 	// Rename source path from temp src to src
	// 	err1 := os.Rename(tmpSrcPath, sourcePath)
	// 	if err1 != nil {
	// 		return fmt.Errorf("couldn't open dest file: %v: failed renaming final file back to source: %s", err, err1)
	// 	}

	// 	return fmt.Errorf("couldn't open dest file: %s", err)
	// }
	// defer outputFile.Close()
	// _, err = io.Copy(outputFile, inputFile)
	// inputFile.Close()
	// if err != nil {
	// 	// Rename source path from temp src to src
	// 	err1 := os.Rename(tmpSrcPath, sourcePath)
	// 	if err1 != nil {
	// 		return fmt.Errorf("writing to output file failed: %v: failed renaming final file back to source: %s", err, err1)
	// 	}

	// 	return fmt.Errorf("writing to output file failed: %s", err)
	// }

	// COPY USING RSYNC
	err = RunRsyncCmd(tmpSrcPath, destDir)
	if err != nil {
		// Rename source path from temp src to src
		err1 := os.Rename(tmpSrcPath, sourcePath)
		if err1 != nil {
			return fmt.Errorf("failed copying file using rsync: %v: failed renaming final file back to source: %s", err, err1)
		}

		return fmt.Errorf("failed copying file using rsync: %s", err)
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(tmpSrcPath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}

	// Rename dest path from tmp to root
	err = os.Rename(tmpDestPath, destPath)
	if err != nil {
		return fmt.Errorf("failed renaming final file: %s", err)
	}

	// Delete transfer lock file
	err = os.Remove(destDir + transferLock)
	if err != nil {
		return fmt.Errorf("failed deleting transfer lock file: %s", err)
	}

	return nil
}

func MoveFileV2(sourcePath, destPath, filename string) error {
	if !strings.HasSuffix(sourcePath, "/") {
		sourcePath = sourcePath + "/"
	}

	if !strings.HasSuffix(destPath, "/") {
		destPath = destPath + "/"
	}

	// Make moving dir if does not exist
	// for source path
	MakeTempDir(sourcePath, "moving")

	// Make tmp dir if does not exist
	// for destination path
	MakeTempDir(destPath, "tmp")

	tmpSrcPath := sourcePath + "moving/" + filename
	sourcePath = sourcePath + filename
	destDir := destPath
	tmpDestPath := destDir + "tmp/" + filename
	destPath = destPath + filename

	// Rename source path from root to moving
	err := os.Rename(sourcePath, tmpSrcPath)
	if err != nil {
		return fmt.Errorf("failed renaming final file: %s", err)
	}

	// Create lock files
	// transfer-lock-[unix-nano-time]
	transferLock := fmt.Sprintf("%s-%v", TransferLockName, time.Now().UnixNano())
	tfl, err := os.Create(destDir + transferLock)
	if err != nil {
		return fmt.Errorf("failed creating transfer lock file: %s", err)
	}
	tfl.Close()

	inputFile, err := os.Open(tmpSrcPath)
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
	err = os.Remove(tmpSrcPath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}

	// Rename dest path from tmp to root
	err = os.Rename(tmpDestPath, destPath)
	if err != nil {
		return fmt.Errorf("failed renaming final file: %s", err)
	}

	// Delete transfer lock file
	err = os.Remove(destDir + transferLock)
	if err != nil {
		return fmt.Errorf("failed deleting transfer lock file: %s", err)
	}

	return nil
}

func MakeTempDir(dir string, name string) error {
	dir = dir + name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0777)
		return fmt.Errorf("failed to make folder: %s", err)
	}
	return nil
}

func RunRsyncCmd(sourcePath, destPath string) error {
	rsyncDestPath, err := GetRsyncDestPath(destPath)
	if err != nil {
		return fmt.Errorf("err running rsync cmd: %v", err)
	}

	cmd := exec.Command("sudo", "rsync",
		"-a", sourcePath, rsyncDestPath)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("err running rsync cmd: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("err running rsync cmd: %v", err)
	}
	return nil
}
