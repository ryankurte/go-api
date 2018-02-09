package formats

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

const XMLResourceType string = "application/xml"

type XML struct {
}

func NewXML() XML {
	return XML{}
}

func (j XML) Encode(o interface{}) (string, error) {
	js, err := xml.Marshal(o)
	if err != nil {
		return "", fmt.Errorf("XML encoding error: %s", err)
	}
	return string(js), nil
}

func (j XML) Decode(r *http.Request, i interface{}) error {
	err := xml.NewDecoder(r.Body).Decode(i)
	if err != nil {
		return fmt.Errorf("XML decoding error: %s", err)
	}
	return nil
}
