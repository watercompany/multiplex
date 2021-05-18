package main

import (
	"flag"

	"github.com/watercompany/multiplex/worker"
	"github.com/watercompany/multiplex/worker/client"
)

var (
	LogName       string
	WorkerPortNum string
	TaskName      string
)

// for flags
func init() {
	flag.StringVar(&LogName, "log-name", "test", "output log name")
	flag.StringVar(&WorkerPortNum, "worker-port", "9090", "port number a worker is listening to")
	flag.StringVar(&TaskName, "task", "plot", "name of the task the worker will do")
}

func main() {
	flag.Parse()

	testCfg := client.CallWorkerConfig{
		LogName:  LogName,
		TaskName: TaskName,
	}

	var res *worker.Result
	client.CallWorker(testCfg, ":"+WorkerPortNum, res)
}
