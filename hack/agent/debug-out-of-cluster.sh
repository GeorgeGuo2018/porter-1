#!/bin/bash

#!/bin/bash

set -e
master=root@192.168.98.2
echo "Building binary"
OOS=linux GOARCH=amd64 go build -o bin/agent cmd/agent/main.go
echo "transport to remote"
scp bin/agent $master:/root/
echo "Starting"
ssh $master "./agent " ##> agent_log.txt 2>&1 &