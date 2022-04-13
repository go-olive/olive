package internal

import (
	_ "github.com/go-olive/olive/platform/bilibili"
	_ "github.com/go-olive/olive/platform/douyin"
	_ "github.com/go-olive/olive/platform/huya"
	_ "github.com/go-olive/olive/platform/youtube"

	_ "github.com/go-olive/olive/monitor"

	_ "github.com/go-olive/olive/parser/ffmpeg"
	_ "github.com/go-olive/olive/parser/streamlink"
)
