package formats

import (
	"fmt"
	"net/http"
)

// Formatter Defines a formatter interface
type Formatter interface {
	Encode(o interface{}) (string, error)
	Decode(r *http.Request, i interface{}) error
}

// Default formatting adaptors
var formatters = map[string]Formatter{
	JSONResourceType: NewJSON(),
	XMLResourceType:  NewXML(),
	YAMLResourceType: NewYAML(),
	FormResourceType: NewForm(),
}

// DefaultResponseFormatter defines the default encoding in the absence of headers
var DefaultResponseFormatter = JSONResourceType

// DefaultRequestFormatter defines the default decoding in the absence of headers
var DefaultRequestFormatter = JSONResourceType

// Decode Generic decode function (uses http accept header)
func Decode(t string, r *http.Request, i interface{}) error {
	if t == "" {
		t = DefaultRequestFormatter
	}

	if i == nil {
		return fmt.Errorf("Unable to decode nil type")
	}

	// Find formatter and decode
	if f, ok := formatters[t]; ok {
		return f.Decode(r, i)
	}

	// No formatter found
	return fmt.Errorf("No decoder found matching type: %s", t)
}

// Encode Generic encode function
func Encode(accepts string, i interface{}) (string, string, error) {
	// Parse accept types
	types := ParseAcceptHeader(accepts)

	// No accept types, fall back to default
	if len(types) == 0 {
		if f, ok := formatters[DefaultResponseFormatter]; ok {
			s, e := f.Encode(i)
			return s, DefaultResponseFormatter, e
		}
	}

	// List of types, locate and use the first matching formatter
	for _, t := range types {
		if f, ok := formatters[t]; ok {
			s, e := f.Encode(i)
			return s, t, e
		}
	}

	// No formatter found
	return "", "", fmt.Errorf("No encoder found matching types: %s", types)
}

// BindFormatter Bind an alternate formatter implementation
func BindFormatter(t string, f Formatter) {
	formatters[t] = f
}

// RemoveFormatter Remove a formatter implementation
func RemoveFormatter(t string) {
	delete(formatters, t)
}
