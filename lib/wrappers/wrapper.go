package wrappers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
)

func BuildEndpoint(method string, ctx interface{}, f interface{}) (interface{}, error) {
	vf := reflect.ValueOf(f)
	ftype := vf.Type()

	if ftype.Kind() != reflect.Func {
		return nil, fmt.Errorf("Method '%s' should be type %s but got %s", ftype.Name(), reflect.Func, ftype.Kind())
	}

	argCount := ftype.NumIn()
	if argCount < 2 || argCount > 3 {
		return nil, fmt.Errorf("Function %s invalid input parameter count", ftype.Name())
	}
	returnCount := ftype.NumOut()
	if returnCount < 1 || returnCount > 4 {
		return nil, fmt.Errorf("Function %s invalid output parameter count", ftype.Name())
	}

	if returnCount > 2 && ftype.Out(1) != reflect.TypeOf(int(0)) {
		return nil, fmt.Errorf("Function %s second output parameter should be of type 'int' not '%s'", ftype.Name(), ftype.Out(1).Name())
	}

	if returnCount > 3 && ftype.Out(2) != reflect.TypeOf(http.Header{}) {
		return nil, fmt.Errorf("Function %s third output parameter should be of type 'http.Header' not '%s'", ftype.Name(), ftype.Out(2).Name())
	}

	if ftype.Out(returnCount-1) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, fmt.Errorf("Function %s final output parameter (%d) is %+v however should be of Error type", ftype.Name(), returnCount-1, ftype.Out(1))
	}

	// Parse input and output types
	inputType := ftype.In(1)
	outputType := ftype.Out(0)

	// Generate a wrapper function for binding
	w := generateWrapper(method, vf, inputType, outputType)

	return w, nil
}

// Generate a gocraft compatible endpoint wrapper function with object mapping and validation
func generateWrapper(method string, vf reflect.Value, inputType reflect.Type, outputType reflect.Type) interface{} {
	numIn := vf.Type().NumIn()
	numOut := vf.Type().NumOut()

	return func(ctxIn interface{}, rw http.ResponseWriter, req *http.Request) {
		var err error

		// Coerce context and input type
		input := reflect.New(inputType)
		err = decodeRequest(method, req, input.Interface())
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "Data decoding error %s", err)
			return
		}

		// Validate input fields
		_, err = govalidator.ValidateStruct(input.Interface())
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "Data validation error %s", err)
			return
		}

		// Generate input values
		var inputs []reflect.Value
		switch numIn {
		case 1:
			inputs = []reflect.Value{input.Elem()}
		case 2:
			inputs = []reflect.Value{reflect.ValueOf(ctxIn), input.Elem()}
		case 3:
			inputs = []reflect.Value{reflect.ValueOf(ctxIn), input.Elem(), reflect.ValueOf(req.Header)}
		default:
			fmt.Printf("INVALID INPUT COUNT %d", numIn)
			return
		}

		fmt.Printf("Inputs: %+v", inputs)

		// Call reflected function
		outputs := vf.Call(inputs)

		// Parse function call errors
		err, _ = reflect.ValueOf(outputs[numOut-1]).Interface().(error)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "Internal Server Error", err)
			return
		}

		// Fetch status code if returned
		statusCode := http.StatusOK
		if numOut > 2 {
			statusCode = outputs[1].Interface().(int)
		}

		// Fetch header if returned
		if numOut > 3 {
			respHeaders := outputs[2].Interface().(http.Header)
			for k, v := range respHeaders {
				rw.Header().Set(k, strings.Join(v, ";"))
			}
		}

		// Coerce and write output type
		output := outputs[0].Interface()
		err = encodeResponse(req, output, statusCode, rw)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "Data encoding error %s", err)
			return
		}
	}
}
