package serve

import (
	"errors"
	"net/http"
	"net/url"
	"rutgo/serve/trie"
	"strings"
)

type router struct {
	Routes []*Route

	disableTrieCompression bool
	index                  map[*Route]int
	trie                   *trie.Trie
}

// MakeRouter returns the router app. Given a set of Routes, it dispatches the request to the
// HandlerFunc of the first route that matches. The order of the Routes matters.
func MakeRouter(routes ...*Route) (App, error) {
	r := &router{
		Routes: routes,
	}
	err := r.start()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Handle the REST routing and run the user code.
func (rt *router) AppFunc() HandlerFunc {
	return func(writer ResponseWriter, request *Request) {
		// find the route
		route, params, pathMatched := rt.findRouteFromURL(request.Method, request.URL)
		if route == nil {

			if pathMatched {
				// no route found, but path was matched: 405 Method Not Allowed
				Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// no route found, the path was not matched: 404 Not Found
			NotFound(writer, request)
			return
		}

		// a route was found, set the PathParams
		request.PathParams = params

		// run the user code
		handler := route.Func
		handler(writer, request)
	}
}

// This is run for each new request, perf is important.
func escapedPath(urlObj *url.URL) string {
	// the escape method of url.URL should be public
	// that would avoid this split.
	parts := strings.SplitN(urlObj.RequestURI(), "?", 2)
	return parts[0]
}

var preEscape = strings.NewReplacer("*", "__SPLAT_PLACEHOLDER__", "#", "__RELAXED_PLACEHOLDER__")

var postEscape = strings.NewReplacer("__SPLAT_PLACEHOLDER__", "*", "__RELAXED_PLACEHOLDER__", "#")

// This is run at init time only.
func escapedPathExp(pathExp string) (string, error) {

	// PathExp validation
	if pathExp == "" {
		return "", errors.New("empty PathExp")
	}
	if pathExp[0] != '/' {
		return "", errors.New("PathExp must start with /")
	}
	if strings.Contains(pathExp, "?") {
		return "", errors.New("PathExp must not contain the query string")
	}

	// Get the right escaping
	// XXX a bit hacky

	pathExp = preEscape.Replace(pathExp)

	urlObj, err := url.Parse(pathExp)
	if err != nil {
		return "", err
	}

	// get the same escaping as find requests
	pathExp = urlObj.RequestURI()

	pathExp = postEscape.Replace(pathExp)

	return pathExp, nil
}

// This validates the Routes and prepares the Trie data structure.
// It must be called once the Routes are defined and before trying to find Routes.
// The order matters, if multiple Routes match, the first defined will be used.
func (rt *router) start() error {

	rt.trie = trie.New()
	rt.index = map[*Route]int{}

	for i, route := range rt.Routes {

		// work with the PathExp urlencoded.
		pathExp, err := escapedPathExp(route.PathExp)
		if err != nil {
			return err
		}

		// insert in the Trie
		err = rt.trie.AddRoute(
			strings.ToUpper(route.HTTPMethod), // work with the HttpMethod in uppercase
			pathExp,
			route,
		)
		if err != nil {
			return err
		}

		// index
		rt.index[route] = i
	}

	if rt.disableTrieCompression == false {
		rt.trie.Compress()
	}

	return nil
}

// return the result that has the route defined the earliest
func (rt *router) ofFirstDefinedRoute(matches []*trie.Match) *trie.Match {
	minIndex := -1
	var bestMatch *trie.Match

	for _, result := range matches {
		route := result.Route.(*Route)
		routeIndex := rt.index[route]
		if minIndex == -1 || routeIndex < minIndex {
			minIndex = routeIndex
			bestMatch = result
		}
	}

	return bestMatch
}

// Return the first matching Route and the corresponding parameters for a given URL object.
func (rt *router) findRouteFromURL(httpMethod string, urlObj *url.URL) (*Route, map[string]string, bool) {

	// lookup the routes in the Trie
	matches, pathMatched := rt.trie.FindRoutesAndPathMatched(
		strings.ToUpper(httpMethod), // work with the httpMethod in uppercase
		escapedPath(urlObj),         // work with the path urlencoded
	)

	// short cuts
	if len(matches) == 0 {
		// no route found
		return nil, nil, pathMatched
	}

	if len(matches) == 1 {
		// one route found
		return matches[0].Route.(*Route), matches[0].Params, pathMatched
	}

	// multiple routes found, pick the first defined
	result := rt.ofFirstDefinedRoute(matches)
	return result.Route.(*Route), result.Params, pathMatched
}

// Parse the url string (complete or just the path) and return the first matching Route and the corresponding parameters.
func (rt *router) findRoute(httpMethod, urlStr string) (*Route, map[string]string, bool, error) {

	// parse the url
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, false, err
	}

	route, params, pathMatched := rt.findRouteFromURL(httpMethod, urlObj)
	return route, params, pathMatched, nil
}

// Route pattern to handler
type Route struct {

	// Any HTTP method. It will be used as uppercase to avoid common mistakes.
	HTTPMethod string

	// A string like "/resource/:id.json".
	// Placeholders supported are:
	// :paramName that matches any char to the first '/' or '.'
	// #paramName that matches any char to the first '/'
	// *paramName that matches everything to the end of the string
	// (placeholder names must be unique per PathExp)
	PathExp string

	// Code that will be executed when this route is taken.
	Func HandlerFunc
}

// MakePath generates the path corresponding to this Route and the provided path parameters.
// This is used for reverse route resolution.
func (route *Route) MakePath(pathParams map[string]string) string {
	path := route.PathExp
	for paramName, paramValue := range pathParams {
		paramPlaceholder := ":" + paramName
		relaxedPlaceholder := "#" + paramName
		splatPlaceholder := "*" + paramName
		r := strings.NewReplacer(paramPlaceholder, paramValue, splatPlaceholder, paramValue, relaxedPlaceholder, paramValue)
		path = r.Replace(path)
	}
	return path
}

// Head is a shortcut method that instantiates a HEAD route. See the Route object the parameters definitions.
// Equivalent to &Route{"HEAD", pathExp, handlerFunc}
func Head(pathExp string, handlerFunc HandlerFunc) *Route {
	return &Route{
		HTTPMethod: "HEAD",
		PathExp:    pathExp,
		Func:       handlerFunc,
	}
}

// Get is a shortcut method that instantiates a GET route. See the Route object the parameters definitions.
// Equivalent to &Route{"GET", pathExp, handlerFunc}
func Get(pathExp string, handlerFunc HandlerFunc) *Route {
	return &Route{
		HTTPMethod: "GET",
		PathExp:    pathExp,
		Func:       handlerFunc,
	}
}

// Post is a shortcut method that instantiates a POST route. See the Route object the parameters definitions.
// Equivalent to &Route{"POST", pathExp, handlerFunc}
func Post(pathExp string, handlerFunc HandlerFunc) *Route {
	return &Route{
		HTTPMethod: "POST",
		PathExp:    pathExp,
		Func:       handlerFunc,
	}
}

// Put is a shortcut method that instantiates a PUT route.  See the Route object the parameters definitions.
// Equivalent to &Route{"PUT", pathExp, handlerFunc}
func Put(pathExp string, handlerFunc HandlerFunc) *Route {
	return &Route{
		HTTPMethod: "PUT",
		PathExp:    pathExp,
		Func:       handlerFunc,
	}
}

// Patch is a shortcut method that instantiates a PATCH route.  See the Route object the parameters definitions.
// Equivalent to &Route{"PATCH", pathExp, handlerFunc}
func Patch(pathExp string, handlerFunc HandlerFunc) *Route {
	return &Route{
		HTTPMethod: "PATCH",
		PathExp:    pathExp,
		Func:       handlerFunc,
	}
}

// Delete is a shortcut method that instantiates a DELETE route. Equivalent to &Route{"DELETE", pathExp, handlerFunc}
func Delete(pathExp string, handlerFunc HandlerFunc) *Route {
	return &Route{
		HTTPMethod: "DELETE",
		PathExp:    pathExp,
		Func:       handlerFunc,
	}
}

// Options is a shortcut method that instantiates an OPTIONS route.  See the Route object the parameters definitions.
// Equivalent to &Route{"OPTIONS", pathExp, handlerFunc}
func Options(pathExp string, handlerFunc HandlerFunc) *Route {
	return &Route{
		HTTPMethod: "OPTIONS",
		PathExp:    pathExp,
		Func:       handlerFunc,
	}
}
