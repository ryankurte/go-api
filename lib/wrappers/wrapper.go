package wrappers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/asaskevich/govalidator"
)

// HTTPHandler is a standard http endpoint handler for binding into a http mux
type HTTPHandler func(ctx interface{}, rw http.ResponseWriter, req *http.Request)

// ErrorHandler type for handling errors in wrapped functions or encoders/decoders
type ErrorHandler func(ctx interface{}, rw http.ResponseWriter, req *http.Request, format string, args ...interface{})

// ValidateHandler provides structure field validation
type ValidateHandler func(s interface{}) (bool, error)

// Decoder handles decoding of a request body into the provided interface based on the provided type
type Decoder func(method string, req *http.Request, input interface{}) error

// Encoder handles the encoding of a provided interface into an accepted form
type Encoder func(rw http.ResponseWriter, req *http.Request, output interface{}, status int) error

// DefaultErrorHandler ErrorHandler used if no error handling argument is passed to BuildEndpoint
var DefaultErrorHandler = func(ctx interface{}, rw http.ResponseWriter, req *http.Request, format string, args ...interface{}) {
	rw.WriteHeader(http.StatusInternalServerError)
	msg := fmt.Sprintf(format, args)
	log.Println(msg)
	rw.Write([]byte(msg))
}

// DefaultValidateHandler ValidateHandler used if no validation handling argument is passed to BuildEndpoint
var DefaultValidateHandler = govalidator.ValidateStruct

// DefaultDecoder Decoder used if no decoder argument is passed to BuildEndpoint
var DefaultDecoder Decoder = decodeRequest

// DefaultEncoder Encoder used if no encoder argument is passed to BuildEndpoint
var DefaultEncoder Encoder = encodeResponse

// BuildEndpoint Build and return and endpoint handler for the provided function and method
// Supports handler functions with (i InputType), (i InputType, h http.Header) or (ctx interface{}, i InputType, http.header) input parameters
// and (OutputType, error), (OutputType, int, error) or (OutputType, int, http.Header, error) output parameters where int is a http.Status code.
// args may include an ErrorHandler to override the DefaultError handler.
func BuildEndpoint(method string, fn interface{}, args ...interface{}) (HTTPHandler, error) {

	// Validate function prior to binding
	err := validateFn(fn)
	if err != nil {
		return nil, err
	}

	// Generate a wrapper function for binding
	w := generateWrapper(method, fn)

	return w, nil
}

func validateFn(fn interface{}) error {
	vf := reflect.ValueOf(fn)
	ftype := vf.Type()

	// Validate function meets one of the supported specifications
	if ftype.Kind() != reflect.Func {
		return fmt.Errorf("Method '%s' should be type `%s` but got `%s` (%+V)", ftype.Name(), reflect.Func, ftype.Kind(), fn)
	}

	argCount := ftype.NumIn()
	if argCount < 1 || argCount > 3 {
		return fmt.Errorf("Function %s invalid input parameter count", ftype.Name())
	}
	returnCount := ftype.NumOut()
	if returnCount < 1 || returnCount > 4 {
		return fmt.Errorf("Function %s invalid output parameter count", ftype.Name())
	}

	if returnCount > 2 && ftype.Out(1) != reflect.TypeOf(int(0)) {
		return fmt.Errorf("Function %s second output parameter should be of type 'int' not '%s'", ftype.Name(), ftype.Out(1).Name())
	}

	if returnCount > 3 && ftype.Out(2) != reflect.TypeOf(http.Header{}) {
		return fmt.Errorf("Function %s third output parameter should be of type 'http.Header' not '%s'", ftype.Name(), ftype.Out(2).Name())
	}

	if ftype.Out(returnCount-1) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("Function %s final output parameter (%d) is %+v however should be of Error type", ftype.Name(), returnCount-1, ftype.Out(1))
	}

	// Parse input and output types
	inputType, outputType := GetTypes(fn)

	var a interface{}
	if inputType == reflect.TypeOf(a) {
		return fmt.Errorf("Function %s inputType may not be interface{}", ftype.Name())
	}
	if outputType == reflect.TypeOf(a) {
		return fmt.Errorf("Function %s outputType may not be interface{}", ftype.Name())
	}

	return nil
}

// GetTypes fetches the associated input and output types for a supported handler function
func GetTypes(fn interface{}) (input, output reflect.Type) {
	vf := reflect.ValueOf(fn)
	ftype := vf.Type()

	// Parse input and output types
	var inputType reflect.Type
	if ftype.NumIn() == 1 {
		inputType = ftype.In(0)
	} else {
		inputType = ftype.In(1)
	}
	outputType := ftype.Out(0)

	return inputType, outputType
}

// Generate a gocraft or gorilla/mux compatible endpoint wrapper function with object mapping and validation
func generateWrapper(method string, fn interface{}, args ...interface{}) HTTPHandler {
	vf := reflect.ValueOf(fn)
	numIn := vf.Type().NumIn()
	numOut := vf.Type().NumOut()

	// Parse input and output types
	inputType, _ := GetTypes(fn)

	// Process varadic arguments
	errorHandler := DefaultErrorHandler
	validateHander := DefaultValidateHandler
	decoder, encoder := DefaultDecoder, DefaultEncoder
	for _, a := range args {
		switch a := a.(type) {
		// Bind error handler argument if present
		case ErrorHandler:
			errorHandler = a
		case ValidateHandler:
			validateHander = a
		case Decoder:
			decoder = a
		case Encoder:
			encoder = a
		}
	}

	return func(ctxIn interface{}, rw http.ResponseWriter, req *http.Request) {
		var err error

		// Coerce context and input type
		input := reflect.New(inputType)
		err = decoder(method, req, input.Interface())
		if err != nil {
			errorHandler(ctxIn, rw, req, "Data decoding error %s", err)
			return
		}

		// Validate input fields
		ok, err := validateHander(input.Interface())
		if err != nil {
			errorHandler(ctxIn, rw, req, "Input data validation error %s", err)
			return
		}
		if !ok {
			errorHandler(ctxIn, rw, req, "Input data validation failed")
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

		// Validate output fields
		ok, err = validateHander(output)
		if err != nil {
			errorHandler(ctxIn, rw, req, "Output data validation error %s", err)
			return
		}
		if !ok {
			errorHandler(ctxIn, rw, req, "Output data validation failed")
			return
		}

		// Encode outputs
		err = encoder(rw, req, output, statusCode)
		if err != nil {
			errorHandler(ctxIn, rw, req, "Data encoding error %s", err)
			return
		}
	}
}
