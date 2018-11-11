package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"golang.org/x/sync/errgroup"
)

var proxyServer, apiServer http.Server

func initLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&easy.Formatter{
		LogFormat: "%msg%\n",
	})
}

func initSignals() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Infof("shutting down...")
		err := proxyServer.Shutdown(context.Background())
		if err != nil {
			log.Errorf("error shutting down proxy server: %s", err)
		}
		err = apiServer.Shutdown(context.Background())
		if err != nil {
			log.Errorf("error shutting down API server: %s", err)
		}
	}()
}

func main() {

	initLogger()
	initSignals()

	var g errgroup.Group

	g.Go(func() error {
		proxyServer.Addr = ":8080"
		proxyServer.Handler = &Proxy{}
		err := proxyServer.ListenAndServe()
		return err
	})

	g.Go(func() error {
		apiServer.Addr = ":8081"
		apiServer.Handler = NewAPIHandler()
		err := apiServer.ListenAndServe()
		return err
	})

	err := g.Wait()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("fatal error: %s", err)
	}
}
