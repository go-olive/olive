package dispatcher

import (
	"log"

	"github.com/luxcgo/lifesaver/enum"
)

var SharedManager = &Manager{}

type Manager struct {
	savers map[enum.DispatcherTypeID]Dispatcher
}

func (m *Manager) Register(ds ...Dispatcher) {
	if m.savers == nil {
		m.savers = map[enum.DispatcherTypeID]Dispatcher{}
	}

	for _, d := range ds {
		_, ok := m.savers[d.DispatcherType()]
		if ok {
			log.Printf("[%T]Type(%d)已注册\n", d, d.DispatcherType())
		}
		m.savers[d.DispatcherType()] = d
	}
}

func (m *Manager) Dispatcher(typ enum.DispatcherTypeID) (Dispatcher, bool) {
	v, ok := m.savers[typ]
	return v, ok
}
