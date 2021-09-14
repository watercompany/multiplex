package worker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

	plotGraphData, err := RunExecutable(args.TaskName, args.NumactlArg, finalArg...)
	if err != nil {
		log.Printf("error running exec: %v", err)
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error running exec: %v", err))
	}

	// Save plot graph data json
	jsonOutputName := fmt.Sprintf("%v-%v-graph-data", timeNowFormatted, args.LogName)
	err = SavePlotGraphDataToJSON(plotGraphData, wCfg.OutputDir+jsonOutputName)
	if err != nil {
		log.Printf("error saving plot graph data json: %v", err)
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("error saving plot graph data json: %v", err))
	}

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

func RunExecutable(taskName string, numactlArgs []string, args ...string) (PlotGraph, error) {
	var plotGraph PlotGraph

	preArgs := fmt.Sprintf("%s %s %s time", numactlArgs[0], numactlArgs[1], numactlArgs[2])
	log.Printf("Command executed: %s %v\n", preArgs, args)
	execArgs := []string{}
	execArgs = append(execArgs, numactlArgs[1])
	execArgs = append(execArgs, numactlArgs[2])
	execArgs = append(execArgs, "time")
	execArgs = append(execArgs, args...)
	// execArgs := append(preArgs, args...)
	cmd := exec.Command(numactlArgs[0], execArgs...)

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("error creating StdoutPipe for cmd: %v", err)
		return PlotGraph{}, fmt.Errorf("error creating StdoutPipe for cmd: %v", err)
	}
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("error creating StderrPipe for cmd: %v", err)
		return PlotGraph{}, fmt.Errorf("error creating StderrPipe for cmd: %v", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	errScanner := bufio.NewScanner(cmdErrReader)
	go func(taskName string) {
		var scannedStr string
		for scanner.Scan() {
			scannedStr = scanner.Text()
			log.Printf("%s\n", scannedStr)

			if strings.Contains(scannedStr, phaseStringTemplateV2) {
				phase, hours, err := ParseGraphData(scannedStr, taskName)
				if err != nil {
					log.Printf("err=%s\n", err)
					return
				}

				// append data
				plotGraph.Data = append(plotGraph.Data, PlotData{Phase: phase, Hours: hours})
			}
		}

		for errScanner.Scan() {
			scannedStr = errScanner.Text()
			log.Printf("%s\n", scannedStr)
		}
	}(taskName)

	err = cmd.Start()
	if err != nil {
		log.Printf("error starting cmd:%v\n", err)
		return PlotGraph{}, fmt.Errorf("error starting cmd: %v", err)
	}

	log.Printf("Process PID: %v\n", cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		log.Printf("error waiting for cmd: %v", err)
		return PlotGraph{}, fmt.Errorf("error waiting for cmd: %v", err)
	}

	return plotGraph, nil
}
