package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func newRequest() *http.Request {
	req, err := http.NewRequest(
		http.MethodPost,
		"https://github.com/loderunner/gorby",
		strings.NewReader(`{"message":"Hello World"}`),
	)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err.Error())
	}
	return req
}

func newResponse() *http.Response {
	resp := &http.Response{}
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
	_, err := AddRequest(ts, req)
	if err != nil {
		t.Fatalf("couldn't add request to DB: %s", err)
	}
}

func TestAddResponse(t *testing.T) {
	ts := time.Date(2018, time.November, 2, 23, 38, 0, int(20*time.Millisecond), time.UTC)
	req := newRequest()
	_, err := AddRequest(ts, req)
	if err != nil {
		t.Fatalf("couldn't add request to DB: %s", err)
	}
}
