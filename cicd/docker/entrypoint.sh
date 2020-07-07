#!/bin/bash

set -e
set -v

# install compiledaemon for hot reload
go get github.com/githubnemo/CompileDaemon

cd /app
CompileDaemon -build="make dev" -command="./url-shortener"
