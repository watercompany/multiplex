# multiplex
The Multiplex

### Dependencies Installation
```
go mod verify 
go mod tidy
go mod download
```

### Deployment of Worker
After setting up the executable and output directory in worker.yml, 
use this command to deploy workers.

```
WORKERS=[number of workers to deploy] make deploy-workers
```
#### Example

```
WORKERS=5 make deploy-workers
```

### Give job to worker
```
go run main.go  -log-name=[log-output-name] -worker-port=[port number of a worker] &
```
#### Example

```
go run main.go  -log-name=test-drive-1 -worker-port=9090 &
```
