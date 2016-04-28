// Package main initializes a web server.
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // import for side effects

	"github.com/hack4impact/transcribe4all/web"
)

// Config object
var Config AppConfig

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
