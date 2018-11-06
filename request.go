package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	Timestamp     time.Time           `json:"timestamp"`
	Proto         string              `json:"proto"`
	Method        string              `json:"method"`
	Host          string              `json:"host"`
	Path          string              `json:"path"`
	ContentLength int64               `json:"content_length"`
	Header        map[string][]string `json:"header"`
	Body          []byte              `json:"body"`
	Trailer       map[string][]string `json:"trailer"`
	Query         map[string][]string `json:"query"`
	Form          map[string][]string `json:"form"`
}

func copyMap(src map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for k, v := range src {
		dst[k] = make([]string, len(v))
		copy(dst[k], v)
	}
	return dst
}

func NewRequest(ts time.Time, req *http.Request, body io.ReadCloser) (*Request, error) {
	r := &Request{
		Timestamp:     ts,
		Proto:         req.Proto,
		Method:        req.Method,
		Host:          req.Host,
		Path:          req.URL.Path,
		ContentLength: req.ContentLength,
	}

	if req.Header != nil {
		r.Header = copyMap(req.Header)
	}
	if req.Trailer != nil {
		r.Trailer = copyMap(req.Trailer)
	}
	q := req.URL.Query()
	if len(q) > 0 {
		r.Query = req.URL.Query()
	}
	if body != nil {
		defer body.Close()
		var err error
		r.Body, err = ioutil.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("couldn't read request body: %s", err)
		}
		if ct := req.Header.Get("Content-Type"); ct == "application/x-www-form-urlencoded" {
			r.Form, _ = url.ParseQuery(string(r.Body))
		}
	}

	return r, nil
}
