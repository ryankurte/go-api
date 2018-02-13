package api

import (
	"errors"
	"github.com/gocraft/web"

	log "github.com/sirupsen/logrus"

	"github.com/ryankurte/go-api-server/lib/options"
	"github.com/ryankurte/go-api-server/lib/router"
	"github.com/ryankurte/go-api-server/lib/servers"
)

// API is a core API server instance
type API struct {
	router.Router
	options *options.Base
	logger  log.FieldLogger
	server  servers.Handler
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

	// Create server instance
	var server servers.Handler
	switch o.Mode {
	case options.ModeHTTP:
		server = servers.NewHTTP(o, base)
	case options.ModeLambda:
		server = servers.NewLambda(o, base)
	default:
		log.Errorf("Unhandled mode: '%s'", o.Mode)
		return nil, errors.New("unhandled server mode")
	}

	if server == nil {
		log.Errorf("Error creating server: %s", err)
		return nil, errors.New("error creating server")
	}

	a.server = server

	return &a, nil
}

// Run launches an API server
func (api *API) Run() {
	api.server.Run()
}

// Close closes an API server (if bound)
func (api *API) Close() {
	api.server.Close()
}
