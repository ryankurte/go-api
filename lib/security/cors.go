package security

import (
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/ryankurte/go-api/lib/options"
)

// CORS builds a Cross-Origin Resource Sharing (COS) handler around the provided handler
// with the specified options
func CORS(h http.Handler, o *options.Base) http.Handler {
	if len(o.AllowedOrigins) == 0 {
		o.AllowedOrigins = []string{o.ExternalAddress}
	}

	config := []handlers.CORSOption{
		handlers.AllowedOrigins(o.AllowedOrigins),
		handlers.AllowedMethods(o.AllowedMethods),
		handlers.AllowedHeaders(o.AllowedHeaders),
	}

	if o.AllowCredentials {
		config = append(config, handlers.AllowCredentials())
	}

	CORSHandler := handlers.CORS(config...)

	return CORSHandler(h)
}
