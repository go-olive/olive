package platform

import (
	l "github.com/go-olive/olive/src/log"
)

var SharedManager = &Manager{}

type Manager struct {
	ctrlMap map[string]PlatformCtrl
}

func (p *Manager) RegisterCtrl(ctrls ...PlatformCtrl) {
	if p.ctrlMap == nil {
		p.ctrlMap = map[string]PlatformCtrl{}
	}
	for _, ctrl := range ctrls {
		_, ok := p.ctrlMap[ctrl.Type()]
		if ok {
			l.Logger.Error("[%T]Type(%s)已注册\n", p, ctrl.Type())
		}
		p.ctrlMap[ctrl.Type()] = ctrl
	}
}

func (p *Manager) Ctrl(typ string) (PlatformCtrl, bool) {
	v, ok := p.ctrlMap[typ]
	return v, ok
}

type PlatformCtrl interface {
	Type() string
	Name() string
	ParserType() string

	WithRoomOn() Option
	WithStreamURL() Option

	StreamURL(c PlatformCtrl, roomID string) (string, error)
	Snapshot(c PlatformCtrl, roomID string, options ...Option) (*Snapshot, error)
}
