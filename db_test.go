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

	fixturesSQL := `INSERT INTO request (id,timestamp,proto,method,host,path,header,content_length,body,trailer,query,form)
	VALUES (1,` + strconv.FormatInt(tsReq.Unix(), 10) + `,"HTTP/1.1","POST","https://example.com","/","[]",11,"Hello World",NULL,NULL,NULL);
	INSERT INTO response (id,timestamp,proto,status,status_code,header,content_length,body,trailer,request) 
	VALUES (1,` + strconv.FormatInt(tsRes.Unix(), 10) + `,"HTTP/1.1","200 OK",200,"[]",11,"Hello World",NULL,1);`
	_, err := db.Exec(fixturesSQL)
	if err != nil {
		panic(err.Error())
	}
}

func tearDown() {
}

func TestAddRequest(t *testing.T) {
	req := newTestRequest(tsReq)
	_, err := AddRequest(req)
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
