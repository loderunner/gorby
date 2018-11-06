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
		"Content-Length: 20\r\n" +
		"Content-Type: application/x-www-form-urlencoded; charset=utf-8\r\n" +
		"\r\n" +
		"message=Hello+World%21"
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
		ContentLength: 20,
		Header: map[string][]string{
			"Content-Length": {"20"},
			"Content-Type":   {"application/x-www-form-urlencoded; charset=utf-8"},
		},
		Body:    []byte(`message=Hello+World%21`),
		Trailer: nil,
		Form:    map[string][]string{"message": {"Hello World!"}},
	}
	return r
}

func TestNewResponse(t *testing.T) {
	ts := time.Now()
	resp := newTestHTTPResponse()
	r, err := NewResponse(ts, resp, resp.Body)
	if err != nil {
		t.Fatalf("error creating response: %s", err)
	}
	expect := newTestResponse(ts)
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("expected %#v, got %#v", expect, r)
	}
}
