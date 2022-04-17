package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 20 * time.Second,
	// Transport: &http.Transport{
	// 	//控制主机的最大空闲连接数，0为没有限制
	// 	MaxIdleConns:        1000,
	// 	MaxIdleConnsPerHost: 1000,
	// 	//长连接在关闭之前，保持空闲的最长时间，0表示没限制。
	// 	IdleConnTimeout: 60 * time.Second,
	// },
}

type HttpRequest struct {
	URL          string
	Method       string
	Param        io.Reader
	RequestData  map[string]interface{}
	ResponseData interface{}
	Header       map[string]string
	ContentType  string
}

func (this *HttpRequest) Send() error {
	param, err := this.buildParam()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(this.Method, this.URL, param)
	if err != nil {
		return fmt.Errorf("create http request failed: %s", err.Error())
	}

	req.Header.Set("Content-Type", this.ContentType)
	for k, v := range this.Header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send http request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	switch this.ResponseData.(type) {
	case string:
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		this.ResponseData = string(respBody)
		return nil
	default:
		return json.NewDecoder(resp.Body).Decode(this.ResponseData)
	}
}

func (this *HttpRequest) buildParam() (io.Reader, error) {
	switch this.ContentType {
	case "application/form-data":
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		for k, v := range this.RequestData {
			if err := writer.WriteField(k, fmt.Sprint(v)); err != nil {
				return nil, fmt.Errorf("write form field %s:%s failed: %s", k, v, err.Error())
			}
		}
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("close multipart form writer failed: %s", err.Error())
		}
		this.ContentType = writer.FormDataContentType()
		return body, nil
	case "application/x-www-form-urlencoded":
		body := make(url.Values)
		for k, v := range this.RequestData {
			body[k] = []string{fmt.Sprint(v)}
		}
		paramData := strings.NewReader(body.Encode())
		return paramData, nil
	case "application/json":
		jsonParamBytes, err := json.Marshal(this.RequestData)
		if err != nil {
			return nil, fmt.Errorf("Marshal json error:%s ", err.Error())
		}
		paramData := bytes.NewBuffer(jsonParamBytes)
		return paramData, nil
	default:
		return nil, fmt.Errorf("ContentType = %s not supported", this.ContentType)
	}
}
