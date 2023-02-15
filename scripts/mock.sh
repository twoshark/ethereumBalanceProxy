#!/usr/bin/env bash
set -eo pipefail

echo "Remove Old Mocks ..."
rm -rf mocks 2>> /dev/null
echo "Mocking interfaces ..."
mockgen -source src/upstream/ethereum/client.go -destination mocks/mockClient.go
