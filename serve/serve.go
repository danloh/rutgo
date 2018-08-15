// as so far, credit to go-json-rest: https://github.com/ant0ine/go-json-rest

package serve

import (
	"net/http"
)

// Serve defines a stack of Middlewares and an App.
type Serve struct {
	stack []Middleware
	app   App
}

// NewServe makes a new Api object. The Middleware stack is empty, and the App is nil.
func NewServe() *Serve {
	return &Serve{
		stack: []Middleware{},
		app:   nil,
	}
}

// Use pushes one or multiple middlewares to the stack for middlewares
// maintained in the Api object.
func (srv *Serve) Use(middlewares ...Middleware) {
	srv.stack = append(srv.stack, middlewares...)
}

// SetApp sets the App in the Api object.
func (srv *Serve) SetApp(app App) {
	srv.app = app
}

// MakeHandler wraps all the Middlewares of the stack and the App together, and returns an
// http.Handler ready to be used. If the Middleware stack is empty the App is used directly. If the
// App is nil, a HandlerFunc that does nothing is used instead.
func (srv *Serve) MakeHandler() http.Handler {
	var appFunc HandlerFunc
	if srv.app != nil {
		appFunc = srv.app.AppFunc()
	} else {
		appFunc = func(w ResponseWriter, r *Request) {}
	}
	return http.HandlerFunc(
		adapterFunc(
			WrapMiddlewares(srv.stack, appFunc),
		),
	)
}

// Run to start serve
func (srv *Serve) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, srv.MakeHandler())
	return
}

// ServeHTTP conforms to the http.Handler interface.
func (srv *Serve) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// todo
}
