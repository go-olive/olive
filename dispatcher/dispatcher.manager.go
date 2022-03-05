package dispatcher

import (
	"errors"

	"github.com/luxcgo/lifesaver/enum"
	l "github.com/luxcgo/lifesaver/log"
)

var SharedManager = &Manager{}

type Manager struct {
	savers             map[enum.DispatcherTypeID]Dispatcher
	dispatchFuncSavers map[enum.EventTypeID]Dispatcher
}

func (m *Manager) Register(ds ...Dispatcher) {
	if m.savers == nil {
		m.savers = map[enum.DispatcherTypeID]Dispatcher{}
	}

	if m.dispatchFuncSavers == nil {
		m.dispatchFuncSavers = map[enum.EventTypeID]Dispatcher{}
	}

	for _, d := range ds {
		_, ok := m.savers[d.DispatcherType()]
		if ok {
			l.Logger.WithField("type", d).Warn("dispatcher 已注册")
		}
		m.savers[d.DispatcherType()] = d

		for _, v := range d.DispatchTypes() {
			m.RegisterFunc(v, d)
		}
	}
}

func (m *Manager) RegisterFunc(typ enum.EventTypeID, d Dispatcher) {
	if m.dispatchFuncSavers == nil {
		m.dispatchFuncSavers = map[enum.EventTypeID]Dispatcher{}
	}
	_, ok := m.dispatchFuncSavers[typ]
	if ok {
		l.Logger.WithField("type", typ.String()).Warn("dipatch func 已注册")
	}
	m.dispatchFuncSavers[typ] = d
}

func (m *Manager) Dispatcher(typ enum.DispatcherTypeID) (Dispatcher, bool) {
	v, ok := m.savers[typ]
	return v, ok
}

func (m *Manager) Dispatch(e *Event) error {
	dispatchFunc, ok := m.dispatchFuncSavers[e.Type]
	if !ok {
		return errors.New("dispatch func not exist")
	}

	return dispatchFunc.Dispatch(e)
}
