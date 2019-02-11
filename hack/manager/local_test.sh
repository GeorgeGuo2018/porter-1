#!/bin/bash

set -e

go build -o bin/manager cmd/manager/main.go
./bin/manager -f config/bgp/config.toml