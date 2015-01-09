// Package rv fills and validates values into structs from HTTP
// requests via reflection.
package rv

import (
	"encoding/json"
	"net/url"
)

// Request is the interface to implement to allow req to read values
// from a HTTP request.
type Request interface {
	QueryArgs() (url.Values, error)
	PathArgs() (map[string]string, error)
	BodyJson() (map[string]interface{}, error)
	BodyForm() (url.Values, error)
}

// BasicRequest implements the Request interface and can be used for
// testing or parsing requests from unsupported request types.
type BasicRequest struct {
	Query string
	Path  map[string]string
	Body  string
}

func (r *BasicRequest) QueryArgs() (url.Values, error) {
	return url.ParseQuery(r.Query)
}
func (r *BasicRequest) PathArgs() (map[string]string, error) {
	return r.Path, nil
}
func (r *BasicRequest) BodyJson() (map[string]interface{}, error) {
	if r.Body == "" {
		return nil, nil
	}

	parsed := make(map[string]interface{})
	if err := json.Unmarshal([]byte(r.Body), &parsed); err != nil {
		return nil, err
	}
	return parsed, nil
}
func (r *BasicRequest) BodyForm() (url.Values, error) {
	return url.ParseQuery(r.Body)
}
