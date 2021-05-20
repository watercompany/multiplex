# multiplex
The Multiplex

## Dependencies Installation
```
go mod verify 
go mod tidy
go mod download
```

## Usage 
Follow these steps for usage:
    - Deploy workers first
    - Run Redis Database Instance
    - Run Dispatcher
    - Give Jobs to Workers
    - View Queued Jobs

### Deployment of Worker
Use this command to deploy workers.
```
WORKERS=[number of workers to deploy] make deploy-workers
```
###### Example

```
WORKERS=5 make deploy-workers
```

### Run Redis Instance
```
docker pull redis
docker run --name redis-test-instance -p 6379:6379 -d redis
```

### Run Dispatcher
Use this command to run the dispatcher.
```
go run ./job/dispatcher/main.go
```

### Give job to worker
After configuring worker.yml for executable and logging directory
and pos.yml for proof of space task, run using this command
```
go run main.go  -task=[task name] -log-name=[log-output-name] -worker-port=[port number of a worker] &
```
#### Example
For test exec
```
go run main.go  -task=test -log-name=test-drive-1 -worker-port=9090 &
```

For Proof of Space exec
```
go run main.go  -task=pos -log-name=test-drive-1 -worker-port=9090 &
```

