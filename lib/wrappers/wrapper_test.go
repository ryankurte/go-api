package wrappers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockFunc func(ctx APICtx, test map[string]string, h http.Header) (map[string]string, int, http.Header, error)

type apiReq struct {
	Query string
}

type Test struct {
	name string
	fn   interface{}
	err  error
}

type Input struct {
	V string
}

type APICtx struct {
}

func TestWrappers(t *testing.T) {
	tests := []Test{
		{
			"Wraps context input",
			func(ctx APICtx) (Input, error) {
				return Input{V: "test"}, nil
			},
			nil,
		}, {
			"Wraps context + object inputs",
			func(ctx APICtx, test Input) (Input, error) {
				return test, nil
			},
			nil,
		}, {
			"Wraps context + object + header inputs",
			func(ctx APICtx, test Input, h http.Header) (Input, error) {
				return test, nil
			},
			nil,
		}, {
			"Wraps object + error outputs",
			func(ctx APICtx, test Input) (Input, error) {
				return test, nil
			},
			nil,
		}, {
			"Wraps object + status + header outputs",
			func(ctx APICtx, test Input) (Input, int, error) {
				return test, http.StatusOK, nil
			},
			nil,
		}, {
			"Wraps object + status + header + error outputs",
			func(ctx APICtx, test Input) (Input, int, http.Header, error) {
				return test, http.StatusOK, http.Header{}, nil
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h, err := BuildEndpoint("post", test.fn)
			if test.err != nil {
				require.NotNil(t, err)
			} else {
				require.EqualValues(t, test.err, err)
			}

			ctx := APICtx{}
			r := Input{V: "test"}
			b, err := json.Marshal(r)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			require.Nil(t, err)
			resp := httptest.NewRecorder()

			h(ctx, resp, req)
		})
	}
}
