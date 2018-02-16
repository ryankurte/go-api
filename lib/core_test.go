package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryankurte/go-api-server/lib/options"
)

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
func (c *AppContext) FakeEndpoint(i Request, h http.Header) (Response, error) {
	o := Response{i.Message}

	return o, nil
}

// APIContext sub context
type APIContext struct {
	*AppContext
}

// FakeEndpoint APIContext Endpoint handler function
func (c *APIContext) FakeEndpoint(i Request, h http.Header) (Response, error) {
	o := Response{i.Message}

	return o, nil
}

func TestCore(t *testing.T) {
	o := options.Base{}
	o.Mode = options.ModeHTTP
	o.BindAddress = "127.0.0.1"
	o.Port = "9002"
	o.NoTLS = true

	addr := fmt.Sprintf("http://%s:%s/", o.BindAddress, o.Port)

	// Create API instance with base context
	ctx := AppContext{"Whoop whoop"}
	api, err := New(ctx, &o)
	assert.Nil(t, err)

	// Attach base endpoint
	err = api.RegisterEndpoint("/", "GET", (*AppContext).FakeEndpoint)
	assert.Nil(t, err)

	err = api.RegisterEndpoint("/", "POST", (*AppContext).FakeEndpoint)
	assert.Nil(t, err)

	go api.Run()
	defer api.Close()

	client := http.DefaultClient

	resp, err := client.Get(addr + "?message=test")
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Nil(t, err)

	assert.Equal(t, string(body), "{\"Message\":\"test\"}")

}
