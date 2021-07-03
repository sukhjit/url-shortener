#!/bin/bash

set -euo pipefail

# # install compiledaemon for hot reload
# go get github.com/githubnemo/CompileDaemon

# cd /app
# CompileDaemon -build="make dev" -command="./url-shortener"

echo "Container is up and running..."
tail -f /dev/null
