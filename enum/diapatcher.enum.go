package enum

type DispatcherTypeID uint32

var DispatcherType = struct {
	Monitor  DispatcherTypeID
	Recorder DispatcherTypeID
}{
	Monitor:  100,
	Recorder: 200,
}
