package main

import (
	"bufio"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func newTestHTTPRequest() *http.Request {
	requestString := "POST / HTTP/1.1\r\n" +
		"Accept: application/json, */*\r\n" +
		"Accept-Encoding: gzip, deflate\r\n" +
		"Connection: keep-alive\r\n" +
		"Content-Length: 27\r\n" +
		"Content-Type: application/json\r\n" +
		"Host: example.com\r\n" +
		"\r\n" +
		"{\"message\": \"Hello World!\"}"
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(requestString)))
	if err != nil {
		panic(err.Error())
	}
	return req
}

func newTestRequest(ts time.Time) *Request {
	r := &Request{
		Timestamp:     ts,
		Proto:         "HTTP/1.1",
		Method:        http.MethodPost,
		Host:          "example.com",
		Path:          "/",
		ContentLength: 27,
		Header: map[string][]string{
			"Accept":          {"application/json, */*"},
			"Accept-Encoding": {"gzip, deflate"},
			"Connection":      {"keep-alive"},
			"Content-Length":  {"27"},
			"Content-Type":    {"application/json"},
		},
		Body:    []byte(`{"message": "Hello World!"}`),
		Trailer: nil,
		Query:   nil,
		Form:    nil,
	}
	return r
}

func TestNewRequest(t *testing.T) {
	ts := time.Now()
	req := newTestHTTPRequest()
	r, err := NewRequest(ts, req, req.Body)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}
	expect := newTestRequest(ts)
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("expected %#v, got %#v", expect, r)
	}
}
