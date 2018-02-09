package api

import (
	"errors"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"

	"github.com/ryankurte/go-api-server/lib/options"
	"github.com/ryankurte/go-api-server/lib/servers"
)

// API is a core API server instance
type API struct {
	options *options.Base
	logger  log.FieldLogger
	server  servers.Handler
}

// NewAPI creates a new API server
func NewAPI(o *options.Base) (*API, error) {
	var err error
	a := API{
		options: o,
		logger:  log.New().WithField("module", "core"),
	}

	// TODO: create API router
	r := mux.NewRouter()

	// Create server instance
	var server servers.Handler
	switch o.Mode {
	case options.ModeHTTP:
		server = servers.NewHTTP(o, r)
	case options.ModeLambda:
		server = servers.NewLambda(o, r)
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
