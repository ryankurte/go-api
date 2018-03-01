package api

import (
	"errors"
	"net/http"
	"path"

	"github.com/gocraft/web"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"

	"github.com/ryankurte/go-api/lib/options"
	"github.com/ryankurte/go-api/lib/router"
	"github.com/ryankurte/go-api/lib/security"
	"github.com/ryankurte/go-api/lib/servers"
)

// API is a core API server instance
type API struct {
	router.Router
	options      *options.Base
	logger       log.FieldLogger
	server       servers.Handler
	sessionStore sessions.Store
}

// New creates a new API server
func New(ctx interface{}, o *options.Base) (*API, error) {
	var err error
	a := API{
		options: o,
		logger:  log.New().WithField("module", "core"),
	}

	// Create an API router
	base := web.New(ctx)
	a.Router = router.New(base, ctx, "")

	// Attach session storage
	if o.CookieSecret == "" {
		o.CookieSecret, err = options.GenerateSecret(256)
		if err != nil {
			return nil, err
		}
	}
	sessionStore := sessions.NewCookieStore([]byte(o.CookieSecret))
	sessionStore.Options.Secure = true
	sessionStore.Options.HttpOnly = true
	if o.ExternalAddress != "" {
		sessionStore.Options.Domain = o.ExternalAddress
	}
	a.sessionStore = sessionStore

	// Enable static file hosting if configured
	if o.StaticDir != "" {
		staticPath := path.Clean(o.StaticDir)
		base = base.Middleware(web.StaticMiddleware(staticPath))
		log.Printf("Serving static content from: %s\n", staticPath)
	}

	// Enable endpoint logging if specified
	if o.LogEndpoints {
		base = base.Middleware(web.LoggerMiddleware)
	}

	// Setup handlers
	var h http.Handler = base
	h = security.CORS(base, o)
	h = security.CSP(h, o)

	// Create server instance
	var server servers.Handler
	switch o.Mode {
	case options.ModeHTTP:
		server = servers.NewHTTP(o, h)
	case options.ModeLambda:
		server = servers.NewLambda(o, h)
	default:
		log.Errorf("Unhandled mode: '%s'", o.Mode)
		return nil, errors.New("unhandled server mode")
	}
	a.server = server

	return &a, nil
}

// SessionStore fetches the API service session store
func (api *API) SessionStore() sessions.Store {
	return api.sessionStore
}

// Run launches an API server
func (api *API) Run() {
	api.server.Run()
}

// Close closes an API server (if bound)
func (api *API) Close() {
	api.server.Close()
}
