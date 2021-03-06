// Copyright (c) 2013-2016 Antoine Imbert
// The MIT License: https://github.com/ant0ine/go-json-rest/blob/master/LICENSE

// inherit from http.Request
// to add additional method

package serve

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
	ErrJSONPayloadEmpty = errors.New("JSON payload is empty")
)

// Request inherits from http.Request, and provides additional methods.
type Request struct {
	*http.Request

	// Map of parameters that have been matched in the URL Path.
	PathParams map[string]string

	// Environment used by middlewares to communicate.
	Env map[string]interface{}
}

// PathParam provides a convenient access to the PathParams map.
func (r *Request) PathParam(name string) string {
	return r.PathParams[name]
}

// DecodeJSONPayload reads the request body and decodes the JSON using json.Unmarshal.
// store data to v
func (r *Request) DecodeJSONPayload(v interface{}) error {
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return ErrJSONPayloadEmpty
	}
	err = json.Unmarshal(content, v) // store json-data to v
	if err != nil {
		return err
	}
	return nil
}

// BaseURL returns a new URL object with the Host and Scheme taken from the request.
// (without the trailing slash in the host)
func (r *Request) BaseURL() *url.URL {
	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	// HTTP sometimes gives the default scheme as HTTP even when used with TLS
	// Check if TLS is not nil and given back https scheme
	if scheme == "http" && r.TLS != nil {
		scheme = "https"
	}

	host := r.Host
	if len(host) > 0 && host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

// URLFor returns the URL object from UriBase with the Path set to path, and the query
// string built with queryParams.
func (r *Request) URLFor(path string, queryParams map[string][]string) *url.URL {
	baseURL := r.BaseURL()
	baseURL.Path = path
	if queryParams != nil {
		query := url.Values{}
		for k, v := range queryParams {
			for _, vv := range v {
				query.Add(k, vv)
			}
		}
		baseURL.RawQuery = query.Encode()
	}
	return baseURL
}

// CorsInfo contains the CORS request info derived from a s.Request.
type CorsInfo struct {
	IsCors      bool
	IsPreflight bool
	Origin      string
	OriginURL   *url.URL

	// The header value is converted to uppercase to avoid common mistakes.
	AccessControlRequestMethod string

	// The header values are normalized with http.CanonicalHeaderKey.
	AccessControlRequestHeaders []string
}

// GetCorsInfo derives CorsInfo from Request.
func (r *Request) GetCorsInfo() *CorsInfo {

	origin := r.Header.Get("Origin")

	var originURL *url.URL
	var isCors bool

	if origin == "" {
		isCors = false
	} else if origin == "null" {
		isCors = true
	} else {
		var err error
		originURL, err = url.ParseRequestURI(origin)
		isCors = err == nil && r.Host != originURL.Host
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")

	reqHeaders := []string{}
	rawReqHeaders := r.Header[http.CanonicalHeaderKey("Access-Control-Request-Headers")]
	for _, rawReqHeader := range rawReqHeaders {
		if len(rawReqHeader) == 0 {
			continue
		}
		// net/http does not handle comma delimited headers for us
		for _, reqHeader := range strings.Split(rawReqHeader, ",") {
			reqHeaders = append(reqHeaders, http.CanonicalHeaderKey(strings.TrimSpace(reqHeader)))
		}
	}

	isPreflight := isCors && r.Method == "OPTIONS" && reqMethod != ""

	return &CorsInfo{
		IsCors:                      isCors,
		IsPreflight:                 isPreflight,
		Origin:                      origin,
		OriginURL:                   originURL,
		AccessControlRequestMethod:  strings.ToUpper(reqMethod),
		AccessControlRequestHeaders: reqHeaders,
	}
}
