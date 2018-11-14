package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const timestampFormat = "2006-01-02 15:04:05 -0700 MST"

func establishTunnel(req *http.Request, clientConn net.Conn) error {
	ctx := req.Context()
	srv, ok := ctx.Value(http.ServerContextKey).(*http.Server)
	if !ok {
		return fmt.Errorf("couldn't get server from request")
	}
	serverConn, err := net.Dial("tcp", req.Host)
	if err != nil {
		return err
	}

	srv.RegisterOnShutdown(func() {
		log.Debugf("shutting down hijacked connections")
		clientConn.Close()
		serverConn.Close()
	})

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

type Proxy struct {
	listeners []chan RequestResponse
	mu        sync.Mutex
}

func NewProxyHandler() *Proxy {
	return &Proxy{listeners: make([]chan RequestResponse, 0)}
}

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
	reqBody, addErr := ioutil.ReadAll(req.Body)
	if addErr == nil {
		r, addErr = NewRequest(reqTS, req, ioutil.NopCloser(bytes.NewBuffer(reqBody)))
	}
	if addErr == nil {
		r.ID, addErr = AddRequest(r)
	}
	if addErr != nil {
		log.Errorf("error adding request: %s", addErr)
	}
	go p.dispatch(r, nil)

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
		err = establishTunnel(req, conn)
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
		r2.ID, respErr = AddResponse(r2, r.ID)
	}
	if respErr != nil {
		log.Errorf("error adding response: %s", respErr)
	}
	go p.dispatch(r, r2)

	fmt.Fprintf(&logMessage, " - %s %d", resp.Status, resp.ContentLength)
	for k, h := range resp.Header {
		for _, v := range h {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, ioutil.NopCloser(bytes.NewBuffer(respBody)))
}

func (p *Proxy) Subscribe() <-chan RequestResponse {
	c := make(chan RequestResponse, 1)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.listeners = append(p.listeners, c)
	return c
}

func (p *Proxy) dispatch(req *Request, resp *Response) {
	rr := RequestResponse{req, resp}
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < len(p.listeners); i++ {
		c := p.listeners[i]
		select {
		case c <- rr:
			log.Debugf("request & response sent to listener %d", i)
		case <-time.After(time.Second):
			log.Warnf("timed out sending request & response to listener %d, closing and unsubscribing", i)
			p.listeners = append(p.listeners[:i], p.listeners[i+1:]...)
			i--
			close(c)
		}
	}
}
