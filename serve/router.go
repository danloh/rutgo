package serve

import "net/http"

// RouterGroup is configure router
// path, an array of handlers (middleware).
type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	Serve    *Serve
	root     bool
}

// HandlerFunc to handle request
type HandlerFunc func(http.ResponseWriter, http.Request)

// HandlersChain is array of HandlerFunc
type HandlersChain []HandlerFunc

// Use to add middleware to the Router
func (group *RouterGroup) Use(middleware ...HandlerFunc) *RouterGroup {
	group.Handlers = append(group.Handlers, middleware...)
	return group
}
