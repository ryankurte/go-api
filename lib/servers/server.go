package servers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gcontext "github.com/gorilla/context"

	"github.com/ryankurte/go-api/lib/options"
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
	bindAddress := fmt.Sprintf("%s:%s", s.options.BindAddress, s.options.Port)

	s.server = http.Server{Addr: bindAddress, Handler: contextHandler}

	s.logger.Infof("Starting http server at %s (bind: %s)", s.options.ExternalAddress, bindAddress)

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

func (s *HTTP) Start() {
	go s.Run()
}

// Close exits a server instance
func (s *HTTP) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
	cancel()
}
