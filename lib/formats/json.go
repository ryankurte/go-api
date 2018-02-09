package formats

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const jsonResourceType string = "application/json"

type JSON struct {
}

func NewJSON() JSON {
	return JSON{}
}

func (j JSON) Encode(o interface{}) (string, error) {
	js, err := json.Marshal(o)
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %s", err)
	}
	return string(js), nil
}

func (j JSON) Decode(r *http.Request, i interface{}) error {
	err := json.NewDecoder(r.Body).Decode(i)
	if err != nil {
		return fmt.Errorf("JSON decoding error: %s", err)
	}
	return nil
}
