package worker

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/watercompany/multiplex/mover"
)

func (pw *ProgramWorker) RunWorker(args *Args) (Result, *erpc.Status) {
	var wCfg WorkerCfg = args.WorkerCfg
	var addArgs []string = args.AdditionalArgs

	timeNow := time.Now()
	timeNowFormatted := timeNow.Format(time.RFC3339)
	outputName := fmt.Sprintf("%v-%v-plot-generation-log", timeNowFormatted, args.LogName)
	f, err := os.OpenFile(wCfg.OutputDir+outputName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	log.SetOutput(f)

	paramLogName := fmt.Sprintf("%v-%v-param-and-time-log", timeNowFormatted, args.LogName)
	paramLog, err := os.OpenFile(wCfg.OutputDir+paramLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error opening file: %v", err))
	}
	defer paramLog.Close()

	var finalArg []string
	finalArg = append(finalArg, wCfg.ExecDir)
	finalArg = append(finalArg, addArgs...)

	_, err = paramLog.Write([]byte(fmt.Sprintf("Command Executed: time %v\nTime Started: %v\n", args, timeNow)))
	if err != nil {
		log.Fatal(err)
	}

	// Check if temp dir have enough storage
	// 300GB free per plot
	freeSpaceInMB := mover.GetFreeDiskSpaceInMB(args.POSCfg.TempDir)
	if freeSpaceInMB < 300000 {
		log.Printf("error not enough free space: %v MB", freeSpaceInMB)
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error not enough free space: %v MB", freeSpaceInMB))
	}

	err = RunExecutable(finalArg...)
	if err != nil {
		log.Printf("error running exec: %v", err)
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error running exec: %v", err))
	}

	startMoveTime := time.Now()
	// Move
	// Copy final plot to somewhere
	// Delete final plot
	err = moveFinalPlot(args)
	if err != nil {
		log.Printf("error moving final plot: %v", err)
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error moving final plot: %v", err))
	}
	log.Printf("Final Plot has been moved to final destination dir.\nMoving Plot Took: %v minutes\n", time.Since(startMoveTime).Minutes())

	_, err = paramLog.Write([]byte(fmt.Sprintf("Time Finished: %v\n", time.Now())))
	if err != nil {
		log.Fatal(err)
	}

	_, err = paramLog.Write([]byte(fmt.Sprintf("Duration: %v minutes\n", time.Since(timeNow).Minutes())))
	if err != nil {
		log.Fatal(err)
	}

	res := Result{}
	return res, nil
}

func RunExecutable(args ...string) error {
	log.Printf("Command executed: time %v\n", args)
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
		return err
	}

	// Check if destpath have available size for that
	destFreeSpace := mover.GetFreeDiskSpaceInMB(args.POSCfg.FinalDestDir)
	if dirSize > int64(destFreeSpace) {
		return errors.New("file size greater than destination free space")
	}

	// Moves and Deletes
	err = mover.MoveFile(args.POSCfg.FinalDir, args.POSCfg.FinalDestDir, args.POSCfg.FileName)
	if err != nil {
		return err
	}

	return nil
}
