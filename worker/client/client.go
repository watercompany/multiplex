package client

import (
	"time"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/watercompany/multiplex/worker"
)

type CallWorkerConfig struct {
	LogName        string           `json:"log_name"`
	TaskName       string           `json:"task_name"`
	WorkerCfg      worker.WorkerCfg `json:"worker_cfg"`
	AdditionalArgs []string         `json:"additional_args"`
	TempDir        string           `json:"temp_dir"`
	FinalDir       string           `json:"final_dir"`
	FinalDestDir   string           `json:"final_dest_dir"`
}

func CallWorker(cWorker CallWorkerConfig, workerAddr string, result *worker.Result) {
	erpc.SetLoggerLevel("OFF")()
	cli := erpc.NewPeer(erpc.PeerConfig{RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()
	cli.SetTLSConfig(erpc.GenerateTLSConfigForClient())
	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(workerAddr)
	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}
	defer sess.Close()

	// Set worker args
	args := &worker.Args{
		LogName:        cWorker.LogName,
		TaskName:       cWorker.TaskName,
		WorkerCfg:      cWorker.WorkerCfg,
		AdditionalArgs: cWorker.AdditionalArgs,
		TempDir:        cWorker.TempDir,
		FinalDir:       cWorker.FinalDir,
		FinalDestDir:   cWorker.FinalDestDir,
	}

	stat = sess.Call(
		worker.RunProgram,
		args,
		&result,
	).Status()

	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}
}

// Push push handler
type Push struct {
	erpc.PushCtx
}

// Push handles '/push/status' message
func (p *Push) Status(arg *string) *erpc.Status {
	erpc.Printf("%s", *arg)
	return nil
}
