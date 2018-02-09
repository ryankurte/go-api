package formats

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-yaml/yaml"
)

const yamlResourceType string = "application/yaml"

type YAML struct {
}

func NewYAML() YAML {
	return YAML{}
}

func (j YAML) Encode(o interface{}) (string, error) {
	js, err := yaml.Marshal(o)
	if err != nil {
		return "", fmt.Errorf("YAML encoding error: %s", err)
	}
	return string(js), nil
}

func (j YAML) Decode(r *http.Request, i interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("YAML decoding error: %s", err)
	}

	err = yaml.Unmarshal(data, i)
	if err != nil {
		return fmt.Errorf("YAML decoding error: %s", err)
	}
	return nil
}
