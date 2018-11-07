package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"golang.org/x/sync/errgroup"
)

var proxyServer http.Server

func initLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&easy.Formatter{
		LogFormat: "%msg%\n",
	})
}

func main() {

	initLogger()

	var g errgroup.Group

	g.Go(func() error {
		proxyServer.Addr = ":8080"
		proxyServer.Handler = &Proxy{}
		err := proxyServer.ListenAndServe()
		return err
	})

	err := g.Wait()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("fatal error: %s", err)
	}
}
