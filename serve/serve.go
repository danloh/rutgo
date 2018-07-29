package serve

import (
	"net/http"
)

// Serve defines a stack of Middlewares and an App.
type Serve struct {
	RouterGroup
}

// New makes a new Api object. The Middleware stack is empty, and the App is nil.
func New() *Serve {
	return &Serve{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
	}
}

// Use attachs a middleware to the router
func (s *Serve) Use(middleware ...HandlerFunc) *Serve {
	s.RouterGroup.Use(middleware...)
	return s
}

// Run to start serve
func (s *Serve) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, s)
	return
}

// ServeHTTP conforms to the http.Handler interface.
func (s *Serve) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// todo
}
