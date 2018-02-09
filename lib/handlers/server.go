package handlers

import (
	"context"
	"net/http"
	"time"

	gcontext "github.com/gorilla/context"

	"github.com/ryankurte/go-api-server/lib/options"
)

// Server is an HTTP server based http handler
type Server struct {
	Base
	server http.Server
}

// NewServer creates a new HTTP server with the provided options
func NewServer(o options.Options, h http.Handler) Server {
	return Server{
		Base: NewBase("server", h, o),
	}
}

// Run starts a server instance (this only returns on error or exit)
func (s *Server) Run() error {
	var err error

	contextHandler := gcontext.ClearHandler(s.handler)

	s.server = http.Server{Addr: s.options.BindAddress, Handler: contextHandler}

	if s.options.TLSCert != "" && s.options.TLSKey != "" {
		err = s.server.ListenAndServeTLS(s.options.TLSCert, s.options.TLSKey)
	} else {
		err = s.server.ListenAndServe()
	}

	return err
}

// Close exits a server instance
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
	cancel()
}
