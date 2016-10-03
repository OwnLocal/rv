package goji

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/OwnLocal/rv"
	"goji.io/pattern"
)

// Request holds a http.Request and knows how to pull Goji path params from the request context.
type Request struct {
	*http.Request
	bodyRead bool
}

// QueryArgs pulls the standard query args from the request.
func (r *Request) QueryArgs() (url.Values, error) {
	vals, err := url.ParseQuery(r.Request.URL.RawQuery) // Doesn't use URL.Query because we want to see errors.
	if err != nil {
		return nil, err
	}
	return vals, nil
}

// PathArgs extracts all Goji path variables from the request context.
func (r *Request) PathArgs() (map[string]string, error) {
	if pathVars, ok := r.Context().Value(pattern.AllVariables).(map[pattern.Variable]interface{}); ok {
		pathMap := make(map[string]string, len(pathVars))
		for k, v := range pathVars {
			pathMap[string(k)] = v.(string)
		}
		return pathMap, nil
	}
	return nil, nil
}

// BodyJSON parses and returns a map of values from the body.
func (r *Request) BodyJSON() (map[string]interface{}, error) {
	if r.bodyRead {
		return nil, errors.New("body already read")
	}

	r.bodyRead = true
	return rv.ParseJSONBody(r.Request.Body)
}

// BodyForm parses and returns any values from a body form.
func (r *Request) BodyForm() (url.Values, error) {
	if r.bodyRead {
		return nil, errors.New("body already read")
	}

	r.bodyRead = true

	if r.Request.Body == nil {
		return nil, nil
	}

	body, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return nil, err
	}

	return url.ParseQuery(string(body))
}

// Ensure *gocract.Request meets the rv.Request interface
var _ rv.Request = (*Request)(nil)
