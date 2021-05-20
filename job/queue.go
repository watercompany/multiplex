package job

import (
	"context"

	workerclient "github.com/watercompany/multiplex/worker/client"
)

func AddJob(wCfg workerclient.CallWorkerConfig) error {
	ctx := context.Background()
	dbClient, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	pushWorkerCfg(ctx, dbClient, wCfg)

	return nil
}

func GetJob() (*workerclient.CallWorkerConfig, error) {
	ctx := context.Background()
	dbClient, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	var workerCfg *workerclient.CallWorkerConfig
	workerCfg, err = popWorkerCfg(ctx, dbClient)
	if err != nil {
		return workerCfg, err
	}

	return workerCfg, nil
}

func ListAllJobs() (map[string]string, error) {
	ctx := context.Background()
	dbClient, err := ConnectDB()
	if err != nil {
		return map[string]string{}, err
	}

	kv, err := getAll(ctx, dbClient)
	if err != nil {
		return map[string]string{}, err
	}

	return kv, nil
}
