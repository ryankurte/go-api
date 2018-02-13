package wrappers

import (
	"net/http"

	"github.com/gorilla/schema"

	"github.com/ryankurte/go-api-server/lib/formats"
)

const ContentTypeKey = "content-type"
const AcceptKey = "accept"

func decodeRequest(method string, req *http.Request, input interface{}) error {
	var err error
	var decoder = schema.NewDecoder()

	// Fetch content type header
	contentType := req.Header.Get(ContentTypeKey)

	//Decode input object/params
	if method == http.MethodGet {
		// Handle get request params for get method
		err = decoder.Decode(input, req.URL.Query())
	} else {
		// Handle data in body for other methods
		err = formats.Decode(contentType, req, input)
	}

	return err
}

func encodeResponse(rw http.ResponseWriter, req *http.Request, output interface{}, status int) error {
	// Fetch accept header
	acceptType := req.Header.Get(AcceptKey)

	// Attempt encoding to specified types
	out, encodedType, err := formats.Encode(acceptType, output)
	if err != nil {
		return err
	}

	// Write output data
	rw.Header().Set(ContentTypeKey, encodedType)
	rw.WriteHeader(status)
	rw.Write([]byte(out))

	return nil
}
