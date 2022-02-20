package recorder

import (
	"errors"
	"sync"

	"github.com/luxcgo/lifesaver/engine"
)

type Manager struct {
	mu     sync.RWMutex
	savers map[engine.ID]Recorder
}

func NewManager() *Manager {
	return &Manager{
		savers: make(map[engine.ID]Recorder),
	}
}

func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, recorder := range m.savers {
		recorder.Stop()
		delete(m.savers, id)
	}
}

func (m *Manager) addRecorder(show engine.Show) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.savers[show.GetID()]; ok {
		return errors.New("exist")
	}
	recorder, err := NewRecorder(show)
	if err != nil {
		return err
	}
	m.savers[show.GetID()] = recorder
	return recorder.Start()
}

func (m *Manager) removeRecorder(show engine.Show) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	recorder, ok := m.savers[show.GetID()]
	if !ok {
		return errors.New("recorder not exist")
	}
	recorder.Stop()
	delete(m.savers, show.GetID())
	return nil
}
