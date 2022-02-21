package enum

type EventTypeID uint32

var EventType = struct {
	AddMonitor    EventTypeID
	RemoveMonitor EventTypeID

	AddRecorder    EventTypeID
	RemoveRecorder EventTypeID
}{
	AddMonitor:    101,
	RemoveMonitor: 102,

	AddRecorder:    201,
	RemoveRecorder: 202,
}
