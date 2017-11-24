package main

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type LineApi struct {
	Authorization string
}

func NewLineApi(authorization string) *LineApi {
	l := &LineApi{
		Authorization: authorization,
	}
	return l
}

func (l *LineApi) sendNotify(text string) (response interface{}, err error) {
	buf := bytes.Buffer{}
	mw := multipart.NewWriter(&buf)
	mw.WriteField("message", text)
	mw.Close()

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Add("Authorization", "Bearer "+l.Authorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response = string(body)
	resp.Body.Close()

	return response, nil
}
