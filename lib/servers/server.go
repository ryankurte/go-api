package servers

import (
	"context"
	"net/http"
	"time"

	gcontext "github.com/gorilla/context"

	"github.com/ryankurte/go-api-server/lib/options"
)

// HTTP is an HTTP server based http handler
type HTTP struct {
	Base
	server http.Server
}

// NewHTTP creates a new HTTP server with the provided options
func NewHTTP(o *options.Base, h http.Handler) *HTTP {
	return &HTTP{
		Base: NewBase(options.ModeHTTP, h, o),
	}
}

// Run starts a server instance (this only returns on error or exit)
func (s *HTTP) Run() {
	var err error

	contextHandler := gcontext.ClearHandler(s.handler)

	s.server = http.Server{Addr: s.options.BindAddress, Handler: contextHandler}

	if s.options.NoTLS {
		s.logger.Warn("TLS IS DISABLED. USE EXTERNAL TLS TERMINATION.")
		err = s.server.ListenAndServe()
	} else if s.options.TLSCert != "" && s.options.TLSKey != "" {
		s.logger.Info("Starting http server with TLS")
		err = s.server.ListenAndServeTLS(s.options.TLSCert, s.options.TLSKey)
	} else {
		s.logger.Error("TLS enabled but missing certificate or key argument")
	}

	if err != nil {
		s.logger.Errorf("ListenAndServe error: %s", err)
	}
}

// Close exits a server instance
func (s *HTTP) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
	cancel()
}