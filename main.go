package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
)

func initLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&easy.Formatter{
		LogFormat: "%msg%\n",
	})
}

func main() {

	initLogger()

	var p Proxy
	err := http.ListenAndServe(":8080", &p)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
