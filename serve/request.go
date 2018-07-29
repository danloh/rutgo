// inherit from http.Request
// to add additional method

package serve

import (
	"net/http"
)

// Request inherits from http.Request, and more methods.
type Request struct {
	*http.Request
	// to add
}
