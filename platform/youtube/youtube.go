package youtube

import (
	"fmt"
	"log"
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

func (c *youtubeCtrl) StreamURL(roomID string) (string, error) {
	// youtube possibly have multiple lives in one channel,
	// curruently the program returns the first one.
	s, err := c.Snapshot(roomID)
	if err != nil {
		return "", err
	}
	return s.StreamURL, nil
}

func (c *youtubeCtrl) Snapshot(roomID string) (*platform.Snapshot, error) {
	snapShot := &platform.Snapshot{}
	channelURL := fmt.Sprintf("https://www.youtube.com/channel/%s", roomID)
	content, err := util.GetURLContent(channelURL)
	if err != nil {
		return nil, err
	}
	snapShot.RoomOn = strings.Contains(content, `icon":{"iconType":"LIVE"}}`)
	if !snapShot.RoomOn {
		return snapShot, nil
	}

	streamID, err := util.Match(`"videoRenderer":{"videoId":"([^"]+)",`, content)
	if err != nil {
		return nil, err
	}

	roomURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", streamID)
	snapShot.StreamURL = roomURL

	roomContent, err := util.GetURLContent(roomURL)
	if err != nil {
		log.Println(err)
		return snapShot, nil
	}

	title, err := util.Match(`name="title" content="([^"]+)"`, roomContent)
	if err != nil {
		log.Println(err)
		return snapShot, nil
	}
	snapShot.RoomName = title
	return snapShot, nil
}

func (c *youtubeCtrl) ParserType() string {
	return "streamlink"
}
