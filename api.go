package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	var flusher http.Flusher
	accept := req.Header.Get("Accept")
	if accept == "text/event-stream" {
		var ok bool
		flusher, ok = w.(http.Flusher)
		if !ok {
			log.Error("cannot stream response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	startStr := req.FormValue("start")
	var startTime time.Time
	if len(startStr) == 0 {
		startTime = time.Time{}
	} else {
		startTime, err = time.Parse(time.RFC3339Nano, startStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid start time")
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid end time")
			return
		}
	}

	limitStr := req.FormValue("limit")
	var limit int64 = -1
	if len(limitStr) > 0 {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid limit")
			return
		}
	}

	requests, responses, err := ListRequests(startTime, endTime, limit)
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

	if accept == "text/event-stream" {
		flusher.Flush()
		HandleStreamRequests(w, req, flusher, requestResponses)
	} else {
		b, err := json.Marshal(requestResponses)
		if err != nil {
			log.Errorf("couldn't marshal requests and responses: %s", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

func HandleStreamRequests(w http.ResponseWriter, req *http.Request, f http.Flusher, requestResponses []requestResponse) {
	log.Debugf("handling server-sent event stream")

	closer, ok := w.(http.CloseNotifier)
	if !ok {
		log.Warning("couldn't notify connection closed by client")
		return
	}

	for _, rr := range requestResponses {
		b, err := json.Marshal(rr)
		if err != nil {
			log.Errorf("couldn't marshal request and response: %s", err)
			continue
		}
		w.Write([]byte("data:"))
		w.Write(b)
		w.Write([]byte("\n\n"))
		f.Flush()
	}

	for {
		select {
		case t := <-time.Tick(3 * time.Second):
			fmt.Fprintf(w, "data: %s\n\n", t.Format(time.RFC3339))
			log.Debugf("sending timestamp event %s", t)
			f.Flush()
		case <-closer.CloseNotify():
			log.Debugf("closing")
			return
		}
	}
}
