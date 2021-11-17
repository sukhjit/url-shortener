#!/bin/bash

set -euo pipefail

if ! which CompileDaemon &> /dev/null ; then
    go install github.com/githubnemo/CompileDaemon@v1.4.0
fi

go build -o api main.go
CompileDaemon -build="go build -o api main.go" -command="./api"
