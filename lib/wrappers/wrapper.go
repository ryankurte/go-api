package wrappers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
)

// HTTPHandler is a standard http endpoint handler for binding into a http mux
type HTTPHandler func(ctx interface{}, rw http.ResponseWriter, req *http.Request)

// ErrorHandler type for handling errors in wrapped functions or encoders/decoders
type ErrorHandler func(ctx interface{}, rw http.ResponseWriter, req *http.Request, format string, args ...interface{})

// DefaultError ErrorHandler used if no error handling argument is passed to BuildEndpoint
func DefaultError(ctx interface{}, rw http.ResponseWriter, req *http.Request, format string, args ...interface{}) {
	rw.WriteHeader(http.StatusInternalServerError)
	msg := fmt.Sprintf(format, args)
	log.Println(msg)
	rw.Write([]byte(msg))
}

// BuildEndpoint Build and return and endpoint handler for the provided function and method
// Supports handler functions with (i InputType), (i InputType, h http.Header) or (ctx interface{}, i InputType, http.header) input parameters
// and (OutputType, error), (OutputType, int, error) or (OutputType, int, http.Header, error) output parameters where int is a http.Status code.
// args may include an ErrorHandler to override the DefaultError handler.
func BuildEndpoint(method string, fn interface{}, args ...interface{}) (HTTPHandler, error) {
	vf := reflect.ValueOf(fn)
	ftype := vf.Type()

	// Validate function meets one of the supported specifications
	if ftype.Kind() != reflect.Func {
		return nil, fmt.Errorf("Method '%s' should be type %s but got %s", ftype.Name(), reflect.Func, ftype.Kind())
	}

	argCount := ftype.NumIn()
	if argCount < 1 || argCount > 3 {
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
	var inputType reflect.Type
	if argCount == 1 {
		inputType = ftype.In(0)
	} else {
		inputType = ftype.In(1)
	}
	outputType := ftype.Out(0)

	// Generate a wrapper function for binding
	w := generateWrapper(method, vf, inputType, outputType)

	return w, nil
}

// Generate a gocraft or gorilla/mux compatible endpoint wrapper function with object mapping and validation
func generateWrapper(method string, vf reflect.Value, inputType, outputType reflect.Type, args ...interface{}) HTTPHandler {
	numIn := vf.Type().NumIn()
	numOut := vf.Type().NumOut()

	// Process varadic arguments
	errorHandler := DefaultError
	for _, a := range args {
		switch a := a.(type) {
		// Bind error handler argument if present
		case ErrorHandler:
			errorHandler = a
		}
	}

	return func(ctxIn interface{}, rw http.ResponseWriter, req *http.Request) {
		var err error

		// Coerce context and input type
		input := reflect.New(inputType)
		err = decodeRequest(method, req, input.Interface())
		if err != nil {
			errorHandler(ctxIn, rw, req, "Data decoding error %s", err)
			return
		}

		// Validate input fields
		_, err = govalidator.ValidateStruct(input.Interface())
		if err != nil {
			errorHandler(ctxIn, rw, req, "Data validation error %s", err)
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
			errorHandler(ctxIn, rw, req, "Invalid input parameter count")
			return
		}

		// Call reflected function
		outputs := vf.Call(inputs)

		// Parse function call errors
		err, _ = reflect.ValueOf(outputs[numOut-1]).Interface().(error)
		if err != nil {
			errorHandler(ctxIn, rw, req, "Internal Server Error %s", err)
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
			errorHandler(ctxIn, rw, req, "Data encoding error %s", err)
			return
		}
	}
}
