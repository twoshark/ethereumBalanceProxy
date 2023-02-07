package upstream

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/twoshark/balanceproxy/upstream/ethereum"
)

var (
	errNoUpstreamAvailable     = errors.New("unable to connect to any upstream endpoint")
	errNoHealthyUpstreamClient = errors.New("no healthy upstream client available")
)

var quitHealthCheck chan bool

// Manager maintains health status for Clients and provides Clients to calling code.
type Manager struct {
	Clients []ethereum.IClient
}

func NewManager(endpoints []string) *Manager {
	mgr := new(Manager)
	mgr.Clients = make([]ethereum.IClient, len(endpoints))
	return mgr
}

func (m *Manager) LoadClients(endpoints []string) {
	for i := 0; i < len(endpoints); i++ {
		m.Clients[i] = ethereum.NewClient(endpoints[i])
	}
}

var wg sync.WaitGroup

// ConnectAll attempts to ready all clients for calling.
func (m *Manager) ConnectAll() error {
	clientCount := len(m.Clients)
	wg.Add(clientCount)
	availableChans := make(chan bool, clientCount)
	for i := range m.Clients {
		go func() {
			defer wg.Done()
			availableChans <- m.Connect(m.Clients[i])
		}()
	}
	wg.Wait()
	var anyAvailable bool
	for available := range availableChans {
		if available {
			anyAvailable = true
			break
		}
	}

	if !anyAvailable {
		return errNoUpstreamAvailable
	}

	return nil
}

// Connect dials a Client and if successful, checks its health
// If both succeed the client is marked healthy
// This returns true for a connected and healthy client, otherwise false
func (m *Manager) Connect(client ethereum.IClient) bool {
	if err := client.Dial(); err != nil {
		log.Error("client failed to connect: ", err)
		return false
	}
	err := client.HealthCheck()
	if err != nil {
		log.Error("client failed health check and will not be available for calls until (Manager).Connect() is run again: ", err)
		return false
	}
	return true
}

// GetClient provides the first available healthy ethereum to satisfy a caller's request.
func (m *Manager) GetClient() (ethereum.IClient, error) {
	for i := 0; i < len(m.Clients); i++ {
		if !m.Clients[i].Healthy() {
			go func() {
				err := m.Clients[i].Dial()
				if err != nil {
					log.Error("Redial failed for: ", m.Clients[i])
				}
			}()
			continue
		}
		return m.Clients[i], nil
	}
	return nil, errNoHealthyUpstreamClient
}

func (m *Manager) StartHealthCheck() chan bool {
	failureCount := 0
	failureLimit := viper.GetInt("HEALTH_FAILURE_THRESHOLD")
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
				log.Trace("HealthCheck Heartbeat")
				for _, client := range m.Clients {
					err := client.HealthCheck()
					if err != nil {
						log.Error("UpstreamClient Health Check Failure: ", err)
						failureCount++
						successStreak = 0
						if failureCount > failureLimit {
							client.SetHealth(false)
						}
						continue
					}

					successStreak++
					if !client.Healthy() && successStreak >= successThreshold {
						client.SetHealth(true)
					}
				}
			}
		}
	}()
	return quitHealthCheck
}
