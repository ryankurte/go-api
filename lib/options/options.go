package options

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/jessevdk/go-flags"
)

// Base are base API server options
type Base struct {
	Mode            string `short:"m" long:"mode" description:"Server mode" choice:"http" choice:"lambda" default:"http"`
	BindAddress     string `short:"b" long:"address" description:"Address to bind API server" default:"0.0.0.0"`
	Port            string `short:"p" long:"port" description:"Port on which to bind API server" default:"10001"`
	ExternalAddress string `short:"e" long:"external-address" description:"External address for connection to server" default:"localhost:10001"`
	StaticDir       string `short:"s" long:"static-dir" description:"Directory to serve static content from (if specified)"`

	TLS `namespace:"tls" group:"Transport Layer Security (TLS) options"`

	CookieSecret string `long:"cookie-secret" description:"Secret for session cookie encryption (defaults to a random key)"`

	LogEndpoints bool `long:"log-endpoints" description:"Enable endpoint logging"`

	CORS `namespace:"cors" group:"Cross Origin Resource Sharing (CORS) settings"`
	CSP  `namespace:"csp" group:"Content Security Policy (CSP) settings"`
}

func (b *Base) GetExternalAddress() string {
	if b.TLS.NoTLS {
		return fmt.Sprintf("http://%s", b.ExternalAddress)
	} else {
		return fmt.Sprintf("https://%s", b.ExternalAddress)
	}
}

func (b *Base) GetBindAddress() string {
	if b.TLS.NoTLS {
		return fmt.Sprintf("http://%s:%s", b.BindAddress, b.Port)
	} else {
		return fmt.Sprintf("https://%s:%s", b.BindAddress, b.Port)
	}
}

// TLS configuration options
type TLS struct {
	TLSCert string `short:"c" long:"cert" description:"TLS certificate file"`
	TLSKey  string `short:"k" long:"key" description:"TLS key file"`
	NoTLS   bool   `long:"disable" description:"Disable TLS"`
}

// CORS configuration options
type CORS struct {
	AllowedOrigins   []string `long:"allowed-origins" description:"Allowed origins (defaults to external address or bind address)"`
	AllowedMethods   []string `long:"allowed-methods" description:"Allowed http methods" default:"GET,POST,PUT,UPDATE,OPTIONS"`
	AllowedHeaders   []string `long:"allowed-headers" description:"Allowed headers" default:"Content-Type"`
	AllowCredentials bool     `long:"allowed-credentials" description:"Allowed credentials"`
	NoCORS           bool     `long:"disable" description:"Disable CORS headers"`
}

// CSP configuration options
type CSP struct {
	ReportOnly  bool     `long:"report-only" description:"Sets CSP to report only mode"`
	DefaultSrc  []string `long:"default-src" description:"Default allowed sources" default:"'self'"`
	ScriptSrc   []string `long:"script-src" description:"Allowed script sources"`
	StyleSrc    []string `long:"style-src" description:"Allowed style sources"`
	ImgSrc      []string `long:"img-src" description:"Allowed img sources"`
	FontSrc     []string `long:"font-src" description:"Allowed font sources"`
	ChildSrc    []string `long:"child-src" description:"Allowed child sources"`
	ConnectSrc  []string `long:"connect-src" description:"Allowed connect sources"`
	FrameSrc    []string `long:"frame-src" description:"Allowed frame sources"`
	ManifestSrc []string `long:"manifest-src" description:"Allowed manifest sources"`
	MediaSrc    []string `long:"media-src" description:"Allowed media sources"`
	ObjectSrc   []string `long:"object-src" description:"Allowed object sources"`
	WorkerSrc   []string `long:"worker-src" description:"Allowed worker sources"`
	ReportTo    string   `long:"report-to" description:"ReportTo address" default:"/csp-report"`
	NoCSP       bool     `long:"disable" description:"Disable CSP headers"`
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
