package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func newTestHTTPResponse() *http.Response {
	const responseString = "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 27\r\n" +
		"Content-Type: application/json\r\n" +
		"\r\n" +
		"{\"message\": \"Hello World!\"}"
	reader := bufio.NewReader(strings.NewReader(responseString))
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		panic(err.Error())
	}
	resp.Body = ioutil.NopCloser(reader)
	return resp
}

func newTestResponse(ts time.Time) *Response {
	r := &Response{
		Timestamp:     ts,
		Proto:         "HTTP/1.1",
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		ContentLength: 27,
		Header: map[string][]string{
			"Content-Length": {"27"},
			"Content-Type":   {"application/json"},
		},
		Body:    []byte(`{"message": "Hello World!"}`),
		Trailer: nil,
		Form:    nil,
	}
	return r
}

func TestNewResponse(t *testing.T) {
	ts := time.Now()
	resp := newTestHTTPResponse()
	r, err := NewResponse(ts, resp, resp.Body)
	if err != nil {
		t.Fatalf("error creating respuest: %s", err)
	}
	expect := newTestResponse(ts)
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("expected %#v, got %#v", expect, r)
	}
}
