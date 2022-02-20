package dispatcher

import (
	"github.com/luxcgo/lifesaver/enum"
)

type Dispatcher interface {
	Dispatch(event *Event) error
	DispatcherType() enum.DispatcherTypeID
}
