package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Timestamp     time.Time           `json:"timestamp"`
	Proto         string              `json:"proto"`
	Status        string              `json:"status"`
	StatusCode    int                 `json:"status_code"`
	ContentLength int64               `json:"content_length"`
	Header        map[string][]string `json:"header"`
	Body          []byte              `json:"body"`
	Trailer       map[string][]string `json:"trailer"`
	Form          map[string][]string `json:"form"`
}

func NewResponse(ts time.Time, resp *http.Response, body io.ReadCloser) (*Response, error) {
	r := &Response{
		Timestamp:     ts,
		Proto:         resp.Proto,
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
	}

	if resp.Header != nil {
		r.Header = copyMap(resp.Header)
	}
	if resp.Trailer != nil {
		r.Trailer = copyMap(resp.Trailer)
	}
	if body != nil {
		defer body.Close()
		var err error
		r.Body, err = ioutil.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("couldn't read response body: %s", err)
		}
		if ct := resp.Header.Get("Content-Type"); ct == "application/x-www-form-urlencoded" {
			r.Form, _ = url.ParseQuery(string(r.Body))
		}
	}

	return r, nil
}
