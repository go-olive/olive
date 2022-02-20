package monitor

import (
	"github.com/luxcgo/lifesaver/dispatcher"
	"github.com/luxcgo/lifesaver/engine"
	"github.com/luxcgo/lifesaver/enum"
)

var MonitorManager = NewManager()

func init() {
	dispatcher.SharedManager.Register(
		MonitorManager,
	)
}

func (m *Manager) Dispatch(event *dispatcher.Event) error {
	show := event.Object.(engine.Show)
	switch event.Type {
	case enum.EventType.AddMonitor:
		return m.addMonitor(show)
	}
	return nil
}

func (m *Manager) DispatcherType() enum.DispatcherTypeID {
	return enum.DispatcherType.Monitor
}
