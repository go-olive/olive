package recorder

import (
	"github.com/luxcgo/lifesaver/dispatcher"
	"github.com/luxcgo/lifesaver/engine"
	"github.com/luxcgo/lifesaver/enum"
)

var RecorderManager = NewManager()

func init() {
	dispatcher.SharedManager.Register(
		RecorderManager,
	)
}

func (m *Manager) Dispatch(event *dispatcher.Event) error {
	show := event.Object.(engine.Show)
	switch event.Type {
	case enum.EventType.AddRecorder:
		return m.addRecorder(show)
	case enum.EventType.RemoveRecorder:
		return m.removeRecorder(show)
	}
	return nil
}

func (m *Manager) DispatcherType() enum.DispatcherTypeID {
	return enum.DispatcherType.Recorder
}
