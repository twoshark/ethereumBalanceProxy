package main

import (
	"github.com/twoshark/ethbalanceproxy/src/cmd"
	"github.com/twoshark/ethbalanceproxy/src/version"
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
