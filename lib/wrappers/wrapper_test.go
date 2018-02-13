package wrappers

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockFunc func(ctx context.Context, test map[string]string, h http.Header) (map[string]string, int, http.Header, error)

type apiReq struct {
	Query string `query:"testQuery"`
}

type Test struct {
	name string
	fn   interface{}
	err  error
}

type Input struct {
	V string
}

func TestWrappers(t *testing.T) {
	tests := []Test{
		{
			"Wraps object input",
			func(test Input) (Input, int, http.Header, error) {
				return test, http.StatusOK, http.Header{}, nil
			},
			nil,
		}, {
			"Wraps object + header inputs",
			func(test Input, h http.Header) (Input, int, http.Header, error) {
				return test, http.StatusOK, h, nil
			},
			nil,
		}, {
			"Wraps context + object + header inputs",
			func(ctx context.Context, test Input, h http.Header) (Input, int, http.Header, error) {
				return test, http.StatusOK, h, nil
			},
			nil,
		}, {
			"Wraps object + error outputs",
			func(test Input) (Input, error) {
				return test, nil
			},
			nil,
		}, {
			"Wraps object + status + header outputs",
			func(test Input) (Input, int, error) {
				return test, http.StatusOK, nil
			},
			nil,
		}, {
			"Wraps object + status + header + error outputs",
			func(test Input) (Input, int, http.Header, error) {
				return test, http.StatusOK, http.Header{}, nil
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := BuildEndpoint("post", test.fn)
			if test.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.EqualValues(t, test.err, err)
			}
		})
	}
}
