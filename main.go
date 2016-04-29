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

		url, err := uploadFileToBackblaze("testfile4.wav", config.AccountID, config.ApplicationKey, config.BucketName)
		if err != nil {
			// replace this with your actual use of the url of the file
			fmt.Println(url)
		}
	}

	// serve http
	http.Handle("/", middlewareRouter)
	http.ListenAndServe(":8080", nil)
}
