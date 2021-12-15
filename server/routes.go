package main

import (
	mux "github.com/julienschmidt/httprouter"
)

// Route represents a URL that serves a specific resource.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handle  mux.Handle
}

// Routes are a list of Routes for this application.
type Routes []Route

var routes = Routes{
	Route{
		"Create/modify/delete a Share",
		"POST",
		"/v2/edit",
		handleEdit,
	},
	Route{
		"Read a Share",
		"POST",
		"/v2/get",
		handleGet,
	},
}
