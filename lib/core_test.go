package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func (c *APIContext) ContextEndpoint(i Request) (APIContext, error) {
	return *c, nil
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
	require.Nil(t, err)

	// Attach base endpoint
	err = api.RegisterEndpoint("/", "GET", (*AppContext).FakeEndpoint)
	require.Nil(t, err)

	err = api.RegisterEndpoint("/", "POST", (*AppContext).FakeEndpoint)
	require.Nil(t, err)

	go api.Run()
	defer api.Close()

	client := http.DefaultClient

	t.Run("Get with query params", func(t *testing.T) {
		resp, err := client.Get(addr + "?message=test")
		require.Nil(t, err)
		require.NotNil(t, resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.Nil(t, err)

		require.Equal(t, string(body), "{\"Message\":\"test\"}")
	})

	t.Run("Post form", func(t *testing.T) {
		v := url.Values{"message": []string{"test"}}
		resp, err := client.PostForm(addr, v)
		require.Nil(t, err)
		require.NotNil(t, resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.Nil(t, err)

		require.Equal(t, string(body), "{\"Message\":\"test\"}")
	})

	t.Run("Post JSON", func(t *testing.T) {
		r := Request{Message: "test"}
		b, _ := json.Marshal(r)
		resp, err := client.Post(addr, "application/json", bytes.NewReader(b))
		require.Nil(t, err)
		require.NotNil(t, resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.Nil(t, err)

		require.Equal(t, string(body), "{\"Message\":\"test\"}")
	})

	t.Run("Post invalid content-type", func(t *testing.T) {
		r := Request{Message: "test"}
		b, _ := json.Marshal(r)
		resp, err := client.Post(addr, "application/cats", bytes.NewReader(b))
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.NotNil(t, resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.Nil(t, err)

		assert.Contains(t, string(body), "No decoder found matching type")
	})

	t.Run("Get YAML response", func(t *testing.T) {
		req, err := http.NewRequest("GET", addr+"?message=test", nil)
		require.Nil(t, err)
		req.Header.Set("accept", "application/yaml")
		resp, err := client.Do(req)
		require.Nil(t, err)
		require.NotNil(t, resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.Nil(t, err)

		require.Equal(t, string(body), "message: test\n")
	})

	t.Run("Invalid accept header", func(t *testing.T) {
		req, err := http.NewRequest("GET", addr+"?message=test", nil)
		require.Nil(t, err)
		req.Header.Set("accept", "application/cats")
		resp, err := client.Do(req)
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.NotNil(t, resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.Nil(t, err)

		assert.Contains(t, string(body), "No encoder found matching type")
	})

}
