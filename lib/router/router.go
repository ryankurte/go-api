package router

import (
	"fmt"
	"net/http"

	"github.com/gocraft/web"

	log "github.com/sirupsen/logrus"

	"github.com/ryankurte/go-api/lib/wrappers"
)

// Router instance
type Router struct {
	// Underlying gocraft/web router instance
	router *web.Router
	// Router context object
	ctx interface{}
	// Path of the current router
	path string
	// Endpoints attached to the router
	endpoints []endpoint
	// Error handling function attached to the router
	errorHandler wrappers.ErrorHandler
}

// New Creates an API router instance (internal use only)
func New(router *web.Router, ctx interface{}, path string) Router {
	return Router{
		router:       router,
		ctx:          ctx,
		path:         path,
		endpoints:    make([]endpoint, 0),
		errorHandler: wrappers.DefaultErrorHandler,
	}
}

func wrapGocraft(h wrappers.HTTPHandler) func(ctx interface{}, rw web.ResponseWriter, req *web.Request) {
	return func(ctx interface{}, rw web.ResponseWriter, req *web.Request) {
		h(ctx, rw, req.Request)
	}
}

// RegisterEndpoint Register a route to the API router.
// This takes a typed endpoint and generates a wrapper to handle
// translation and validation of input and output structures,
// as well as error handling for the endpoint.
func (r *Router) RegisterEndpoint(route string, method string, f interface{}) error {

	log.Infof("Router '%s' attaching route %s with method %s (f: %+V)", r.path, route, method, f)

	var w interface{}

	// Build endpoint wrapper
	h, err := wrappers.BuildEndpoint(method, f)
	if err != nil {
		return err
	}
	w = wrapGocraft(h)

	// Fetch endpoint input/output instances
	inType, outType := wrappers.GetTypes(f)

	// Save endpoint object for later traversal
	path := fmt.Sprintf("%s/%s:%s", r.path, route, method)
	r.endpoints = append(r.endpoints, endpoint{
		path: path,
		i:    inType,
		o:    outType,
		f:    f,
		w:    w,
	})

	// Bind to router
	return r.Register(route, method, w)
}

// Register registers a basic http or gocraft/web route handler without any modification
// This does not currently call any meta plugins.
func (r *Router) Register(route string, method string, f interface{}) error {
	log.Infof("Router '%s' attaching route '%s' with method '%s' (f: %+V)", r.path, route, method, f)

	switch method {
	case http.MethodGet:
		r.router = r.router.Get(route, f)
	case http.MethodPost:
		r.router = r.router.Post(route, f)
	case http.MethodPut:
		r.router = r.router.Put(route, f)
	case http.MethodDelete:
		r.router = r.router.Delete(route, f)
	case http.MethodPatch:
		r.router = r.router.Patch(route, f)
	case http.MethodHead:
		r.router = r.router.Head(route, f)
	case http.MethodOptions:
		r.router = r.router.Options(route, f)
	default:
		log.Errorf("Invalid HTTP method: %s", method)
		return fmt.Errorf("Invalid HTTP method: %s", method)
	}

	return nil
}

// Subrouter Creates a subrouter with a given context and path.
// As with gocraft, this context must have a pointer to the parent context as it's first field
func (r *Router) Subrouter(ctx interface{}, path string) *Router {
	// Create child from base router
	b := r.router.Subrouter(ctx, path)

	// Create API Router instance
	sr := New(b, ctx, path)

	return &sr
}

// RegisterMiddleware Attach dependency injected middleware to API router.
// This is not yet supported
func (r *Router) RegisterMiddleware() error {
	return fmt.Errorf("Dependency injected middleware not yet supported")
}

// Middleware Attach standard middleware to an API router
func (r *Router) Middleware(fn interface{}) {
	r.router = r.router.Middleware(fn)
}

// GetBaseRouter Fetch the underlying router.
// Note that operations on this will bypass any GoAPI magic
func (r *Router) GetBaseRouter() *web.Router {
	return r.router
}
