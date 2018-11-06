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

func newTestHTTPRequest() *http.Request {
	requestString := "POST /?msg=hello+world HTTP/1.1\r\n" +
		"Accept: application/x-www-form-urlencoded, */*\r\n" +
		"Accept-Encoding: gzip, deflate\r\n" +
		"Connection: keep-alive\r\n" +
		"Content-Length: 20\r\n" +
		"Content-Type: application/x-www-form-urlencoded; charset=utf-8\r\n" +
		"Host: example.com\r\n" +
		"\r\n" +
		"message=Hello+World%21"
	reader := bufio.NewReader(strings.NewReader(requestString))
	req, err := http.ReadRequest(reader)
	if err != nil {
		panic(err.Error())
	}
	req.Body = ioutil.NopCloser(reader)
	return req
}

func newTestRequest(ts time.Time) *Request {
	r := &Request{
		Timestamp:     ts,
		Proto:         "HTTP/1.1",
		Method:        http.MethodPost,
		Host:          "example.com",
		Path:          "/",
		ContentLength: 20,
		Header: map[string][]string{
			"Accept":          {"application/x-www-form-urlencoded, */*"},
			"Accept-Encoding": {"gzip, deflate"},
			"Connection":      {"keep-alive"},
			"Content-Length":  {"20"},
			"Content-Type":    {"application/x-www-form-urlencoded; charset=utf-8"},
		},
		Body:    []byte(`message=Hello+World%21`),
		Trailer: nil,
		Query:   map[string][]string{"msg": {"hello world"}},
		Form:    map[string][]string{"message": {"Hello World!"}},
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
