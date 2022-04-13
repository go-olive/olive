package recorder

import (
	"errors"
	"sync"

	"github.com/go-olive/olive/engine"
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
	for _, recorder := range m.savers {
		recorder.Stop()
		<-recorder.Done()
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
