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
	if o.Session.Secret == "" {
		a.options.Session.Secret, err = options.GenerateSecret(256)
		if err != nil {
			return nil, err
		}
	}

	sessionStore := sessions.NewCookieStore([]byte(o.Session.Secret))
	if a.options.Session.DisableSecure {
		log.Warn("SECURE COOKIE FLAG IS DISABLED. DEVELOPMENT USE ONLY.")
		sessionStore.Options.Secure = false
	} else {
		sessionStore.Options.Secure = true
	}
	sessionStore.Options.HttpOnly = true
	if a.options.ExternalAddress != "" {
		sessionStore.Options.Domain = o.ExternalAddress
	}
	a.sessionStore = sessionStore

	return &a, nil
}

// SessionStore fetches the API service session store
func (api *API) SessionStore() sessions.Store {
	return api.sessionStore
}

// Run launches an API server
func (api *API) Run() error {
	base := api.GetBaseRouter()

	// Enable static file hosting if configured
	if api.options.StaticDir != "" {
		staticPath := path.Clean(api.options.StaticDir)
		base = base.Middleware(web.StaticMiddleware(staticPath))
		log.Printf("Serving static content from: '%s'", staticPath)
	}

	// Enable endpoint logging if specified
	if api.options.LogEndpoints {
		base = base.Middleware(web.LoggerMiddleware)
	}

	// Setup handlers
	var h http.Handler = base
	h = security.CORS(h, api.options)
	h = security.CSP(h, api.options)

	// Create server instance
	var server servers.Handler
	switch api.options.Mode {
	case options.ModeHTTP:
		server = servers.NewHTTP(api.options, h)
	case options.ModeLambda:
		server = servers.NewLambda(api.options, h)
	default:
		log.Errorf("Unhandled mode: '%s'", api.options.Mode)
		return errors.New("unhandled server mode")
	}
	api.server = server

	api.server.Run()

	return nil
}

// Close closes an API server (if bound)
func (api *API) Close() {
	api.server.Close()
}
