package upstream

import (
	"time"

	"github.com/twoshark/balanceproxy/src/metrics"

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
				healthyCount, archiveCount := m.ClientCounts()
				metrics.Metrics().HealthyUpstreams.Set(float64(healthyCount))
				metrics.Metrics().ArchiveUpstreams.Set(float64(archiveCount))
			}
		}
	}()
	return quitHealthCheck
}

func (m *Manager) ClientCounts() (int, int) {
	healthyCount := 0
	archiveCount := 0
	for _, client := range m.Clients {
		if client.Healthy() {
			healthyCount++
		}
		if client.IsArchive() {
			archiveCount++
		}
	}
	return healthyCount, archiveCount
}
