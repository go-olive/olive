package engine

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/go-olive/olive/src/dispatcher"
	"github.com/go-olive/olive/src/enum"
	"github.com/go-olive/olive/src/parser"
	"github.com/go-olive/olive/src/platform"
)

type ID string

type Show interface {
	GetID() ID
	GetPlatform() string
	GetRoomID() string
	GetStreamerName() string
	StreamURL() (string, error)
	Snapshot() (*platform.Snapshot, error)

	AddMonitor() error
	RemoveMonitor() error
	AddRecorder() error
	RemoveRecorder() error

	NewParser() (parser.Parser, error)
}

type show struct {
	ID           ID
	Platform     string
	RoomID       string
	StreamerName string
	enum.ShowTaskStatusID
	stop chan struct{}
	ctrl platform.PlatformCtrl
}

func NewShow(platformType, roomID, streamerName string) (Show, error) {
	pc, valid := platform.SharedManager.Ctrl(platformType)
	if !valid {
		return nil, errors.New("platform not exist")
	}

	s := &show{
		Platform:     platformType,
		RoomID:       roomID,
		StreamerName: streamerName,
		stop:         make(chan struct{}),

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

func (s *show) GetStreamerName() string {
	return s.StreamerName
}

func (s *show) GetPlatform() string {
	return s.Platform
}

func (s *show) genID() ID {
	h := md5.New()
	b := []byte(fmt.Sprintf("%s%s", s.Platform, s.RoomID))
	h.Write(b)
	return ID(hex.EncodeToString(h.Sum(nil)))
}

func (s *show) StreamURL() (string, error) {
	return s.ctrl.StreamURL(s.ctrl, s.RoomID)
}

func (s *show) Snapshot() (*platform.Snapshot, error) {
	return s.ctrl.Snapshot(s.ctrl, s.RoomID)
}

func (s *show) NewParser() (parser.Parser, error) {
	v, ok := parser.SharedManager.Parser(s.ctrl.ParserType())
	if !ok {
		return nil, errors.New("parser not exist")
	}
	return v.New(), nil
}

func (s *show) AddMonitor() error {
	e := dispatcher.NewEvent(enum.EventType.AddMonitor, s)
	return dispatcher.SharedManager.Dispatch(e)
}

func (s *show) RemoveMonitor() error {
	e := dispatcher.NewEvent(enum.EventType.RemoveMonitor, s)
	return dispatcher.SharedManager.Dispatch(e)
}

func (s *show) AddRecorder() error {
	e := dispatcher.NewEvent(enum.EventType.AddRecorder, s)
	return dispatcher.SharedManager.Dispatch(e)
}

func (s *show) RemoveRecorder() error {
	e := dispatcher.NewEvent(enum.EventType.RemoveRecorder, s)
	return dispatcher.SharedManager.Dispatch(e)
}

func (s *show) Stop() {
	dispatcher.SharedManager.Dispatch(dispatcher.NewEvent(enum.EventType.RemoveMonitor, s))
	dispatcher.SharedManager.Dispatch(dispatcher.NewEvent(enum.EventType.RemoveRecorder, s))
}
