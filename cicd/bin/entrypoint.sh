#!/bin/bash

set -euo pipefail

if ! which npm > /dev/null 2>&1; then
    apt-get update && apt-get install -y nodejs npm
fi

echo ""
echo ""
echo "Setup finished. Container is ready to use."
echo ""
tail -f /dev/null
