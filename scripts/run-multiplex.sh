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

sudo WORKERS=$WORKERS make deploy-workers

# Run Redis Instance
sudo docker pull redis
sudo docker run --name redis-test-instance -p 6379:6379 -d redis

# Run Dispatcher
sudo go run ./job/dispatcher/main.go -workers=$WORKERS &

# Run Mover
sudo go run ./mover/move-worker/main.go &

