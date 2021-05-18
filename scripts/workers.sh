#!/bin/bash
  
# turn on bash's job control
set -m
  
if [ "$WORKERS" == "" ]
then
	WORKERS=1
fi

# Start workers and put it in background
basePort=9090
i=0
while [ "$i" -lt $WORKERS ]; do
    go run ./worker/cmd/server.go -port=$(($basePort+$i)) &
    i=$(( i + 1 ))
done 
