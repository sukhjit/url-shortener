#!/bin/bash

set -euo pipefail

if ! which -v CompileDaemon &> /dev/null ; then
	echo "Installing CompileDaemon"
	go get github.com/githubnemo/CompileDaemon
fi

CompileDaemon -build="go build -o url-shortener main.go" -command="./url-shortener"
