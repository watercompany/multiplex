package worker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/watercompany/multiplex/mover"
)

func (pw *ProgramWorker) RunWorker(args *Args) (Result, *erpc.Status) {
	var wCfg WorkerCfg = args.WorkerCfg
	var addArgs []string = args.AdditionalArgs

	outputName := fmt.Sprintf("%v-output-log", args.LogName)
	f, err := os.OpenFile(wCfg.OutputDir+outputName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	log.SetOutput(f)

	var finalArg []string
	finalArg = append(finalArg, wCfg.ExecDir)
	finalArg = append(finalArg, addArgs...)

	// Check if temp dir have enough storage
	// 300GB free per plot
	freeSpaceInMB := mover.GetFreeDiskSpaceInMB(args.POSCfg.TempDir)
	if freeSpaceInMB < 300000 {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error not enough free space: %v MB", freeSpaceInMB))
	}

	err = RunExecutable(finalArg...)
	if err != nil {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error running exec: %v", err))
	}
	// Move
	// Copy final plot to somewhere
	// Delete final plot
	err = moveFinalPlot(args)
	if err != nil {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error moving final plot: %v", err))
	}

	res := Result{}
	return res, nil
}

func RunExecutable(args ...string) error {
	cmd := exec.Command("time", args...)

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating StdoutPipe for cmd: %v", err)
	}
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating StderrPipe for cmd: %v", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for scanner.Scan() {
			log.Printf("%s\n", scanner.Text())
		}

		for errScanner.Scan() {
			log.Printf("%s\n", errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Printf("error starting cmd:%v\n", err)
		return fmt.Errorf("error starting cmd: %v", err)
	}

	log.Printf("Process PID: %v\n", cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for cmd: %v", err)
	}

	return nil
}

func moveFinalPlot(args *Args) error {
	// Check src path size
	dirSize, err := mover.DirSizeInMB(args.POSCfg.FinalDir)
	if err != nil {
		panic(err)
	}

	// Check if destpath have available size for that
	destFreeSpace := mover.GetFreeDiskSpaceInMB(args.POSCfg.FinalDestDir)
	if dirSize > int64(destFreeSpace) {
		panic("file size greater than destination free space")
	}

	// Moves and Deletes
	mover.MoveFile(args.POSCfg.FinalDir, args.POSCfg.FinalDestDir, args.POSCfg.FileName)

	return nil
}
