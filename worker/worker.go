package worker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/henrylee2cn/erpc/v6"
)

const (
	RunProgram     = "/program_worker/run_worker"
	BasePortNumber = 9090
)

type Args struct {
	LogName string
}

type Result struct {
	Output float64
}
type ProgramWorker struct {
	erpc.CallCtx
}

func (pw *ProgramWorker) RunWorker(args *Args) (Result, *erpc.Status) {
	var wCfg WorkerCfg
	_, err := wCfg.GetWorkerCfg()
	if err != nil {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("%v", err))
	}

	err = RunExecutable(&wCfg, args.LogName)
	if err != nil {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("%v", err))
	}

	res := Result{}
	return res, nil
}

func RunExecutable(wCfg *WorkerCfg, logName string) error {
	outputName := fmt.Sprintf("%v-output-log", logName)
	f, err := os.OpenFile(wCfg.OutputDir+outputName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	cmd := exec.Command("time", wCfg.ExecDir)

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
		return fmt.Errorf("error starting cmd: %v", err)
	}

	log.Printf("Process PID: %v\n", cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for cmd: %v", err)
	}

	return nil
}