package engine

type Device struct {
	Shows map[ID]*Show
	// Monitor *Monitor
}

func NewDevice() *Device {
	return &Device{}
}
