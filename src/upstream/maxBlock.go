package upstream

import (
	"time"

	"github.com/twoshark/ethbalanceproxy/src/metrics"
)

func (m *Manager) SetMaxBlock(block uint64) {
	m.blockLock.Lock()
	defer m.blockLock.Unlock()
	m.maxBlockObserved = block
	metrics.Metrics().MaxBlock.Set(float64(block))
}

func (m *Manager) GetMaxBlock() uint64 {
	m.blockLock.Lock()
	defer m.blockLock.Unlock()
	return m.maxBlockObserved
}

func (m *Manager) StartBlockWatcher() chan bool {
	m.SetMaxBlock(0)
	var quit chan bool
	go func() {
		t := time.NewTicker(12 * time.Second)
		for {
			select {
			case <-quit:
				return
			case <-t.C:
				m.CheckClientMaxBlocks()
			}
		}
	}()
	return quit
}

func (m *Manager) CheckClientMaxBlocks() {
	oldMax := m.GetMaxBlock()
	newMax := oldMax
	for _, client := range m.Clients {
		clientMax := client.GetMaxBlock()
		if clientMax > newMax {
			newMax = clientMax
		}
	}
	if newMax > oldMax {
		m.SetMaxBlock(newMax)
	}
}
