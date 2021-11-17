#!/bin/bash

set -euo pipefail

if ! which golangci-lint > /dev/null 2>&1; then
    curl -fsSL https://deb.nodesource.com/setup_14.x | bash
    apt-get update && apt-get install -y npm
fi

echo ""
echo ""
echo "Setup finished. Container is ready to use."
echo ""
tail -f /dev/null
