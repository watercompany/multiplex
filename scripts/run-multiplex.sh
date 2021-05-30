#!/bin/bash

# turn on bash's job control
set -m

# Requirements
# 1. Docker must be running.

# Deploy workers
if [ "$WORKERS" == "" ]
then
	WORKERS=1
fi

WORKERS=$WORKERS make deploy-workers

# Run Redis Instance
docker pull redis
docker run --name redis-test-instance -p 6379:6379 -d redis

# Run Dispatcher
go run ./job/dispatcher/main.go -workers=$WORKERS &

