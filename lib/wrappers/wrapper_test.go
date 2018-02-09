package wrappers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryankurte/go-api-server/lib/formats"
)

type MockFunc func(ctx context.Context, test map[string]string, h http.Header) (map[string]string, int, http.Header, error)

type apiReq struct {
	Query string `query:"testQuery"`
}

type Test struct {
	name   string
	body   string
	header http.Header
	fn     interface{}
	err    error
	resp   string
}

type Input struct {
	V string
}

func TestWrappers(t *testing.T) {
	tests := []Test{
		{
			"Valid full function",
			"{\"V\": \"test\"}",
			http.Header{ContentTypeKey: []string{formats.JSONResourceType}},
			func(ctx context.Context, test Input, h http.Header) (Input, int, http.Header, error) {
				return test, http.StatusOK, h, nil
			},
			nil,
			"{\"a\": \"b\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w, err := BuildEndpoint("get", context.Background(), test.fn)
			assert.EqualValues(t, test.err, err)
			assert.NotNil(t, w)

			req := http.Request{
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(test.body))),
				Header: test.header,
			}
			resp := httptest.NewRecorder()

			h := w.(func(interface{}, http.ResponseWriter, *http.Request))

			h(context.Background(), resp, &req)

			assert.EqualValues(t, test.resp, resp.Body.String())

		})
	}

}
