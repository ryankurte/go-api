package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ryankurte/go-api-server/lib"
	"github.com/ryankurte/go-api-server/lib/options"
)

// AppConfig Application configuration object
type AppConfig struct {
	options.Base
}

// AppContext Application Context object
// Route handlers are called against this
type AppContext struct {
	Mock string
}

// Request Input structure for parsing
type Request struct {
	Message string `valid:"ascii,required"`
	Option  string `valid:"ascii,optional"`
}

// Response Output structure for parsing
type Response struct {
	Message string `valid:"ascii,required"`
}

// FakeEndpoint AppContext Endpoint handler function
func (c *AppContext) FakeEndpoint(i Request) (Response, error) {
	o := Response{i.Message}

	log.Printf("APP Endpoint context: %+v in: %+v out: %+v\n", c, i, o)

	return o, nil
}

// APIContext sub context
type APIContext struct {
	*AppContext
}

// FakeEndpoint APIContext Endpoint handler function
func (c *APIContext) FakeEndpoint(i Request) (Response, error) {
	o := Response{i.Message}

	log.Printf("API Endpoint context: %+v in: %+v out: %+v\n", c, i, o)

	return o, nil
}

func main() {
	// Load application config
	o := AppConfig{}
	err := options.Parse(&o)
	if err != nil {
		os.Exit(0)
	}

	// Create API instance with base context
	ctx := AppContext{"Whoop whoop"}
	api, err := api.New(ctx, &o.Base)
	if err != nil {
		log.Print(err)
		os.Exit(-2)
	}

	// Register logging plugin
	//api.RegisterPlugin(plugins.NewLogPlugin())

	// Register static middleware
	//api.Middleware(web.StaticMiddleware("./static", web.StaticOption{IndexFile: "index.html"}))

	// Attach base endpoint
	api.RegisterEndpoint("/", "POST", (*AppContext).FakeEndpoint)
	api.RegisterEndpoint("/", "GET", (*AppContext).FakeEndpoint)

	// Create subrouter with sub context and register endpoint
	apiCtx := APIContext{}
	sr := api.Subrouter(apiCtx, "/api")
	sr.RegisterEndpoint("/test", "POST", (*APIContext).FakeEndpoint)

	// Start API server
	api.Run()
}
