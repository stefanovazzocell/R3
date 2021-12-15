package main

import (
	"net/http"

	mux "github.com/julienschmidt/httprouter"
)

// NewRouter - Create a Router object
func newRouter() *mux.Router {
	router := mux.New()

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		header.Set("Access-Control-Allow-Origin", CORSorigin)

		w.WriteHeader(http.StatusNoContent)
	})

	for _, route := range routes {
		router.Handle(route.Method, route.Pattern, route.Handle)
	}

	return router
}
