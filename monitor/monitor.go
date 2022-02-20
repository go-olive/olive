package monitor

import (
	"sync/atomic"
	"time"

	"github.com/lthibault/jitterbug/v2"
	"github.com/luxcgo/lifesaver/dispatcher"
	"github.com/luxcgo/lifesaver/engine"
	"github.com/luxcgo/lifesaver/enum"
	"github.com/luxcgo/lifesaver/platform"
)

type Monitor interface {
	Start() error
	Stop()
}

func NewMonitor(show engine.Show) Monitor {
	return &monitor{
		status:   enum.Status.Starting,
		show:     show,
		stop:     make(chan struct{}),
		snapshot: &platform.Snapshot{},
	}
}

type monitor struct {
	status   enum.StatusID
	show     engine.Show
	stop     chan struct{}
	snapshot *platform.Snapshot
}

func (m *monitor) Start() error {
	if !atomic.CompareAndSwapUint32(&m.status, enum.Status.Starting, enum.Status.Pending) {
		return nil
	}
	defer atomic.CompareAndSwapUint32(&m.status, enum.Status.Pending, enum.Status.Running)
	m.refresh()
	go m.run()
	return nil
}

func (m *monitor) Stop() {
	if !atomic.CompareAndSwapUint32(&m.status, enum.Status.Running, enum.Status.Stopping) {
		return
	}
	m.show.RemoveRecorder()
	close(m.stop)
}

func (m *monitor) refresh() {
	latestSnapshot, err := m.show.Snapshot()
	if err != nil {
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

	d, ok := dispatcher.SharedManager.Dispatcher(enum.DispatcherType.Recorder)
	if !ok {
		return
	}
	e := dispatcher.NewEvent(eventType, m.show)
	d.Dispatch(e)
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
			return
		case <-t.C:
			m.refresh()
		}
	}
}
