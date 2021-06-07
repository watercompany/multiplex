package worker_test

import (
	"fmt"
	"testing"

	"github.com/watercompany/multiplex/worker"
)

// TODO: add new tests after crunch mode.
// func TestRunExecutable(t *testing.T) {
// 	tests := []struct {
// 		scenario string
// 	}{
// 		{
// 			scenario: "default worker cfg",
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			var cfg worker.WorkerCfg

// 			cfg.GetWorkerCfg()
// 			fmt.Printf("%+v\n", cfg)

// 			err := worker.RunExecutable(&cfg, "test")
// 			if err != nil {
// 				fmt.Printf("%v", err)
// 			}
// 		})
// 	}
// }

func TestGetPOSCfg(t *testing.T) {
	tests := []struct {
		scenario string
	}{
		{
			scenario: "default worker cfg",
		},
	}
	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			tempDir := "/mnt/ssd1/plotfiles/temp"
			temp2Dir := "/mnt/ssd1/plotfiles/temp2"
			finalDir := "/mnt/ssd1/plotfiles/final"
			args, err := worker.GetPOSArgs(tempDir, temp2Dir, finalDir)
			if err != nil {
				t.Errorf("err=%v\n", err)
			}
			fmt.Printf("%+v\n", args)
		})
	}
}
