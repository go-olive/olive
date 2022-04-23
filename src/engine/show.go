package engine

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-olive/olive/src/config"
	"github.com/go-olive/olive/src/dispatcher"
	"github.com/go-olive/olive/src/enum"
	"github.com/go-olive/olive/src/parser"

	"github.com/go-olive/tv"
)

type ID string

type Show interface {
	GetID() ID
	GetPlatform() string
	GetRoomID() string
	GetStreamerName() string
	GetOutTmpl() string

	AddMonitor() error
	RemoveMonitor() error
	AddRecorder() error
	RemoveRecorder() error

	NewParser() (parser.Parser, error)

	tv.ITv
}

type show struct {
	ID       ID
	Platform string
	RoomID   string
	Streamer string
	OutTmpl  string
	Parser   string
	enum.ShowTaskStatusID
	stop chan struct{}

	*tv.Tv
}

type ShowOption func(*show)

func WithStreamerName(name string) ShowOption {
	return func(s *show) {
		s.Streamer = name
	}
}

func WithOutTmpl(tmpl string) ShowOption {
	return func(s *show) {
		s.OutTmpl = tmpl
	}
}

func WithParser(parser string) ShowOption {
	return func(s *show) {
		s.Parser = parser
	}
}

func NewShow(platformType, roomID string, opts ...ShowOption) (Show, error) {
	var cookie string
	if platformType == "douyin" {
		cookie = config.APP.PlatformConfig.DouyinCookie
	}
	t, err := tv.New(platformType, roomID, tv.SetCookie(cookie))
	if err != nil {
		return nil, fmt.Errorf("Show init failed! err msg: %s", err.Error())
	}

	s := &show{
		Platform: platformType,
		RoomID:   roomID,

		stop: make(chan struct{}),

		Tv: t,
	}
	for _, opt := range opts {
		opt(s)
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
	return s.Streamer
}

func (s *show) GetPlatform() string {
	return s.Platform
}

func (s *show) GetOutTmpl() string {
	return s.OutTmpl
}

func (s *show) GetParser() string {
	return s.Parser
}

func (s *show) genID() ID {
	h := md5.New()
	b := []byte(fmt.Sprintf("%s%s%d", s.Platform, s.RoomID, time.Now().UnixNano()))
	h.Write(b)
	return ID(hex.EncodeToString(h.Sum(nil)))
}

const (
	defaultTyp = "flv"
)

func (s *show) NewParser() (parser.Parser, error) {
	typ := defaultTyp
	if s.SiteID == "youtube" {
		typ = "streamlink"
	}
	if s.GetParser() != "" {
		typ = s.GetParser()
	}

	v, ok := parser.SharedManager.Parser(typ)
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
