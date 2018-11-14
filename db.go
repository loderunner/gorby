package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func init() {
	initDB("/tmp/gorby.sqlite")
}

func initDB(path string) {
	// Create an in-memory SQLite store
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		panic(err.Error())
	}

	// Create schema
	_, err = db.Exec(schemaSQL)
	if err != nil {
		panic(err.Error())
	}
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
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
	query, err := json.Marshal(req.Query)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal request query variables: %s", err)
	}
	form, err := json.Marshal(req.Form)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal request form variables: %s", err)
	}
	res, err := db.Exec(
		`INSERT INTO request (timestamp,proto,method,host,path,header,content_length,body,trailer,query,form) 
        VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		req.Timestamp,
		req.Proto,
		req.Method,
		req.Host,
		req.Path,
		header,
		req.ContentLength,
		req.Body,
		trailer,
		query,
		form,
	)
	if err != nil {
		return 0, fmt.Errorf("SQL error: %s", err)
	}
	reqID, _ := res.LastInsertId()
	return reqID, nil
}

func AddResponse(resp *Response, reqID int64) (int64, error) {
	header, err := json.Marshal(resp.Header)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal response header: %s", err)
	}
	trailer, err := json.Marshal(resp.Trailer)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal response trailer: %s", err)
	}
	form, err := json.Marshal(resp.Form)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal response form variables: %s", err)
	}
	res, err := db.Exec(
		`INSERT INTO response (timestamp,proto,status,status_code,header,content_length,body,trailer,form,request) 
        VALUES (?,?,?,?,?,?,?,?,?,?)`,
		resp.Timestamp,
		resp.Proto,
		resp.Status,
		resp.StatusCode,
		header,
		resp.ContentLength,
		resp.Body,
		trailer,
		form,
		reqID,
	)
	if err != nil {
		return 0, fmt.Errorf("SQL error: %s", err)
	}
	respID, _ := res.LastInsertId()
	return respID, nil
}

func ListRequests(start, end time.Time, limit int64) ([]*Request, []*Response, error) {
	res, err := db.Query(
		`SELECT * FROM request_response 
		WHERE [req.timestamp] > ? AND [req.timestamp] <= ?
		LIMIT ?`,
		start, end, limit,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("SQL error: %s", err)
	}
	defer res.Close()

	reqs := make([]*Request, 0)
	resps := make([]*Response, 0)
	for res.Next() {
		var req Request
		var resp Response

		var reqID, respID sql.NullInt64
		var reqHeader, reqTrailer, reqQuery, reqForm, respHeader, respTrailer, respForm []byte
		var respTS NullTime
		var respProto, respStatus sql.NullString
		var respStatusCode, respContentLength sql.NullInt64

		err = res.Scan(
			&reqID,
			&req.Timestamp,
			&req.Proto,
			&req.Method,
			&req.Host,
			&req.Path,
			&reqHeader,
			&req.ContentLength,
			&req.Body,
			&reqTrailer,
			&reqQuery,
			&reqForm,
			&respID,
			&respTS,
			&respProto,
			&respStatus,
			&respStatusCode,
			&respHeader,
			&respContentLength,
			&resp.Body,
			&respTrailer,
			&respForm,
		)

		if err != nil {
			log.Warningf("couldn't scan row: %s", err)
			continue
		}

		if reqID.Valid {
			req.ID = reqID.Int64
		}
		if reqHeader != nil {
			err = json.Unmarshal(reqHeader, &req.Header)
			if err != nil {
				log.Warningf("couldn't unmarshal request header: %s", err)
			}
		}
		if reqTrailer != nil {
			err = json.Unmarshal(reqTrailer, &req.Trailer)
			if err != nil {
				log.Warningf("couldn't unmarshal request trailer: %s", err)
			}
		}
		if reqQuery != nil {
			err = json.Unmarshal(reqQuery, &req.Query)
			if err != nil {
				log.Warningf("couldn't unmarshal request query variables: %s", err)
			}
		}
		if reqForm != nil {
			err = json.Unmarshal(reqForm, &req.Form)
			if err != nil {
				log.Warningf("couldn't unmarshal request form variables: %s", err)
			}
		}
		reqs = append(reqs, &req)

		if respID.Valid {
			resp.ID = respID.Int64
			if respTS.Valid {
				resp.Timestamp = respTS.Time
			} else {
				log.Warningf("invalid response timestamp")
				continue
			}
			if respProto.Valid {
				resp.Proto = respProto.String
			} else {
				log.Warningf("invalid response proto")
				continue
			}
			if respStatus.Valid {
				resp.Status = respStatus.String
			} else {
				log.Warningf("invalid response status")
				continue
			}
			if respStatusCode.Valid {
				resp.StatusCode = int(respStatusCode.Int64)
			} else {
				log.Warningf("invalid response status code")
				continue
			}
			if respContentLength.Valid {
				resp.ContentLength = respContentLength.Int64
			} else {
				log.Warningf("invalid response content length")
				continue
			}
			if respHeader != nil {
				err = json.Unmarshal(respHeader, &resp.Header)
				if err != nil {
					log.Warningf("couldn't read response header: %s", err)
				}
			}
			if respTrailer != nil {
				err = json.Unmarshal(respTrailer, &resp.Trailer)
				if err != nil {
					log.Warningf("couldn't unmarshal response trailer: %s", err)
				}
			}
			if respForm != nil {
				err = json.Unmarshal(respForm, &resp.Form)
				if err != nil {
					log.Warningf("couldn't unmarshal response query variables: %s", err)
				}
			}
			resps = append(resps, &resp)
		} else {
			resps = append(resps, nil)
		}

	}

	if err := res.Err(); err != nil {
		return nil, nil, fmt.Errorf("row error: %s", err)
	}

	return reqs, resps, nil
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
	trailer        BLOB,
	query          BLOB,
	form           BLOB
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
	form           BLOB,
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
	[req.query],
	[req.form],
    [res.id],
    [res.timestamp],
    [res.proto],
    [res.status],
    [res.status_code],
	[res.header],
	[res.content_length],
    [res.body],
	[res.trailer],
	[res.form]
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
		   req.query,
		   req.form,
           res.id,
           res.timestamp,
           res.proto,
           res.status,
           res.status_code,
		   res.header,
		   res.content_length,
           res.body,
		   res.trailer,
		   res.form
      FROM request AS req
           LEFT OUTER JOIN
           response AS res ON req.id = res.request
     ORDER BY req.timestamp ASC;`
