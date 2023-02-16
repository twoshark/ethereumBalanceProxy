package main

import (
	"github.com/twoshark/balanceproxy/src/cmd"
	"github.com/twoshark/balanceproxy/src/version"
)

var (
	Version        string
	CommitHash     string
	BuildTimestamp string
)

func main() {
	version.Set(Version, CommitHash, BuildTimestamp)
	cmd.Execute()
}
