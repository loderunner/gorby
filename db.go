package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	initDB()
}

func initDB() {
	// Create an in-memory SQLite store
	var err error
	db, err = sql.Open("sqlite3", "/tmp/gorby.sqlite")
	if err != nil {
		panic(err.Error())
	}

	// Create schema
	_, err = db.Exec(schemaSQL)
	if err != nil {
		panic(err.Error())
	}
}

func AddRequest(req *Request) (int64, error) {
	header, err := json.Marshal(req.Header)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal request header: %s", err)
	}
	trailer, err := json.Marshal(req.Trailer)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal request trailer: %s", err)
	}
	res, err := db.Exec(
		`INSERT INTO request (timestamp,proto,method,host,path,header,content_length,body,trailer) 
        VALUES (?,?,?,?,?,?,?,?,?)`,
		req.Timestamp,
		req.Proto,
		req.Method,
		req.Host,
		req.Path,
		header,
		req.ContentLength,
		req.Body,
		trailer,
	)
	if err != nil {
		return 0, fmt.Errorf("SQL error: %s", err)
	}
	reqID, _ := res.LastInsertId()
	return reqID, nil
}

func AddResponse(ts time.Time, resp *http.Response, body io.ReadCloser, reqID int64) (int64, error) {
	header, err := json.Marshal(resp.Header)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal response header: %s", err)
	}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return 0, fmt.Errorf("couldn't read response body: %s", err)
	}
	trailer, err := json.Marshal(resp.Trailer)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal response trailer: %s", err)
	}
	res, err := db.Exec(
		`INSERT INTO response (timestamp,proto,status,status_code,header,content_length,body,trailer,request) 
        VALUES (?,?,?,?,?,?,?,?,?)`,
		ts,
		resp.Proto,
		resp.Status,
		resp.StatusCode,
		header,
		resp.ContentLength,
		b,
		trailer,
		reqID,
	)
	if err != nil {
		return 0, fmt.Errorf("SQL error: %s", err)
	}
	respID, _ := res.LastInsertId()
	return respID, nil
}

const schemaSQL = `-- Table: request
DROP TABLE IF EXISTS request;
CREATE TABLE request (
    id             INTEGER  PRIMARY KEY
                            NOT NULL
                            UNIQUE,
    timestamp      DATETIME NOT NULL,
    proto          STRING   NOT NULL,
    method         STRING   NOT NULL,
    host           STRING   NOT NULL,
    path           STRING   NOT NULL,
    header         BLOB,
    content_length INTEGER  NOT NULL,
    body           BLOB,
    trailer        BLOB
);


-- Table: response
DROP TABLE IF EXISTS response;
CREATE TABLE response (
    id             INTEGER  PRIMARY KEY
                            UNIQUE
                            NOT NULL,
    timestamp      DATETIME NOT NULL,
    proto          STRING   NOT NULL,
    status         STRING   NOT NULL,
    status_code    INTEGER  NOT NULL,
    header         BLOB,
    content_length INTEGER  NOT NULL,
    body           BLOB,
    trailer        BLOB,
    request        INTEGER  REFERENCES request (id) 
                            NOT NULL
                            UNIQUE
);


-- Index: idx_request_host
DROP INDEX IF EXISTS idx_request_host;
CREATE INDEX idx_request_host ON request (
    host
);


-- Index: idx_request_timestamp
DROP INDEX IF EXISTS idx_request_timestamp;
CREATE INDEX idx_request_timestamp ON request (
    timestamp
);


-- Index: idx_response_timestamp
DROP INDEX IF EXISTS idx_response_timestamp;
CREATE INDEX idx_response_timestamp ON response (
    timestamp
);


-- View: request_response
DROP VIEW IF EXISTS request_response;
CREATE VIEW request_response (
    [req.id],
    [req.timestamp],
    [req.proto],
    [req.method],
    [req.host],
    [req.path],
    [req.header],
    [req.content_length],
    [req.body],
    [req.trailer],
    [res.id],
    [res.timestamp],
    [res.proto],
    [res.status],
    [res.status_code],
    [res.header],
    [res.body],
    [res.trailer]
)
AS
    SELECT req.id,
           req.timestamp,
           req.proto,
           req.method,
           req.host,
           req.path,
           req.header,
           req.content_length,
           req.body,
           req.trailer,
           res.id,
           res.timestamp,
           res.proto,
           res.status,
           res.status_code,
           res.header,
           res.body,
           res.trailer
      FROM request AS req
           LEFT OUTER JOIN
           response AS res ON req.id = res.request
     ORDER BY req.timestamp ASC;`
