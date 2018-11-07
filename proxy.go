package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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

type Proxy struct{}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var logMessage strings.Builder
	defer func() {
		log.Info(logMessage.String())
	}()
	reqTS := time.Now().Local()
	fmt.Fprintf(
		&logMessage, "[%s] %s %s %s",
		reqTS.Format(timestampFormat),
		req.Proto,
		req.Method,
		req.Host,
	)

	var r *Request
	var reqID int64
	reqBody, addErr := ioutil.ReadAll(req.Body)
	if addErr == nil {
		r, addErr = NewRequest(reqTS, req, ioutil.NopCloser(bytes.NewBuffer(reqBody)))
	}
	if addErr == nil {
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
	respTS := time.Now()

	var r2 *Response
	respBody, respErr := ioutil.ReadAll(resp.Body)
	if respErr == nil {
		r2, respErr = NewResponse(respTS, resp, ioutil.NopCloser(bytes.NewBuffer(respBody)))
	}
	if respErr == nil {
		_, respErr = AddResponse(r2, reqID)
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
