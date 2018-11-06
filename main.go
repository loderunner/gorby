package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
)

const timestampFormat = "2006-01-02 15:04:05 -0700 MST"

func establishTunnel(host string, clientConn net.Conn) error {
	serverConn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	forward := func(src, dst net.Conn) {
		var err error
		defer func() {
			if err != net.ErrWriteToConnected {
				log.Debugf("closing tunnel %s->%s", src.RemoteAddr(), dst.RemoteAddr())
				src.Close()
				dst.Close()
			}
		}()

		var buf [1024]byte
		for {
			var writeErr error
			bytesRead, readErr := src.Read(buf[:])
			if bytesRead > 0 {
				log.Debugf("read %d bytes from %s", bytesRead, src.RemoteAddr())
				var bytesWritten int
				bytesWritten, writeErr = dst.Write(buf[:bytesRead])
				log.Debugf("wrote %d bytes to %s", bytesWritten, dst.RemoteAddr())
			}
			for _, err = range []error{readErr, writeErr} {
				if err != nil {
					if netErr, ok := err.(net.Error); ok {
						if !netErr.Temporary() {
							log.Errorf("network error: %s", err)
							return
						}
						log.Debugf("network error: %s", err)
					} else {
						if err != io.EOF {
							log.Errorf("error: %s", err)
						}
						return
					}
				}
			}
		}
	}

	go forward(clientConn, serverConn)
	go forward(serverConn, clientConn)

	return nil
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var logMessage strings.Builder
	defer func() {
		log.Info(logMessage.String())
	}()
	ts := time.Now().Local()
	fmt.Fprintf(
		&logMessage, "[%s] %s %s %s",
		ts.Format(timestampFormat),
		req.Proto,
		req.Method,
		req.Host,
	)

	var r *Request
	var reqID int64
	reqBody, addErr := ioutil.ReadAll(req.Body)
	if addErr == nil {
		r, addErr = NewRequest(ts, req, ioutil.NopCloser(bytes.NewBuffer(reqBody)))
	}
	if addErr != nil {
		reqID, addErr = AddRequest(r)
	}
	if addErr != nil {
		log.Errorf("error adding request: %s", addErr)
	}

	if req.Method == http.MethodConnect {
		h, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "couldn't open TCP connection", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		conn, _, err := h.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conn.SetDeadline(time.Time{})
		err = establishTunnel(req.Host, conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Fprintf(&logMessage, " - error: %s", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	respBody, respErr := ioutil.ReadAll(resp.Body)
	if respErr == nil {
		_, respErr = AddResponse(ts, resp, ioutil.NopCloser(bytes.NewBuffer(respBody)), reqID)
	}
	if respErr != nil {
		log.Errorf("error adding response: %s", respErr)
	}

	fmt.Fprintf(&logMessage, " - %s %d", resp.Status, resp.ContentLength)
	for k, h := range resp.Header {
		for _, v := range h {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, ioutil.NopCloser(bytes.NewBuffer(respBody)))
}

func initLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&easy.Formatter{
		LogFormat: "%msg%\n",
	})
}

func main() {

	initLogger()

	var h handler
	err := http.ListenAndServe(":8080", &h)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
