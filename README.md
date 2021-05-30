# multiplex
The Multiplex

## Dependencies Installation
```
go mod verify 
go mod tidy
go mod download
```

## Manual Usage 
Follow these steps for manual usage:
1. Deploy workers first
2. Run Redis Database Instance
3. Run Dispatcher
4. Give Jobs to Workers
5. View Queued Jobs

Or run automated script that follows step 1-3:
WORKERS=[number of workers] make run-multiplex

### Deployment of Worker
Use this command to deploy workers.
```
WORKERS=[number of workers to deploy] make deploy-workers
```
###### Example
```
WORKERS=8 make deploy-workers
```

### Run Redis Instance
```
docker pull redis
docker run --name redis-test-instance -p 6379:6379 -d redis
```

### Run Dispatcher
Use this command to run the dispatcher.
```
go run ./job/dispatcher/main.go -workers=[number of available workers]
```
###### Example
```
go run ./job/dispatcher/main.go -workers=8 &
```

### Give job to worker
After configuring 
- worker.yml for executable and logging directory
- pos.yml for proof of space task, run using this command
```
go run main.go add-job --task=[task name] --log-name=[log-output-name]
```
##### Example
For test exec
```
go run main.go add-job  --task=test --log-name=test-drive-1 
```

For Proof of Space exec
```
go run main.go add-job  --task=pos --log-name=test-drive-1 
```

### View Queued Jobs
```
go run main.go view-jobs
```

## Docker Usage 
Follow these steps for dockerized usage:
1. Run make push-docker
2. Run docker-compose up
3. Configure pos.yml and worker.yml
4. Add job to worker