package wrappers

import (
	"log"
	"reflect"
)

type Wrapper struct {
	in  []reflect.Type
	out []reflect.Type
}

func BuildEndpoint(method string, ctx interface{}, f interface{}) interface{} {
	vf := reflect.ValueOf(f)
	ftype := vf.Type()

	if ftype.Kind() != reflect.Func {
		log.Panicf("`f` should be %s but got %s", reflect.Func, ftype.Kind())
	}
	if ftype.NumIn() != 2 {
		log.Panicf("`f` should have 2 input parameters, but it has %d", ftype.NumIn())
	}
	if ftype.NumOut() != 2 {
		log.Panicf("`f` should have 2 output parameters but it has %d", ftype.NumOut())
	}
	if ftype.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		log.Panicf("`f` second output parameter is %+v however should be of Error type", ftype.Out(1))
	}

	return nil
}
