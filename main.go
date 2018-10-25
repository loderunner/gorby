package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var logMessage strings.Builder
	defer func() {
		log.Println(logMessage.String())
	}()
	fmt.Fprintf(&logMessage, "[%s] %s %s %s", time.Now().Local(), req.Proto, req.Method, req.Host)

	if req.Method == http.MethodConnect {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Fprintf(&logMessage, " - error: %s", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer res.Body.Close()

	fmt.Fprintf(&logMessage, " - %d %s %d", res.StatusCode, res.Status, res.ContentLength)
	for k, h := range res.Header {
		for _, v := range h {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func main() {
	var h handler
	err := http.ListenAndServe(":8080", &h)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
