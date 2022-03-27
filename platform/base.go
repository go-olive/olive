package platform

type Base struct {
}

func (b *Base) ParserType() string {
	return "ffmpeg"
}

func (b *Base) DefaultOptions(pc PlatformCtrl) []Option {
	return []Option{
		pc.WithRoomOn(),
		pc.WithStreamURL(),
	}
}

func (b *Base) StreamURL(pc PlatformCtrl, roomID string) (string, error) {
	s, err := pc.Snapshot(pc, roomID)
	return s.StreamURL, err
}

func (b *Base) Snapshot(pc PlatformCtrl, roomID string, options ...Option) (*Snapshot, error) {
	s := &Snapshot{
		RoomID: roomID,
	}

	if options == nil {
		options = b.DefaultOptions(pc)
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return s, err
		}
	}

	return s, nil
}
