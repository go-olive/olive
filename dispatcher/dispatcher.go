package dispatcher

import (
	"github.com/go-olive/olive/enum"
)

type Dispatcher interface {
	Dispatch(event *Event) error
	DispatcherType() enum.DispatcherTypeID
	DispatchTypes() []enum.EventTypeID
}
