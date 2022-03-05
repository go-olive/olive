package huya

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/luxcgo/lifesaver/platform"
)

const userAgent = "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko); Chrome/75.0.3770.100 Mobile Safari/537.36 "

func init() {
	platform.SharedManager.RegisterCtrl(
		new(huyaCtrl),
	)
}

type huyaCtrl struct {
	platform.Base
}

func (c *huyaCtrl) Type() string {
	return "huya"
}

func (c *huyaCtrl) Name() string {
	return "虎牙"
}

func (c *huyaCtrl) StreamURL(roomID string) (string, error) {
	snapshot, err := c.Snapshot(roomID)
	if err != nil {
		return "", err
	}
	if snapshot.RoomOn {
		return snapshot.StreamURL, nil
	}
	return "", errors.New("not on air")
}

func (c *huyaCtrl) streamURL(roomID string) (string, error) {
	roomURL := fmt.Sprintf("https://m.huya.com/%s", roomID)
	req, err := http.NewRequest("GET", roomURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	respBody := string(content)
	re := regexp.MustCompile(`liveLineUrl":"([^"]+)",`)
	submatch := re.FindAllStringSubmatch(respBody, -1)
	res := make([]string, 0)
	for _, v := range submatch {
		res = append(res, string(v[1]))
	}
	if len(res) < 1 {
		return "", errors.New("find stream url failed")
	}
	a, _ := base64.RawStdEncoding.DecodeString(res[0])
	return c.proc(string(a)), nil
}

func (*huyaCtrl) proc(in string) string {
	ib := strings.Split(in, "?")
	i, b := ib[0], ib[1]
	r := strings.Split(i, "/")
	s := strings.ReplaceAll(r[len(r)-1], ".flv", "")
	s = strings.ReplaceAll(s, ".m3u8", "")
	c := strings.SplitN(b, "&", 4)
	y := c[len(c)-1]
	n := make(map[string]string)
	for _, v := range c {
		if v == "" {
			continue
		}
		vs := strings.SplitN(v, "=", 2)
		n[vs[0]] = vs[1]
	}
	fm := url.PathEscape(n["fm"])
	ub, _ := base64.RawStdEncoding.DecodeString(fm)
	u := string(ub)
	p := strings.Split(u, "_")[0]
	f := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	l := n["wsTime"]
	t := "0"
	h := strings.Join([]string{p, t, s, f, l}, "_")
	m := md5.New()
	io.WriteString(m, h)
	url := fmt.Sprintf("%s?wsSecret=%x&wsTime=%s&u=%s&seqid=%s&%s", i, m.Sum(nil), l, t, f, y)
	url = "https:" + url
	url = strings.ReplaceAll(url, "hls", "flv")
	url = strings.ReplaceAll(url, "m3u8", "flv")
	return url
}

func (c *huyaCtrl) Snapshot(roomID string) (*platform.Snapshot, error) {
	webUserAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
	roomURL := fmt.Sprintf("https://www.huya.com/%s", roomID)
	req, err := http.NewRequest("GET", roomURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", webUserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	roomOn := strings.Contains(string(content), `"isOn":true`)
	snapshot := &platform.Snapshot{
		RoomOn: roomOn,
	}
	if snapshot.RoomOn {
		snapshot.StreamURL, err = c.streamURL(roomID)
		return snapshot, err
	}
	return snapshot, nil
}
