package monitor

import (
	"sync/atomic"
	"time"

	"github.com/lthibault/jitterbug/v2"
	"github.com/luxcgo/lifesaver/dispatcher"
	"github.com/luxcgo/lifesaver/engine"
	"github.com/luxcgo/lifesaver/enum"
	l "github.com/luxcgo/lifesaver/log"
	"github.com/luxcgo/lifesaver/platform"
	"github.com/sirupsen/logrus"
)

type Monitor interface {
	Start() error
	Stop()
	Done() <-chan struct{}
}

func NewMonitor(show engine.Show) Monitor {
	return &monitor{
		status:   enum.Status.Starting,
		show:     show,
		stop:     make(chan struct{}),
		snapshot: &platform.Snapshot{},
		done:     make(chan struct{}),
	}
}

type monitor struct {
	status   enum.StatusID
	show     engine.Show
	stop     chan struct{}
	snapshot *platform.Snapshot
	done     chan struct{}
}

func (m *monitor) Start() error {
	if !atomic.CompareAndSwapUint32(&m.status, enum.Status.Starting, enum.Status.Pending) {
		return nil
	}

	l.Logger.WithFields(logrus.Fields{
		"pf": m.show.GetPlatform(),
		"id": m.show.GetRoomID(),
	}).Info("monitor start")

	defer atomic.CompareAndSwapUint32(&m.status, enum.Status.Pending, enum.Status.Running)
	m.refresh()

	go m.run()

	return nil
}

func (m *monitor) Stop() {
	if !atomic.CompareAndSwapUint32(&m.status, enum.Status.Running, enum.Status.Stopping) {
		return
	}
	close(m.stop)
}

func (m *monitor) refresh() {
	latestSnapshot, err := m.show.Snapshot()
	if err != nil {
		if err.Error() == "not on air" {
			return
		}
		l.Logger.Error(err)
		return
	}
	defer func() {
		m.snapshot = latestSnapshot
	}()
	var eventType enum.EventTypeID
	if !m.snapshot.RoomOn && latestSnapshot.RoomOn {
		eventType = enum.EventType.AddRecorder
	} else if m.snapshot.RoomOn && !latestSnapshot.RoomOn {
		eventType = enum.EventType.RemoveRecorder
	} else {
		return
	}

	l.Logger.WithFields(logrus.Fields{
		"pf":  m.show.GetPlatform(),
		"id":  m.show.GetRoomID(),
		"old": m.snapshot.RoomOn,
		"new": latestSnapshot.RoomOn,
	}).Info("live status changed")

	d, ok := dispatcher.SharedManager.Dispatcher(enum.DispatcherType.Recorder)
	if !ok {
		return
	}
	e := dispatcher.NewEvent(eventType, m.show)
	if err := d.Dispatch(e); err != nil {
		l.Logger.Error(err)
	}

}

func (m *monitor) run() {
	t := jitterbug.New(
		time.Second*15,
		&jitterbug.Norm{Stdev: time.Second * 3},
	)
	defer t.Stop()

	for {
		select {
		case <-m.stop:
			close(m.done)
			l.Logger.WithFields(logrus.Fields{
				"pf": m.show.GetPlatform(),
				"id": m.show.GetRoomID(),
			}).Info("monitor stop")
			return
		case <-t.C:
			m.refresh()
		}
	}
}

func (m *monitor) Done() <-chan struct{} {
	return m.done
}
