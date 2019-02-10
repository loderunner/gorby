package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/loderunner/popt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/t-tomalak/logrus-easy-formatter"
	"golang.org/x/sync/errgroup"
)

var proxyServer, apiServer http.Server

func initConfiguration() error {
	if err := popt.AddAndBindOptions(options, pflag.CommandLine); err != nil {
		panic(err.Error())
	}
	pflag.CommandLine.SortFlags = false

	pflag.Parse()

	viper.SetConfigName("gorby")
	confPathFlag := pflag.Lookup("conf")
	if confPathFlag != nil && confPathFlag.Changed {
		viper.SetConfigFile(confPathFlag.Value.String())
	}
	err := viper.ReadInConfig()
	if err != nil && confPathFlag != nil && confPathFlag.Changed {
		return err
	}

	return nil
}

func initLogger() {
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	} else if viper.GetBool("quiet") {
		log.SetLevel(log.PanicLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
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

	err := initConfiguration()
	if err != nil {
		log.Fatal(err.Error())
	}
	initLogger()
	initSignals()

	var g errgroup.Group

	initDB(viper.GetString("db"))

	proxy := NewProxyHandler()
	proxyServer.Addr = viper.GetString("address")
	proxyServer.Handler = proxy

	apiServer.Addr = viper.GetString("api-address")
	apiServer.Handler = NewAPIHandler(proxy)

	g.Go(func() error {
		log.Infof("starting HTTP proxy on %s", proxyServer.Addr)
		err := proxyServer.ListenAndServe()
		return err
	})

	g.Go(func() error {
		log.Infof("starting API on %s", apiServer.Addr)
		err := apiServer.ListenAndServe()
		return err
	})

	err = g.Wait()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("fatal error: %s", err)
	}
}
