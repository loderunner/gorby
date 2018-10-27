package main

import (
	"fmt"
	"io"
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
		defer func() {
			src.Close()
			dst.Close()
			log.Debugf("closed tunnel %s->%s", src.RemoteAddr(), dst.RemoteAddr())
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
			for _, err := range []error{readErr, writeErr} {
				if err != nil {
					if netErr, ok := err.(net.Error); ok {
						log.Debugf("network error: %s", err)
						if !netErr.Temporary() {
							return
						}
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
	fmt.Fprintf(
		&logMessage, "[%s] %s %s %s",
		time.Now().Local().Format(timestampFormat),
		req.Proto,
		req.Method,
		req.Host,
	)

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
