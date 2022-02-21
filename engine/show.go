package engine

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/luxcgo/lifesaver/dispatcher"
	"github.com/luxcgo/lifesaver/enum"
	"github.com/luxcgo/lifesaver/parser"
	"github.com/luxcgo/lifesaver/platform"
)

var Shows = []*show{
	{
		Platform: "huya",
		RoomID:   "291252",
	},
}

type ID string

type Show interface {
	GetID() ID
	GetRoomID() string
	StreamURL() (string, error)
	Snapshot() (*platform.Snapshot, error)

	AddMonitor() error
	RemoveMonitor() error
	AddRecorder() error
	RemoveRecorder() error

	NewParser() (parser.Parser, error)
}

type show struct {
	ID       ID
	Platform string
	RoomID   string
	enum.ShowTaskStatusID
	stop chan struct{}
	ctrl platform.PlatformCtrl
}

func NewShow(platformType, roomID string) (Show, error) {
	pc, valid := platform.SharedManager.Ctrl(platformType)
	if !valid {
		return nil, errors.New("not exist")
	}

	s := &show{
		Platform: platformType,
		RoomID:   roomID,
		stop:     make(chan struct{}),

		ctrl: pc,
	}
	s.ID = s.genID()
	return s, nil
}

func (s *show) GetID() ID {
	return s.ID
}

func (s *show) GetRoomID() string {
	return s.RoomID
}

func (s *show) genID() ID {
	h := md5.New()
	b := []byte(fmt.Sprintf("%s%s", s.Platform, s.RoomID))
	h.Write(b)
	return ID(hex.EncodeToString(h.Sum(nil)))
}

func (s *show) StreamURL() (string, error) {
	return s.ctrl.StreamURL(s.RoomID)
}

func (s *show) Snapshot() (*platform.Snapshot, error) {
	return s.ctrl.Snapshot(s.RoomID)
}

func (s *show) NewParser() (parser.Parser, error) {
	v, ok := parser.SharedManager.Parser(s.ctrl.ParserType())
	if !ok {
		return nil, errors.New("parser not exist")
	}
	return v.New(), nil
}

func (s *show) AddMonitor() error {
	d, ok := dispatcher.SharedManager.Dispatcher(enum.DispatcherType.Monitor)
	if !ok {
		return errors.New("internal error")
	}
	e := dispatcher.NewEvent(enum.EventType.AddMonitor, s)
	return d.Dispatch(e)
}

func (s *show) RemoveMonitor() error {
	d, ok := dispatcher.SharedManager.Dispatcher(enum.DispatcherType.Monitor)
	if !ok {
		return errors.New("internal error")
	}
	e := dispatcher.NewEvent(enum.EventType.RemoveMonitor, s)
	return d.Dispatch(e)
}

func (s *show) AddRecorder() error {
	d, ok := dispatcher.SharedManager.Dispatcher(enum.DispatcherType.Recorder)
	if !ok {
		return errors.New("internal error")
	}
	e := dispatcher.NewEvent(enum.EventType.AddRecorder, s)
	return d.Dispatch(e)
}

func (s *show) RemoveRecorder() error {
	d, ok := dispatcher.SharedManager.Dispatcher(enum.DispatcherType.Recorder)
	if !ok {
		return errors.New("internal error")
	}
	e := dispatcher.NewEvent(enum.EventType.RemoveRecorder, s)
	return d.Dispatch(e)
}
