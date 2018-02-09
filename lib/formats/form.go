package formats

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

const FormResourceType string = "application/x-www-form-urlencoded"

type Form struct {
}

func NewForm() Form {
	return Form{}
}

var formDecoder = schema.NewDecoder()

func (j Form) Encode(o interface{}) (string, error) {
	return "", fmt.Errorf("FORM encoding not supported")
}

func (j Form) Decode(r *http.Request, i interface{}) error {

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("FORM parsing error", err)
	}

	if err := formDecoder.Decode(i, r.PostForm); err != nil {
		return fmt.Errorf("FORM decoding error: %s", err)
	}

	return nil
}
