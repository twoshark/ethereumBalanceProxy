package common

import "strings"

// AppConfiguration is a convenience wrapper for values from flags and env vars.
type AppConfiguration struct {
	ListenPort string
	Endpoints  []string
}

func NewAppConfiguration(listenPort string, endpointsFlag *string) AppConfiguration {
	return AppConfiguration{
		ListenPort: listenPort,
		Endpoints:  strings.Split(*endpointsFlag, ","),
	}
}
