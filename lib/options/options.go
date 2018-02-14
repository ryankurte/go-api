package options

import (
	"github.com/jessevdk/go-flags"
)

// Base are base API server options
type Base struct {
	Mode string `short:"m" long:"mode" description:"Server mode" choice:"http" choice:"lambda" default:"http"`

	BindAddress string `short:"b" long:"address" description:"Address to bind API server" default:"0.0.0.0"`
	Port        string `short:"p" long:"port" description:"Port on which to bind API server" default:"10001"`

	NoTLS   bool   `long:"no-tls" description:"Disable TLS"`
	TLSCert string `short:"c" long:"tls-cert" description:"TLS certificate file"`
	TLSKey  string `short:"k" long:"tls-key" description:"TLS key file"`
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
