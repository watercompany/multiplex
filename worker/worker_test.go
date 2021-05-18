package worker_test

import (
	"fmt"
	"testing"

	"github.com/watercompany/multiplex/worker"
)

func TestGetAvailableWorkers(t *testing.T) {
	tests := []struct {
		scenario string
	}{
		{
			scenario: "1 worker",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			var cfg worker.WorkerCfg

			cfg.GetWorkerCfg()
			fmt.Printf("%+v\n", cfg)

			err := worker.RunExecutable(&cfg, "test")
			if err != nil {
				fmt.Printf("%v", err)
			}
		})
	}
}
