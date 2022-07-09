package recorder

import (
	"errors"
	"sync"
	"time"

	"github.com/go-olive/olive/src/config"
	"github.com/go-olive/olive/src/engine"
)

type Manager struct {
	mu     sync.RWMutex
	savers map[engine.ID]Recorder
	stop   chan struct{}
}

func NewManager() *Manager {
	return &Manager{
		savers: make(map[engine.ID]Recorder),
		stop:   make(chan struct{}),
	}
}

func (m *Manager) Stop() {
	close(m.stop)
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

type Splitter interface {
	Split()
}

func (m *Manager) Split() {
	isValid := false
	for _, r := range m.savers {
		if r.Show().GetSplitRule().IsValid() {
			isValid = true
			break
		}
	}
	if !isValid {
		return
	}

	t := time.NewTicker(time.Second * time.Duration(config.APP.SplitRestSeconds))
	defer t.Stop()

	for {
		select {
		case <-m.stop:
			return
		case <-t.C:
			for _, r := range m.savers {
				if r.Show().SatisfySplitRule(r.StartTime(), r.Out()) {
					r.Show().RestartRecorder()
				}
			}
		}
	}
}
