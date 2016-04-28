// Package main initializes a web server.
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // import for side effects
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/hack4impact/transcribe4all/config"
	"github.com/hack4impact/transcribe4all/web"
)

// Config object
var Config AppConfig

func init() {
	log.SetOutput(os.Stderr)
	if config.Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	router := web.NewRouter()
	middlewareRouter := web.ApplyMiddleware(router)
	Config, err := parseConfigFile("config.toml")
	if err != nil {
		panic(fmt.Sprintf("%+v\n", *Config))
	}

	// serve http
	http.Handle("/", middlewareRouter)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}
