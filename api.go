package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type requestResponse struct {
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

type APIServer struct{}

func NewAPIHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/requests", HandleListRequests)

	return mux
}

func HandleListRequests(w http.ResponseWriter, req *http.Request) {
	var err error

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	startStr := req.FormValue("start")
	var startTime time.Time
	if len(startStr) == 0 {
		startTime = time.Time{}
	} else {
		startTime, err = time.Parse(time.RFC3339Nano, startStr)
		if err != nil {
			fmt.Fprintf(w, "invalid start time")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	endStr := req.FormValue("end")
	var endTime time.Time
	if len(endStr) == 0 {
		endTime = MaxUnixTime
	} else {
		endTime, err = time.Parse(time.RFC3339Nano, endStr)
		if err != nil {
			fmt.Fprintf(w, "invalid end time")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	requests, responses, err := ListRequests(startTime, endTime)
	if err != nil {
		log.Errorf("couldn't get requests and responses: %s", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if len(requests) != len(responses) {
		log.Errorf("length mismatch: %d requests - %d responses", len(requests), len(responses))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	requestResponses := make([]requestResponse, len(requests))
	for i := 0; i < len(requests); i++ {
		requestResponses[i].Request = requests[i]
		requestResponses[i].Response = responses[i]
	}

	b, err := json.Marshal(requestResponses)
	if err != nil {
		log.Errorf("couldn't marshal requests and responses: %s", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
