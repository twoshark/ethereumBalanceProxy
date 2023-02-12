package upstream

import (
	"github.com/twoshark/balanceproxy/src/metrics"
	"time"

	"github.com/spf13/viper"
)

var quitHealthCheck chan bool

// StartHealthCheck starts a periodic verification of health of each endpoint
// It returns a channel to stop it
func (m *Manager) StartHealthCheck() chan bool {
	period := viper.GetInt("HEALTH_CHECK_PERIOD")
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	quitHealthCheck = make(chan bool)
	go func() {
		for {
			select {
			case <-quitHealthCheck:
				return
			case <-ticker.C:
				for i := range m.Clients {
					m.Clients[i].EvaluatedHealthCheck()
				}
				m.ExportHealthyUpstreamCount()
			}
		}
	}()
	return quitHealthCheck
}

func (m *Manager) ExportHealthyUpstreamCount() {
	count := 0
	for _, client := range m.Clients {
		if client.Healthy() {
			count++
		}
	}
	metrics.Metrics().HealthyUpstreams.Set(float64(count))
}
