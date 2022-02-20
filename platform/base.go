package platform

type Base struct {
}

func (b *Base) ParserType() string {
	return "ffmpeg"
}
