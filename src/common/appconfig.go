package common

import (
	"strings"

	"github.com/twoshark/ethbalanceproxy/src/metrics"
)

// AppConfiguration is a convenience wrapper for values from flags and env vars.
type AppConfiguration struct {
	ListenPort string
	Endpoints  []string
}

func NewAppConfiguration(listenPort string, endpointsFlag *string) AppConfiguration {
	endpoints := strings.Split(*endpointsFlag, ",")
	metrics.Metrics().ConfiguredUpstreams.Set(float64(len(endpoints)))
	return AppConfiguration{
		ListenPort: listenPort,
		Endpoints:  endpoints,
	}
}
