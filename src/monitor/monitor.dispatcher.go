package monitor

import (
	"github.com/go-olive/olive/src/dispatcher"
	"github.com/go-olive/olive/src/engine"
	"github.com/go-olive/olive/src/enum"
	l "github.com/go-olive/olive/src/log"
	"github.com/sirupsen/logrus"
)

var MonitorManager = NewManager()

func init() {
	dispatcher.SharedManager.Register(
		MonitorManager,
	)
}

func (m *Manager) Dispatch(event *dispatcher.Event) error {
	show := event.Object.(engine.Show)

	l.Logger.WithFields(logrus.Fields{
		"pf": show.GetPlatform(),
		"id": show.GetRoomID(),
	}).Info("dispatch ", event.Type)

	switch event.Type {
	case enum.EventType.AddMonitor:
		return m.addMonitor(show)
	case enum.EventType.RemoveMonitor:
		return m.removeMonitor(show)
	}
	return nil
}

func (m *Manager) DispatcherType() enum.DispatcherTypeID {
	return enum.DispatcherType.Monitor
}

func (m *Manager) DispatchTypes() []enum.EventTypeID {
	return []enum.EventTypeID{
		enum.EventType.AddMonitor,
		enum.EventType.RemoveMonitor,
	}
}
