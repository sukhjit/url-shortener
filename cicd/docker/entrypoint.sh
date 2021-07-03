#!/bin/bash

set -euo pipefail

apt-get update

if ! which node &> /dev/null ; then
	echo "Installing nodejs"
	curl -sL https://deb.nodesource.com/setup_12.x | bash -
	apt-get install -y nodejs
fi

echo "Container is up and running..."
tail -f /dev/null
