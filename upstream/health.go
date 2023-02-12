package upstream

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var quitHealthCheck chan bool

// StartHealthCheck starts a periodic verification of health of each endpoint
// It returns a channel to stop it
func (m *Manager) StartHealthCheck() chan bool {
	failureCount := 0
	failureLimit := viper.GetInt("HEALTH_FAILURE_THRESHOLD")
	failureForgive := viper.GetInt("HEALTH_FAILURE_FORGIVENESS_THRESHOLD")
	period := viper.GetInt("HEALTH_CHECK_PERIOD")
	successStreak := 0
	successThreshold := viper.GetInt("HEALTH_SUCCESS_THRESHOLD")
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	quitHealthCheck = make(chan bool)
	go func() {
		for {
			select {
			case <-quitHealthCheck:
				return
			case <-ticker.C:
				for i, client := range m.Clients {
					err := client.HealthCheck()
					if err != nil {
						log.WithFields(log.Fields{
							"upstream": m.endpoints[i],
						}).Error("health check failed: ", err)
						failureCount++
						successStreak = 0
						if failureCount > failureLimit {
							client.SetHealth(false)
							log.WithFields(log.Fields{
								"upstream": m.endpoints[i],
							}).Error("upstream has exceeded failure threshold and is now marked unhealthy")
						}
						continue
					}
					successStreak++
					if !client.Healthy() && successStreak >= successThreshold {
						failureCount = 0
						client.SetHealth(true)
						log.WithFields(log.Fields{
							"upstream": m.endpoints[i],
						}).Error("upstream ")
					}
					if successStreak > failureForgive {
						failureCount = 0
					}
				}
			}
		}
	}()
	return quitHealthCheck
}
