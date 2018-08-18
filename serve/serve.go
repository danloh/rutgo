// as so far, credit to go-json-rest: https://github.com/ant0ine/go-json-rest

package serve

import (
	"log"
	"net/http"
)

// Serve defines a stack of Middlewares and an App.
type Serve struct {
	*Router
	stack []Middleware // like global middware
}

// NewServe makes a new Serve object. The Middleware stack is empty.
func NewServe() *Serve {
	return &Serve{
		Router: NewRouter(),
		stack:  []Middleware{},
	}
}

// Use pushes one or multiple middlewares to the stack
func (srv *Serve) Use(middlewares ...Middleware) {
	srv.stack = append(srv.stack, middlewares...)
	srv.Router.GlobalMiddle = srv.stack
}

// Run to start serve
func (srv *Serve) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, srv.Router)
	log.Fatal(err)
	return
}
