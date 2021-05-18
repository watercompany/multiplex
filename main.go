package main

import (
	"flag"

	"github.com/watercompany/multiplex/worker"
	"github.com/watercompany/multiplex/worker/client"
)

var (
	LogName       string
	WorkerPortNum string
)

// for flags
func init() {
	flag.StringVar(&LogName, "log-name", "test", "output log name")
	flag.StringVar(&WorkerPortNum, "worker-port", "9090", "port number a worker is listening to")
}

func main() {
	flag.Parse()
	testCfg := client.CallWorkerConfig{
		LogName: LogName,
	}

	var res *worker.Result
	client.CallWorker(testCfg, ":"+WorkerPortNum, res)
}
