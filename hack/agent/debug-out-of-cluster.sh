#!/bin/bash

#!/bin/bash

set -e
master=root@192.168.98.2
OOS=linux GOARCH=amd64 go build -o bin/agent cmd/agent/main.go
scp bin/agent $master:/root/
ssh $master "./agent > agent_log.txt 2>&1"