package util

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetURLContent(url string) (string, error) {
	webUserAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", webUserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	content := string(respBody)
	return content, nil
}

func Match(pattern, content string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	submatch := re.FindAllStringSubmatch(content, -1)
	res := make([]string, 0)
	for _, v := range submatch {
		res = append(res, string(v[1]))
	}
	if len(res) < 1 {
		return "", errors.New("pattern not found")
	}
	return res[0], nil
}
