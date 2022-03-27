package platform

type Snapshot struct {
	RoomID    string
	RoomName  string
	RoomOn    bool
	StreamURL string
}

type Option func(*Snapshot) error
