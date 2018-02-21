# go-api

A typed API builder project to remove a bunch of boilerplate and repetition from golang API projects.
This wires together a bunch of common dependencies, uses some reflection based magic to provide typed API functions and common error handling, and is intended to encode some best practices in API security.

## Status

[![Build Status](https://travis-ci.org/ryankurte/go-api.svg?branch=master)](https://travis-ci.org/ryankurte/go-api)
[![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ryankurte/go-api)
[![Release](https://img.shields.io/github/release/ryankurte/go-api.svg)](https://github.com/ryankurte/go-api)


## Overview

- [core](lib/) collects components and exposes the user API
- [formats](lib/formats) provide format encoding/decoding functions
- [options](lib/options) provide base api server options and option parsing
- [plugins](lib/plugins) provide plugins for meta analysis of the API implementation
- [security](lib/security) provide security extensions for the API implementation
- [servers](lib/servers) provide base server handling (ie. http server, AWS lambda)
- [wrappers](lib/wrappers) provide wrapping functions for typed api endpoints

## Usage

Install with `go get github.com/ryankurte/go-api`.

Create an application options object that inherits from `options.Base` and load with `options.Parse(&o)`.
Options are parsed using [jessevdk/go-flags](https://github.com/jessevdk/go-flags).

``` go
import (
    "github.com/ryankurte/go-api/lib/options"
)

type AppConfig struct {
    options.Base
    ...
}

...

o := AppConfig{}
err := options.Parse(&o)
if err != nil {
    os.Exit(0)
} 
```

Create a base application context with type handlers and a base `api.API` router, the attach handlers to the API router.

Handler functions support input parameters:
- `(ctx ContextType)`
- `(ctx ContextType, i InputType)`
- `(ctx ContextType, i InputType, http.header)`

And output parameters:
- `(OutputType, error)`
- `(OutputType, int, error)`
- `(OutputType, int, http.Header, error)`

Where `ContextType`, `InputType` and `OutputType` are user defined structs and int is a `http.Status` code. 
Input and output types are validated after decoding and prior to encoding using [asaskevich/govalidator](https://github.com/asaskevich/govalidator).
The underlying mux is provided by [gocraft/web](https://github.com/gocraft/web).

``` go
import (
    "github.com/ryankurte/go-api/lib/options"
)

type AppContext struct {
	...
}

func (c *AppContext) FakeEndpoint(a AppContext, i Request) (Response, error) {
	...
}

ctx := AppContext{"Whoop whoop"}
api, err := api.New(ctx, &o.Base)
if err != nil {
    log.Print(err)
    os.Exit(-1)
}

err = api.RegisterEndpoint("/", "POST", (*AppContext).FakeEndpoint)

```

You can then launch a server with `api.Run()` and exit wth `api.Close()`.

Check out [example.go](example.go) for a working example.

------

If you have any questions, comments, or suggestions, feel free to open an issue or a pull request.
