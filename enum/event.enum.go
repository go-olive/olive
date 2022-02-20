package enum

type EventTypeID uint32

var EventType = struct {
	AddMonitor EventTypeID

	AddRecorder    EventTypeID
	RemoveRecorder EventTypeID
}{
	AddMonitor: 101,

	AddRecorder:    201,
	RemoveRecorder: 202,
}
