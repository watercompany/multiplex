# multiplex
The Multiplex

### Dependencies Installation
```
go mod verify 
go mod tidy
go mod download
```

### Deployment of Worker
Use this command to deploy workers.

```
WORKERS=[number of workers to deploy] make deploy-workers
```
#### Example

```
WORKERS=5 make deploy-workers
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

