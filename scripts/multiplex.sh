#!/bin/bash
  
# turn on bash's job control
set -m
  
if [ "$WORKERS" == "" ]
then
	WORKERS=16
fi

# Start workers and put it in background
basePort=9090
i=0
while [ "$i" -lt $WORKERS ]; do
    ./server -port=$(($basePort+$i)) &
    i=$(( i + 1 ))
done 

# Run Dispatcher
./main -workers=$WORKERS