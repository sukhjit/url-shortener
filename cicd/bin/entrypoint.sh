#!/bin/bash

set -euo pipefail

if ! which npm > /dev/null 2>&1; then
    curl -fsSL https://deb.nodesource.com/setup_16.x | bash
    apt-get update && apt-get install -y nodejs
fi

echo ""
echo ""
echo "Setup finished. Container is ready to use."
echo ""
tail -f /dev/null
