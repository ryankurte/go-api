package wrappers

import ()

// HandlerC Request object, headers, status code | internal error
type HandlerC func(req interface{}) (resp interface{}, err interface{})
