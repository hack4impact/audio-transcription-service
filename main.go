// Package main initializes a web server.
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // import for side effects

		"github.com/hack4impact/transcribe4all/web"
)

func main() {
	router := web.NewRouter()
	middlewareRouter := web.ApplyMiddleware(router)
	config, configErr := parseConfigFile("config.toml")
	if configErr == nil {
		// replace this with your actual use of config
		fmt.Printf("%+v\n", *config)
	}
	uploadFileToBackblaze("testfile.wav")

	// serve http
	http.Handle("/", middlewareRouter)
	http.ListenAndServe(":8080", nil)
}
