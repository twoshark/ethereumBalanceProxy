package upstream

import (
	"errors"
	"github.com/twoshark/balanceproxy/src/upstream/ethereum"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	errNoUpstreamAvailable     = errors.New("unable to connect to any upstream endpoint")
	errNoHealthyUpstreamClient = errors.New("no healthy upstream client available")
)

// Manager maintains health status for Clients and provides Clients to calling code.
type Manager struct {
	Clients   []ethereum.IClient
	endpoints []string
}

func NewManager(endpoints []string) *Manager {
	mgr := new(Manager)
	mgr.endpoints = endpoints
	mgr.Clients = make([]ethereum.IClient, len(endpoints))
	return mgr
}

// LoadClients instantiates an ethereum.Client in m.Clients for each provided endpoint
func (m *Manager) LoadClients() {
	for i := 0; i < len(m.endpoints); i++ {
		m.Clients[i] = ethereum.NewClient(m.endpoints[i])
	}
}

var wg sync.WaitGroup

// ConnectAll attempts to ready all clients for calling.
func (m *Manager) ConnectAll() error {
	clientCount := len(m.Clients)
	wg.Add(clientCount)
	availableChans := make(chan bool, clientCount)
	for i := range m.Clients {
		index := i // to quiet linter
		go func() {
			defer wg.Done()
			availableChans <- m.Connect(index)
		}()
	}
	wg.Wait()
	close(availableChans)
	var anyAvailable bool
	for available := range availableChans {
		if available {
			anyAvailable = true
			break
		}
	}
	m.ExportHealthyUpstreamCount()

	if !anyAvailable {
		return errNoUpstreamAvailable
	}

	return nil
}

// Connect dials a Client and if successful, checks its health and sets it in the client
// If both succeed the client is marked healthy
// This returns true for a connected and healthy client, otherwise false
func (m *Manager) Connect(index int) bool {
	if err := m.Clients[index].Dial(); err != nil {
		log.Error("client failed to connect: ", err)
		return false
	}

	if err := m.Clients[index].HealthCheck(); err != nil {
		log.Error("client failed health check and will not be available for calls until (Manager).Connect() is run again: ", err)
		return false
	}
	m.Clients[index].SetHealth(true)

	return true
}

// GetClient provides the first available healthy ethereum to satisfy a caller's request.
func (m *Manager) GetClient() (ethereum.IClient, error) {
	for i := 0; i < len(m.Clients); i++ {
		if !m.Clients[i].Healthy() {
			go func(i int) {
				err := m.Clients[i].Dial()
				if err != nil {
					log.Error("Redial failed for: ", m.Clients[i])
				}
			}(i)
			continue
		}
		return m.Clients[i], nil
	}
	return nil, errNoHealthyUpstreamClient
}
