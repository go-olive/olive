package monitor

import (
	"errors"
	"sync"

	"github.com/go-olive/olive/src/engine"
)

type Manager struct {
	mu     sync.RWMutex
	savers map[engine.ID]Monitor
}

func NewManager() *Manager {
	return &Manager{
		savers: make(map[engine.ID]Monitor),
	}
}

func (m *Manager) Stop() {
	for _, monitor := range m.savers {
		monitor.Stop()
		<-monitor.Done()
	}
}

func (m *Manager) addMonitor(show engine.Show) error {
	show.RemoveRecorder()

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.savers[show.GetID()]; ok {
		return errors.New("exist")
	}
	monitor := NewMonitor(show)
	m.savers[show.GetID()] = monitor
	return monitor.Start()
}

func (m *Manager) removeMonitor(show engine.Show) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	monitor, ok := m.savers[show.GetID()]
	if !ok {
		return errors.New("monitor not exist")
	}
	monitor.Stop()
	delete(m.savers, show.GetID())
	return nil
}
