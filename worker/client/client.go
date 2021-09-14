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
	POSCfg         worker.POSCfg    `json:"pos_cfg"`
	NumactlArg     []string
}

func CallWorker(cWorker CallWorkerConfig, workerAddr string, result *worker.Result) {
	// erpc.SetLoggerLevel("ON")()
	cli := erpc.NewPeer(erpc.PeerConfig{RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()
	cli.SetTLSConfig(erpc.GenerateTLSConfigForClient())
	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(workerAddr)
	if !stat.OK() {
		erpc.Criticalf("%v", stat)
		return
	}
	defer sess.Close()

	// Set worker args
	args := &worker.Args{
		LogName:        cWorker.LogName,
		TaskName:       cWorker.TaskName,
		WorkerCfg:      cWorker.WorkerCfg,
		AdditionalArgs: cWorker.AdditionalArgs,
		POSCfg:         cWorker.POSCfg,
		NumactlArg:     cWorker.NumactlArg,
	}

	stat = sess.Call(
		worker.RunProgram,
		args,
		&result,
	).Status()

	if !stat.OK() {
		erpc.Criticalf("%v", stat)
		return
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
