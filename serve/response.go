// inherit from http.Response
// to add additional method

package serve

import (
	"net/http"
)

// A ResponseWriter interface dedicated to JSON HTTP response.
type ResponseWriter interface {

	// Identical to the http.ResponseWriter interface
	Header() http.Header

	// Use EncodeJson to generate the payload
	WriteJson(value interface{}) error

	// Encode the data structure to JSON
	EncodeJson(value interface{}) ([]byte, error)

	WriteHeader(int)
}
