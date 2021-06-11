package job_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/watercompany/multiplex/job"
	"github.com/watercompany/multiplex/worker"
	workerclient "github.com/watercompany/multiplex/worker/client"
)

func TestQueue_AddJob(t *testing.T) {
	tests := []struct {
		scenario  string
		workerCfg workerclient.CallWorkerConfig
	}{
		{
			scenario: "default",
			workerCfg: workerclient.CallWorkerConfig{
				LogName:  "testLog",
				TaskName: "pos",
				WorkerCfg: worker.WorkerCfg{
					ExecDir:   "/execDir",
					OutputDir: "/OutDir",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			flushDB(t)

			err := job.AddJob(tc.workerCfg, "")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			kv, err := job.ListAllJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("All data inside database:\n")
			for key, val := range kv {
				t.Logf("%v:%v\n", key, val)
			}
		})
	}
}

func TestQueue_GetJob(t *testing.T) {
	tests := []struct {
		scenario  string
		workerCfg workerclient.CallWorkerConfig
	}{
		{
			scenario: "default",
			workerCfg: workerclient.CallWorkerConfig{
				LogName:  "testLog",
				TaskName: "pos",
				WorkerCfg: worker.WorkerCfg{
					ExecDir:   "/execDir",
					OutputDir: "/OutDir",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			flushDB(t)

			for i := 1; i < 5; i++ {
				saveCfg := tc.workerCfg
				saveCfg.LogName = tc.workerCfg.LogName + fmt.Sprintf("-%v", i)
				err := job.AddJob(saveCfg, "")
				if err != nil {
					t.Fatalf("err: %v\n", err)
				}
			}

			wCfg, err := job.GetJob()
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}
			t.Logf("got job: %v\n", wCfg)

			kv, err := job.ListAllJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("All data inside database:\n")
			for key, val := range kv {
				t.Logf("%v:%v\n", key, val)
			}
		})
	}
}

func TestQueue_ListAllJob(t *testing.T) {
	tests := []struct {
		scenario  string
		workerCfg workerclient.CallWorkerConfig
	}{
		{
			scenario: "default",
			workerCfg: workerclient.CallWorkerConfig{
				LogName:  "testLog",
				TaskName: "pos",
				WorkerCfg: worker.WorkerCfg{
					ExecDir:   "/execDir",
					OutputDir: "/OutDir",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			flushDB(t)

			for i := 1; i < 5; i++ {
				saveCfg := tc.workerCfg
				saveCfg.LogName = tc.workerCfg.LogName + fmt.Sprintf("-%v", i)
				err := job.AddJob(saveCfg, "")
				if err != nil {
					t.Fatalf("err: %v\n", err)
				}
			}

			kv, err := job.ListAllJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("All data inside database:\n")
			for key, val := range kv {
				t.Logf("%v:%v\n", key, val)
			}
		})
	}
}

func TestQueue_GetNumberOfQueuedJobs(t *testing.T) {
	tests := []struct {
		scenario     string
		workerCfg    workerclient.CallWorkerConfig
		numberOfJobs int
	}{
		{
			scenario: "default",
			workerCfg: workerclient.CallWorkerConfig{
				LogName:  "testLog",
				TaskName: "pos",
				WorkerCfg: worker.WorkerCfg{
					ExecDir:   "/execDir",
					OutputDir: "/OutDir",
				},
			},
			numberOfJobs: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			flushDB(t)

			for i := 1; i <= tc.numberOfJobs; i++ {
				saveCfg := tc.workerCfg
				saveCfg.LogName = tc.workerCfg.LogName + fmt.Sprintf("-%v", i)
				err := job.AddJob(saveCfg, "")
				if err != nil {
					t.Fatalf("err: %v\n", err)
				}
			}

			kv, err := job.ListAllJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("All data inside database:\n")
			for key, val := range kv {
				t.Logf("%v:%v\n", key, val)
			}

			numberOfJobs, err := job.GetNumberOfQueuedJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("Number of Queued Jobs: %v\n", numberOfJobs)
			if tc.numberOfJobs != numberOfJobs {
				t.Errorf("want count %v, got %v", tc.numberOfJobs, numberOfJobs)
			}
		})
	}
}

func TestQueue_GetNumberOfActiveJobs(t *testing.T) {
	tests := []struct {
		scenario     string
		workerCfg    workerclient.CallWorkerConfig
		numberOfJobs int
	}{
		{
			scenario: "default",
			workerCfg: workerclient.CallWorkerConfig{
				LogName:  "testLog",
				TaskName: "pos",
				WorkerCfg: worker.WorkerCfg{
					ExecDir:   "/execDir",
					OutputDir: "/OutDir",
				},
			},
			numberOfJobs: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			flushDB(t)

			for i := 1; i <= tc.numberOfJobs; i++ {
				saveCfg := tc.workerCfg
				saveCfg.LogName = tc.workerCfg.LogName + fmt.Sprintf("-%v", i)
				err := job.IncrActiveJobs()
				if err != nil {
					t.Fatalf("err: %v\n", err)
				}
			}

			kv, err := job.ListAllJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("All data inside database:\n")
			for key, val := range kv {
				t.Logf("%v:%v\n", key, val)
			}

			numberOfJobs, err := job.GetNumberOfActiveJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("Number of Active Jobs: %v\n", numberOfJobs)
			if tc.numberOfJobs != numberOfJobs {
				t.Errorf("want count %v, got %v", tc.numberOfJobs, numberOfJobs)
			}

			// Decrement once
			err = job.DecrActiveJobs()
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			numberOfJobs, err = job.GetNumberOfActiveJobs("")
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			t.Logf("Number of Active Jobs: %v\n", numberOfJobs)
			if tc.numberOfJobs-1 != numberOfJobs {
				t.Errorf("want count %v, got %v", tc.numberOfJobs-1, numberOfJobs)
			}
		})
	}
}

func flushDB(t *testing.T) {
	ctx := context.Background()
	db, err := job.ConnectDB("")
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	err = db.FlushDB(ctx).Err()
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}
}
