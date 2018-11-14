package main

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

var tsReq = time.Now()
var tsResp = tsReq.Add(20 * time.Millisecond)

const testDBLocation = "/tmp/gorby_test.sqlite"

type testHook struct {
	t *testing.T
}

func (h testHook) Levels() []log.Level { return log.AllLevels }
func (h testHook) Fire(e *log.Entry) error {
	h.t.Logf(e.Message)
	return nil
}

func newTestHook(t *testing.T) log.Hook {
	return &testHook{t}
}

func TestMain(m *testing.M) {
	setUp()        // Setup for tests
	res := m.Run() // Run the actual tests
	tearDown()     // Teardown after running the tests
	os.Exit(res)
}

func setUp() {
	db.Close()
	initDB(testDBLocation)

	fixturesSQL := `INSERT INTO request (id,timestamp,proto,method,host,path,header,content_length,body,trailer,query,form)
	VALUES (1,?,"HTTP/1.1","POST","example.com","/","{}",11,"Hello World",NULL,NULL,NULL);
	INSERT INTO response (id,timestamp,proto,status,status_code,header,content_length,body,trailer,form,request) 
	VALUES (1,?,"HTTP/1.1","200 OK",200,"{}",11,"Hello World",NULL,NULL,1);`
	_, err := db.Exec(fixturesSQL, tsReq, tsResp)
	if err != nil {
		panic(err.Error())
	}
}

func tearDown() {
	os.Remove(testDBLocation)
}

func TestAddRequest(t *testing.T) {
	log.AddHook(newTestHook(t))
	defer log.StandardLogger().ReplaceHooks(log.LevelHooks{})

	ts := time.Now()
	req := newTestRequest(ts)
	_, err := AddRequest(req)
	if err != nil {
		t.Fatalf("couldn't add request to DB: %s", err)
	}
}

func TestAddResponse(t *testing.T) {
	log.AddHook(newTestHook(t))
	defer log.StandardLogger().ReplaceHooks(log.LevelHooks{})

	ts := time.Now()
	resp := newTestResponse(ts)
	_, err := AddResponse(resp, 2)
	if err != nil {
		t.Fatalf("couldn't add response to DB: %s", err)
	}
}

func TestListRequests(t *testing.T) {
	log.AddHook(newTestHook(t))
	defer log.StandardLogger().ReplaceHooks(log.LevelHooks{})

	reqs, resps, err := ListRequests(time.Time{}, time.Unix(9999999999, 0), -1)
	if err != nil {
		t.Fatalf("couldn't list requests: %s", err)
	}
	if len(reqs) == 0 {
		t.Fatalf("no requests")
	}
	if len(resps) == 0 {
		t.Fatalf("no responses")
	}

	expectedReq := Request{
		ID:            1,
		Proto:         "HTTP/1.1",
		Method:        http.MethodPost,
		Host:          "example.com",
		Path:          "/",
		ContentLength: 11,
		Header:        map[string][]string{},
		Body:          []byte("Hello World"),
		Trailer:       nil,
		Query:         nil,
		Form:          nil,
	}
	if !tsReq.Equal(reqs[0].Timestamp) {
		t.Errorf("expected request timestamp %s, got %s", tsReq, reqs[0].Timestamp)
	}
	reqs[0].Timestamp = time.Time{}
	if !reflect.DeepEqual(reqs[0], &expectedReq) {
		t.Errorf("expected request %#v, got %#v", expectedReq, *reqs[0])
	}

	expectedResp := Response{
		ID:            1,
		Proto:         "HTTP/1.1",
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		ContentLength: 11,
		Header:        map[string][]string{},
		Body:          []byte("Hello World"),
		Trailer:       nil,
		Form:          nil,
	}
	if !tsResp.Equal(resps[0].Timestamp) {
		t.Errorf("expected response timestamp %s, got %s", tsResp, resps[0].Timestamp)
	}
	resps[0].Timestamp = time.Time{}
	if !reflect.DeepEqual(resps[0], &expectedResp) {
		t.Errorf("expected response %#v, got %#v", expectedResp, *resps[0])
	}
}

func TestListRequestsWithoutResponse(t *testing.T) {
	log.AddHook(newTestHook(t))
	defer log.StandardLogger().ReplaceHooks(log.LevelHooks{})

	ts := time.Now()
	req := newTestRequest(ts)
	var err error
	req.ID, err = AddRequest(req)
	if err != nil {
		t.Fatalf("couldn't add request to DB: %s", err)
	}

	reqs, resps, err := ListRequests(time.Time{}, time.Unix(9999999999, 0), -1)
	if err != nil {
		t.Fatalf("couldn't list requests: %s", err)
	}
	if len(reqs) == 0 {
		t.Fatalf("no requests")
	}
	if len(resps) == 0 {
		t.Fatalf("no responses")
	}

	if !ts.Equal(reqs[len(reqs)-1].Timestamp) {
		t.Errorf("expected request timestamp %s, got %s", ts, reqs[len(reqs)-1].Timestamp)
	}
	req.Timestamp = time.Time{}
	reqs[len(reqs)-1].Timestamp = time.Time{}
	if !reflect.DeepEqual(reqs[len(reqs)-1], req) {
		t.Errorf("expected request %#v, got %#v", req, reqs[len(reqs)-1])
	}

	if resps[len(resps)-1] != nil {
		t.Errorf("expected response <nil>, got %#v", resps[len(resps)-1])
	}
}
