package options

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/jessevdk/go-flags"
)

// Base are base API server options
type Base struct {
	Mode string `short:"m" long:"mode" description:"Server mode" choice:"http" choice:"lambda" default:"http"`

	BindAddress string `short:"b" long:"address" description:"Address to bind API server" default:"0.0.0.0"`
	Port        string `short:"p" long:"port" description:"Port on which to bind API server" default:"10001"`

	ExternalAddress string `short:"e" long:"external-address" description:"External address for connection to server" default:"localhost"`

	NoTLS   bool   `long:"no-tls" description:"Disable TLS"`
	TLSCert string `short:"c" long:"tls-cert" description:"TLS certificate file"`
	TLSKey  string `short:"k" long:"tls-key" description:"TLS key file"`

	StaticDir string `short:"s" long:"static-dir" description:"Directory to serve static content from (if specified)"`

	AllowedOrigins []string `long:"allowed-origins" description:"Allowed origins for CORS headers (defaults to external address or bind address)"`

	CookieSecret string `long:"cookie-secret" description:"Secret for session cookie encryption (defaults to a random key)"`

	LogEndpoints bool `long:"log-endpoints" description:"Enable endpoint logging"`
}

// Server mode constants
const (
	ModeLambda = "lambda"
	ModeHTTP   = "http"
)

// Parse parses command line options
func Parse(i interface{}) error {
	_, err := flags.Parse(i)
	return err
}

// GenerateSecret Helper to generate a default secret to use
func GenerateSecret(len int) (string, error) {
	data := make([]byte, len)
	n, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	if n != len {
		return "", fmt.Errorf("Config: RNG failed")
	}

	return base64.URLEncoding.EncodeToString(data), nil
}
