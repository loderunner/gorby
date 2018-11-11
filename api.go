package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type APIServer struct{}

func NewAPIHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/requests", HandleListRequests)

	return mux
}

func HandleListRequests(w http.ResponseWriter, req *http.Request) {
	var err error

	startStr := req.FormValue("start")
	var startTime time.Time
	if len(startStr) == 0 {
		startTime = time.Time{}
	} else {
		startTime, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			errMsg := fmt.Sprintf("invalid start time: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
	}

	endStr := req.FormValue("end")
	var endTime time.Time
	if len(endStr) == 0 {
		endTime = MaxUnixTime
	} else {
		endTime, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			errMsg := fmt.Sprintf("invalid end time: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
	}

	log.Infof("ListRequests: start=%s,end=%s", startTime, endTime)
}
