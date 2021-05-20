package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	workerclient "github.com/watercompany/multiplex/worker/client"
)

// docker pull redis
// docker run --name redis-test-instance -p 6379:6379 -d redis
func ConnectDB() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := ping(context.Background(), client)
	if err != nil {
		return client, err
	}

	return client, nil
}

func ping(ctx context.Context, client *redis.Client) error {
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	// fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}

func set(ctx context.Context, client *redis.Client, key string, value interface{}) error {
	_, err := client.Set(ctx, key, value, 0).Result()
	if err != nil {
		return err
	}

	return nil
}

func get(ctx context.Context, client *redis.Client, key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func getWorkerCfg(ctx context.Context, client *redis.Client, key string) (*workerclient.CallWorkerConfig, error) {
	objStr, err := client.Get(ctx, key).Result()
	if err != nil {
		return &workerclient.CallWorkerConfig{}, err
	}

	b := []byte(objStr)
	workerCfg := &workerclient.CallWorkerConfig{}
	err = json.Unmarshal(b, workerCfg)
	if err != nil {
		return &workerclient.CallWorkerConfig{}, err
	}
	return workerCfg, nil
}

func getAll(ctx context.Context, client *redis.Client) (map[string]string, error) {
	kv := map[string]string{}
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return map[string]string{}, err
	}

	for _, key := range keys {
		value, err := get(ctx, client, key)
		if err != nil {
			return map[string]string{}, err
		}

		kv[key] = value
	}

	return kv, nil
}

func del(ctx context.Context, client *redis.Client, key string) error {
	_, err := client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func incrCounter(ctx context.Context, client *redis.Client) (string, error) {
	val, err := client.Incr(ctx, "job-last-index").Result()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", val), nil
}

func decrCounter(ctx context.Context, client *redis.Client) (string, error) {
	val, err := client.Decr(ctx, "job-last-index").Result()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", val), nil
}

func pushWorkerCfg(ctx context.Context, client *redis.Client, workerCfg workerclient.CallWorkerConfig) error {
	key, err := incrCounter(ctx, client)
	if err != nil {
		return err
	}

	json, err := json.Marshal(workerCfg)
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set(ctx, key, json, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func popWorkerCfg(ctx context.Context, client *redis.Client) (*workerclient.CallWorkerConfig, error) {
	var workerCfgVal *workerclient.CallWorkerConfig

	// Get Last Index
	lastIdx, err := client.Get(ctx, "job-last-index").Result()
	if err != nil {
		return workerCfgVal, err
	}

	if lastIdx == "0" {
		return workerCfgVal, errors.New("queue is empty")
	}

	// Get Value
	workerCfgVal, err = getWorkerCfg(ctx, client, "1")
	if err != nil {
		return workerCfgVal, err
	}

	// Delete Value
	err = client.Del(ctx, "1").Err()
	if err != nil {
		return workerCfgVal, err
	}

	// Decrement index
	_, err = decrCounter(ctx, client)
	if err != nil {
		return workerCfgVal, err
	}

	// Early return, doesnt need to move
	// indexes because theres no more value
	if lastIdx == "1" {
		return workerCfgVal, nil
	}

	idx, err := strconv.Atoi(lastIdx)
	if err != nil {
		return workerCfgVal, err
	}

	// move indexes, 5->4, 4->3, 3->2, 2->1
	for i := 1; i < idx; i++ {
		moveVal, err := getWorkerCfg(ctx, client, fmt.Sprintf("%v", i+1))
		if err != nil {
			return workerCfgVal, err
		}

		json, err := json.Marshal(moveVal)
		if err != nil {
			fmt.Println(err)
		}

		err = set(ctx, client, fmt.Sprintf("%v", i), json)
		if err != nil {
			return workerCfgVal, err
		}
	}

	err = del(ctx, client, lastIdx)
	if err != nil {
		return workerCfgVal, err
	}

	return workerCfgVal, nil
}

// func push(ctx context.Context, client *redis.Client, value string) error {
// 	key, err := incrCounter(ctx, client)
// 	if err != nil {
// 		return err
// 	}

// 	err = client.Set(ctx, key, value, 0).Err()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
