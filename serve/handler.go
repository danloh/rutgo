// Copyright (c) 2013-2016 Antoine Imbert
// The MIT License: https://github.com/ant0ine/go-json-rest/blob/master/LICENSE

// handler, middleware

package serve

import (
	"net/http"
)

// Handle jus like http.HandlerFunc,
// add  third parameter for the values of wildcards (variables).
// no wrap
type Handle func(http.ResponseWriter, *http.Request, Params)

// HandlerFunc defines the handler function.
// using the inherited ResponseWriter and Request, Wrapped
// and add a third param
type HandlerFunc func(ResponseWriter, *Request, Params)

// Middleware defines the interface to wrap a HandlerFunc, like decorator
type Middleware interface {
	MiddlewareFunc(handler HandlerFunc) HandlerFunc
}

// MiddleFunc is is to MiddleWare just what http.HandlerFunc is to http.Handler
type MiddleFunc func(handler HandlerFunc) HandlerFunc

// MiddlewareFunc makes MiddleFunc implement the Middleware interface.
func (mf MiddleFunc) MiddlewareFunc(handler HandlerFunc) HandlerFunc {
	return mf(handler)
}

// WrapMiddlewares to wrap a set of middlewares
// calls the MiddlewareFunc methods in the reverse order
func WrapMiddlewares(middlewares []Middleware, handler HandlerFunc) HandlerFunc {
	wrapped := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i].MiddlewareFunc(wrapped)
	}
	return wrapped
}

// AdapterFunc Handle the transition serve.HandlerFunc TO serve.Handle
func AdapterFunc(handler HandlerFunc) Handle {

	return func(origW http.ResponseWriter, origReq *http.Request, ps Params) {

		// instantiate the rest objects
		request := &Request{
			origReq,
			nil,
			map[string]interface{}{},
		}

		writer := &responseWriter{
			origW,
			false,
		}

		// call the wrapped handler
		handler(writer, request, ps)
	}
}
