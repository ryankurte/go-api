package wrappers

import (
	"testing"
)

type apiReq struct {
	Query string `query:"testQuery"`
}

type apiResp struct {
	StatusCode int `api:"status"`
}

func TestWrappers(t *testing.T) {

}
