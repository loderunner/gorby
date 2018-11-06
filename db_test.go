package main

import (
	"os"
	"strconv"
	"testing"
	"time"
)

var tsReq = time.Now()
var tsResp = tsReq.Add(20 * time.Millisecond)

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
	VALUES (1,` + strconv.FormatInt(tsResp.Unix(), 10) + `,"HTTP/1.1","200 OK",200,"[]",11,"Hello World",NULL,1);`
	_, err := db.Exec(fixturesSQL)
	if err != nil {
		panic(err.Error())
	}
}

func tearDown() {
}

func TestAddRequest(t *testing.T) {
	ts := time.Now()
	req := newTestRequest(ts)
	_, err := AddRequest(req)
	if err != nil {
		t.Fatalf("couldn't add request to DB: %s", err)
	}
}

func TestAddResponse(t *testing.T) {
	ts := time.Now()
	resp := newTestResponse(ts)
	_, err := AddResponse(resp, 2)
	if err != nil {
		t.Fatalf("couldn't add response to DB: %s", err)
	}
}
