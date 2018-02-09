package formats

import (
	"reflect"
	"testing"
)

func TestMain(t *testing.T) {

	// Test accept header
	var testStr = `text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c`

	t.Run("Parse accept headers", func(t *testing.T) {
		expected := []string{"text/html", "text/x-c", "text/x-dvi", "text/plain"}

		res, err := ParseAcceptHeader(testStr)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Accept headers did not match (received %+v expected %+v)",
				res, expected)
		}

	})

}
