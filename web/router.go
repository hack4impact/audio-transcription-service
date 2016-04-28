// Package web implements http routing logic for the application.
package web

import (
	"net/http"

	logMiddleware "github.com/bakins/logrus-middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// NewRouter creates and returns a mux.Router with default routes.
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// ApplyMiddleware wraps the router in some middleware. This middleware includes
// logging and gzip compression.
func ApplyMiddleware(router http.Handler) http.Handler {
	loggingHandler := func(h http.Handler) http.Handler {
		m := new(logMiddleware.Middleware)
		return m.Handler(h, "")
	}
	middlewareRouter := alice.New(handlers.CompressHandler, loggingHandler).Then(router)
	return middlewareRouter
}
