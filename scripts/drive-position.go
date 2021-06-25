package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	outputStr, err := FindDrives()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Drives lsblk: \n%s\n", outputStr)
	err = ExtractAndPutDrivePositions(outputStr)
	if err != nil {
		panic(err)
	}
}

// run lsblk -o HCTL,NAME,MOUNTPOINT -l
func FindDrives() (string, error) {
	var outputStr string
	cmd := exec.Command("lsblk", "-o", "NAME,HCTL,MOUNTPOINT", "-l")

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("error creating StdoutPipe for cmd: %v", err)
		return "", fmt.Errorf("error creating StdoutPipe for cmd: %v", err)
	}
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("error creating StderrPipe for cmd: %v", err)
		return "", fmt.Errorf("error creating StderrPipe for cmd: %v", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		var scannedStr string
		for scanner.Scan() {
			scannedStr = scanner.Text()
			outputStr = outputStr + fmt.Sprintf("%s\n", scannedStr)
		}

		for errScanner.Scan() {
			scannedStr = errScanner.Text()
			outputStr = outputStr + fmt.Sprintf("%s\n", scannedStr)
		}
	}()

	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("error starting cmd: %v", err)
	}

	time.Sleep(10 * time.Second)
	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("error waiting for cmd: %v", err)
	}

	return outputStr, nil
}

func ExtractAndPutDrivePositions(in string) error {
	namesToFind := []string{
		"sda",
		"sdb",
		"sdc",
		"sdd",
		"sde",
		"sdf",
		"sdg",
		"sdh",
		"sdi",
		"sdj",
		"sdk",
		"sdl",
		"sdm",
		"sdn",
		"sdo",
		"sdp",
		"sdq",
		"sdr",
		"sds",
		"sdt",
		"sdu",
		"sdv",
		"sdw",
		"sdx",
	}

	for _, val := range namesToFind {
		host, err := parseHostFromHCTL(in, val)
		if err != nil {
			return err
		}

		mntPoint, err := parseMountPoint(in, val+"1")
		if err != nil {
			return err
		}

		fmt.Printf("Mount point: %s\n", mntPoint)
		fmt.Printf("Host: %v\n", host)
		rowNum := getRowNumberFromHostNum(host)
		fileStr := fmt.Sprintf("row-%v", rowNum)

		sameNameCount, err := FileCountSubString(mntPoint, fileStr)
		if err != nil {
			return err
		}

		// if theres the same name, dont make file
		if sameNameCount > 0 {
			continue
		} else {
			// make the file
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s", mntPoint, fileStr), []byte(""), 0777)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func parseHostFromHCTL(data, template string) (int, error) {
	index := strings.Index(data, template)
	hctlStr := getValueUntilSpace(data, index+len(template)+7)
	hostStr := hctlStr[0:2]
	host, err := strconv.Atoi(hostStr)
	if err != nil {
		return 0, err
	}
	return host, nil
}

func parseMountPoint(data, template string) (string, error) {
	index := strings.Index(data, template)
	mountPointStr := getValueUntilNewline(data, index+len(template)+17)

	return mountPointStr, nil
}

func getValueUntilSpace(info string, index int) string {
	end := strings.Index(info[index:], " ") + index
	return info[index:end]
}

func getValueUntilNewline(info string, index int) string {
	end := strings.Index(info[index:], "\n") + index
	return info[index:end]
}

func getRowNumberFromHostNum(host int) int {
	switch host {
	case 9, 10, 11, 12, 13, 14, 15, 16:
		return 1
	case 17:
		return 2
	case 18:
		return 3
	default:
		return 3
	}
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
