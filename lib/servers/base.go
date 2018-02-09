package servers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ryankurte/go-api-server/lib/options"
)

// Handler run interface
type Handler interface {
	Run()
	Close()
}

// Base handler type
type Base struct {
	name    string
	options *options.Base
	logger  log.FieldLogger
	handler http.Handler
}

// NewBase creates a new base handler
// TODO: should this take an http.handler OR a more processed (ie. pre-read) object?
func NewBase(name string, handler http.Handler, options *options.Base) Base {
	b := Base{
		name:    name,
		options: options,
		handler: handler,
		logger:  log.New().WithField("module", name),
	}

	return b
}
