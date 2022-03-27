package youtube

import (
	"errors"
	"fmt"
	"strings"

	"github.com/luxcgo/lifesaver/platform"
	"github.com/luxcgo/lifesaver/util"
)

func init() {
	platform.SharedManager.RegisterCtrl(
		new(youtubeCtrl),
	)
}

type youtubeCtrl struct {
	platform.Base
}

func (c *youtubeCtrl) Type() string {
	return "youtube"
}

func (c *youtubeCtrl) Name() string {
	return "油管"
}

func (c *youtubeCtrl) ParserType() string {
	return "streamlink"
}

// func (c *youtubeCtrl) WithRoomOnBak() platform.Option {
// 	return func(s *platform.Snapshot) error {
// 		webUserAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
// 		channelURL := fmt.Sprintf("https://www.youtube.com/channel/%s", s.RoomID)
// 		req := &util.HttpRequest{
// 			URL:          channelURL,
// 			Method:       "GET",
// 			ResponseData: *new(string),
// 			ContentType:  "application/x-www-form-urlencoded",
// 			Header: map[string]string{
// 				"User-Agent": webUserAgent,
// 			},
// 		}
// 		if err := req.Send(); err != nil {
// 			return err
// 		}
// 		content := fmt.Sprint(req.ResponseData)
// 		s.RoomOn = strings.Contains(content, `icon":{"iconType":"LIVE"}}`)
// 		if s.RoomOn {
// 			streamID, err := util.Match(`"videoRenderer":{"videoId":"([^"]+)",`, content)
// 			if err != nil {
// 				return err
// 			}
// 			s.RoomID = streamID
// 		}
// 		return nil
// 	}
// }

func (c *youtubeCtrl) WithRoomOn() platform.Option {
	return func(s *platform.Snapshot) error {
		channelURL := fmt.Sprintf("https://www.youtube.com/channel/%s", s.RoomID)
		content, err := util.GetURLContent(channelURL)
		if err != nil {
			return err
		}
		s.RoomOn = strings.Contains(content, `icon":{"iconType":"LIVE"}}`)

		streamID, err := util.Match(`"videoRenderer":{"videoId":"([^"]+)",`, content)
		if err != nil {
			return err
		}
		s.RoomID = streamID
		return nil
	}
}

func (c *youtubeCtrl) WithStreamURL() platform.Option {
	return func(s *platform.Snapshot) error {
		if !s.RoomOn {
			return errors.New("not on air")
		}
		// youtube possibly have multiple lives in one channel,
		// curruently the program returns the first one.
		roomURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", s.RoomID)
		s.StreamURL = roomURL
		roomContent, err := util.GetURLContent(roomURL)
		if err != nil {
			return err
		}
		title, err := util.Match(`name="title" content="([^"]+)"`, roomContent)
		if err != nil {
			return err
		}
		s.RoomName = title
		return nil
	}
}
