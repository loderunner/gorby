package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func newRequest() *http.Request {
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

func newResponse() *http.Response {
	const responseString = "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 606\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		"<html><body>Hello World!</body><html>"
	reader := bufio.NewReader(strings.NewReader(responseString))
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		panic(err.Error())
	}
	resp.Body = ioutil.NopCloser(reader)
	return resp
}

var tsReq = time.Now()
var tsRes = tsReq.Add(20 * time.Millisecond)

func TestMain(m *testing.M) {
	setUp()        // Setup for tests
	res := m.Run() // Run the actual tests
	tearDown()     // Teardown after running the tests
	os.Exit(res)
}

func setUp() {
	db.Close()
	initDB()

	fixturesSQL := `INSERT INTO request (timestamp,proto,method,host,path,header,content_length,body,trailer)
	VALUES (` + strconv.FormatInt(tsReq.Unix(), 10) + `,"HTTP/1.1","POST","https://example.com","/","[]",11,"Hello World",NULL);
	INSERT INTO response (timestamp,proto,status,status_code,header,content_length,body,trailer,request) 
	VALUES (` + strconv.FormatInt(tsRes.Unix(), 10) + `,"HTTP/1.1","200 OK",200,"[]",11,"Hello World",NULL,1);`
	_, err := db.Exec(fixturesSQL)
	if err != nil {
		panic(err.Error())
	}
}

func tearDown() {
}

func TestAddRequest(t *testing.T) {
	ts := time.Date(2018, time.November, 2, 23, 38, 0, 0, time.UTC)
	req := newRequest()
	_, err := AddRequest(ts, req, req.Body)
	if err != nil {
		t.Fatalf("couldn't add request to DB: %s", err)
	}
}

func TestAddResponse(t *testing.T) {
	ts := time.Date(2018, time.November, 2, 23, 38, 0, int(20*time.Millisecond), time.UTC)
	resp := newResponse()
	_, err := AddResponse(ts, resp, resp.Body, 2)
	if err != nil {
		t.Fatalf("couldn't add response to DB: %s", err)
	}
}
