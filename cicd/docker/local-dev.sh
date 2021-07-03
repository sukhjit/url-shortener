#!/bin/bash

set -euo pipefail

if ! which CompileDaemon &> /dev/null ; then
	echo "Installing CompileDaemon"
	go install github.com/githubnemo/CompileDaemon
fi

CompileDaemon -build="go build -o url-shortener main.go" -command="./url-shortener"
