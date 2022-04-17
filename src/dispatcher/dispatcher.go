package dispatcher

import (
	"github.com/go-olive/olive/src/enum"
)

type Dispatcher interface {
	Dispatch(event *Event) error
	DispatcherType() enum.DispatcherTypeID
	DispatchTypes() []enum.EventTypeID
}
