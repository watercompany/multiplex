package job

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/watercompany/multiplex/worker"
	workerclient "github.com/watercompany/multiplex/worker/client"
)

func TestDatabase_PushWorkerCfg(t *testing.T) {
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
			ctx := context.Background()
			db, err := ConnectDB()
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}

			err = db.FlushDB(ctx).Err()
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			err = pushWorkerCfg(ctx, db, tc.workerCfg)
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}

			kv, err := getAll(ctx, db)
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}
			t.Logf("All data inside database:\n")
			for key, val := range kv {
				t.Logf("%v:%v\n", key, val)
			}

			wCfg, err := getWorkerCfg(ctx, db, "1")
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}

			if !reflect.DeepEqual(tc.workerCfg, *wCfg) {
				t.Errorf("want %v, got %v", tc.workerCfg, wCfg)
			}
		})
	}
}

func TestDatabase_Pop(t *testing.T) {
	testCfg := workerclient.CallWorkerConfig{
		LogName:  "testLog",
		TaskName: "pos",
		WorkerCfg: worker.WorkerCfg{
			ExecDir:   "/execDir",
			OutputDir: "/OutDir",
		},
	}

	tests := []struct {
		scenario string
	}{
		{
			scenario: "default",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.Background()
			db, err := ConnectDB()
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}

			err = db.FlushDB(ctx).Err()
			if err != nil {
				t.Fatalf("err: %v\n", err)
			}

			for i := 0; i < 3; i++ {
				err = pushWorkerCfg(ctx, db, testCfg)
				if err != nil {
					t.Fatalf("err:%v\n", err)
				}
			}

			workerCfg, err := popWorkerCfg(ctx, db)
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}
			t.Logf("popped worker cfg=%v\n", workerCfg)

			workerCfg, err = popWorkerCfg(ctx, db)
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}
			t.Logf("popped worker cfg=%v\n", workerCfg)

			workerCfg, err = popWorkerCfg(ctx, db)
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}
			t.Logf("popped worker cfg=%v\n", workerCfg)

			kv, err := getAll(ctx, db)
			if err != nil {
				t.Fatalf("err:%v\n", err)
			}
			t.Logf("All data inside database:\n")
			for key, val := range kv {
				t.Logf("%v:%v\n", key, val)
			}
		})
	}
}

func TestDatabase_GetEmpty(t *testing.T) {
	db, _ := ConnectDB()
	db.FlushAll(context.Background())
	_, err := get(context.Background(), db, jobLastIndex)
	if err != redis.Nil {
		t.Logf("want err redis.Nil, got %v", err)
	}
}
