package dispatcher

import (
	"github.com/luxcgo/lifesaver/enum"
)

type Event struct {
	Type   enum.EventTypeID
	Object interface{}
}

func NewEvent(typ enum.EventTypeID, object interface{}) *Event {
	return &Event{
		Type:   typ,
		Object: object,
	}
}
