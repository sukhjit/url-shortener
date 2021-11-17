#!/bin/bash

set -euo pipefail

readonly LINT_VERSION="v1.43.0"

if ! which golangci-lint > /dev/null 2>&1; then
    echo "Installing golangci-lint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${LINT_VERSION}
fi

echo "Running golangci-lint"
golangci-lint run ./...
