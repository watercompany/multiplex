package worker

import "github.com/henrylee2cn/erpc/v6"

const (
	RunProgram     = "/program_worker/run_worker"
	BasePortNumber = 9090
)

type Args struct {
	LogName        string
	TaskName       string
	WorkerCfg      WorkerCfg
	AdditionalArgs []string
	POSCfg         POSCfg
}

type Result struct {
	Output float64
}
type ProgramWorker struct {
	erpc.CallCtx
}
