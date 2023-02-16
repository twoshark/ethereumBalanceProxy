#!/usr/bin/env bash
BINARY_NAME=ethBalanceProxy
VERSION="$(cat src/version.txt)"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

LDFLAGS=(
  "-X 'main.Version=${VERSION}'"
  "-X 'main.CommitHash=${COMMIT_HASH}'"
  "-X 'main.BuildTime=${BUILD_TIMESTAMP}'"
)

go mod vendor
go build -ldflags="${LDFLAGS[*]}" -o "${BINARY_NAME}" main.go
./ethBalanceProxy version