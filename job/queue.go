package job

import (
	"context"
	"strconv"

	workerclient "github.com/watercompany/multiplex/worker/client"
)

func AddJob(wCfg workerclient.CallWorkerConfig, databaseAddress string) error {
	ctx := context.Background()
	dbClient, err := ConnectDB(databaseAddress)
	if err != nil {
		panic(err)
	}
	pushWorkerCfg(ctx, dbClient, wCfg)

	return nil
}

func GetJob() (*workerclient.CallWorkerConfig, error) {
	var workerCfg *workerclient.CallWorkerConfig
	ctx := context.Background()
	dbClient, err := ConnectDB("")
	if err != nil {
		return workerCfg, err
	}

	workerCfg, err = popWorkerCfg(ctx, dbClient)
	if err != nil {
		return workerCfg, err
	}

	return workerCfg, nil
}

func GetNumberOfActiveJobs(databaseAddress string) (int, error) {
	ctx := context.Background()
	dbClient, err := ConnectDB(databaseAddress)
	if err != nil {
		return 0, err
	}

	countStr, err := get(ctx, dbClient, activeJobs)
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetNumberOfQueuedJobs(databaseAddress string) (int, error) {
	ctx := context.Background()
	dbClient, err := ConnectDB(databaseAddress)
	if err != nil {
		return 0, err
	}

	countStr, err := get(ctx, dbClient, jobLastIndex)
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func ListAllJobs(databaseAddress string) (map[string]string, error) {
	ctx := context.Background()
	dbClient, err := ConnectDB(databaseAddress)
	if err != nil {
		return map[string]string{}, err
	}

	kv, err := getAll(ctx, dbClient)
	if err != nil {
		return map[string]string{}, err
	}

	return kv, nil
}

func IncrActiveJobs() error {
	ctx := context.Background()
	dbClient, err := ConnectDB("")
	if err != nil {
		return err
	}

	_, err = incrCounter(ctx, dbClient, activeJobs)
	if err != nil {
		return err
	}

	return nil
}

func DecrActiveJobs() error {
	ctx := context.Background()
	dbClient, err := ConnectDB("")
	if err != nil {
		return err
	}

	_, err = decrCounter(ctx, dbClient, activeJobs)
	if err != nil {
		return err
	}

	return nil
}
