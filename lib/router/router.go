package router

import (
	"github.com/gorilla/mux"
)

// Router instance
type Router struct {
	// Underlying gocraft/web router instance
	router *mux.Router
	// Router context object
	ctx interface{}
	// Path of the current router
	path string
	// Endpoints attached to the router
	endpoints []endpoint
	// Error handling function attached to the router
	//errorHandler errors.ErrorHandler
}

// endpoint for internal use
type endpoint struct {
	path string
	// Input object
	i interface{}
	// Output object
	o interface{}
	// Base function
	f interface{}
	// Wrapped function
	w interface{}
}

// Create an API router instance (internal use only)
func newRouter(router *mux.Router, ctx interface{}, path string) Router {
	return Router{
		router:    router,
		ctx:       ctx,
		path:      path,
		endpoints: make([]endpoint, 0),
		//errorHandler: errors.DefaultErrorHandler,
	}
}
