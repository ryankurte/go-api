package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ryankurte/go-api-server/lib/options"
)

// Lambda is an AWS Lambda based http handler
type Lambda struct {
	Base
}

// NewLambda creates a new lambda handler with the given http handler func
func NewLambda(o options.Options, handler http.Handler) Lambda {
	return Lambda{
		Base: NewBase("lambda", handler, o),
	}
}

func (h *Lambda) mapAPIGatewayRequest(req events.APIGatewayProxyRequest) (*http.Request, error) {
	headers := make(http.Header)
	for k, v := range req.Headers {
		headers[k] = strings.Split(v, ";")
	}
	url, err := url.Parse(req.Path)
	if err != nil {
		return nil, err
	}
	body := ioutil.NopCloser(bytes.NewBuffer([]byte(req.Body)))

	return &http.Request{
		Method:        req.HTTPMethod,
		URL:           url,
		Header:        headers,
		Body:          body,
		ContentLength: int64(len(req.Body)),
	}, nil
}

func (h *Lambda) mapAPIGatewayResponse(resp *httptest.ResponseRecorder) (*events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	for k, v := range resp.Header() {
		headers[k] = strings.Join(v, ";")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      resp.Code,
		Headers:         headers,
		Body:            resp.Body.String(),
		IsBase64Encoded: false,
	}, nil
}

type apiGWResp events.APIGatewayProxyResponse

func (lr *apiGWResp) Header() http.Header {
	return lr.Header()
}

func (lr *apiGWResp) Write(body []byte) (int, error) {
	lr.Body += string(body)
	return len(body), nil
}

func (lr *apiGWResp) WriteHeader(status int) {
	lr.StatusCode = status
}

var internalError = events.APIGatewayProxyResponse{
	StatusCode: http.StatusInternalServerError,
	Body:       "Lambda wrapper error",
}

func (h *Lambda) handle(ctx context.Context, gwReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger := h.logger.WithField("request-id", gwReq.RequestContext.RequestID)

	req, err := h.mapAPIGatewayRequest(gwReq)
	if err != nil {
		logger.Errorf("Mapping api request (%s)", err)
		return internalError, err
	}

	resp := &httptest.ResponseRecorder{}
	h.handler.ServeHTTP(resp, req)

	gwResp, err := h.mapAPIGatewayResponse(resp)
	if err != nil {
		logger.Errorf("Mapping api response (%s)", err)
		return internalError, err
	}

	return *gwResp, nil
}

// Run starts a lambda server instance
func (h *Lambda) Run() {
	lambda.Start(h.handle)
}
