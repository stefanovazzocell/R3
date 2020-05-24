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
		"Index",
		"GET",
		"/",
		handleIndex,
	},
	Route{
		"Server Status",
		"GET",
		"/status",
		handleStatus,
	},
	Route{
		"LinkShow",
		"POST",
		"/v1/link",
		handleGet,
	},
	Route{
		"LinkCreate",
		"POST",
		"/v1/links",
		handleCreate,
	},
	Route{
		"LinkEdit",
		"PUT",
		"/v1/link",
		handleEdit,
	},
	Route{
		"LinkDelete",
		"DELETE",
		"/v1/link",
		handleDelete,
	},
}
