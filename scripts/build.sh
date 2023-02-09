#!/bin/sh

go mod vendor
go build -o "${BINARY_NAME}" main.go