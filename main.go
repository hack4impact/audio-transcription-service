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

	// serve http
	http.Handle("/", middlewareRouter)
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	log.Infof("Server is running at http://localhost:%d", config.Config.Port)
	addr := fmt.Sprintf(":%d", config.Config.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error(err)
	}
}
