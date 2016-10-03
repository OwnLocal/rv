// Package rv fills and validates values into structs from HTTP
// requests via reflection.
package rv

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

// Request is the interface to implement to allow rv to read values
// from a HTTP request.
type Request interface {
	// QueryArgs returns the query parameters from the request URL
	QueryArgs() (url.Values, error)
	// PathArgs returns any parameters from the URL path, if supported by the request
	PathArgs() (map[string]string, error)
	// BodyJSON reads the request body and parses it as JSON
	BodyJSON() (map[string]interface{}, error)
	// BodyForm reads the request body and parses it as a form encoding
	BodyForm() (url.Values, error)
}

// BasicRequest implements the Request interface and can be used for
// testing or parsing requests from unsupported request types.
type BasicRequest struct {
	Query string
	Path  map[string]string
	Body  string
}

// QueryArgs parses the Query field
func (r *BasicRequest) QueryArgs() (url.Values, error) {
	return url.ParseQuery(r.Query)
}

// PathArgs returns the Path field value
func (r *BasicRequest) PathArgs() (map[string]string, error) {
	return r.Path, nil
}

// BodyJSON parses the Body string as JSON
func (r *BasicRequest) BodyJSON() (map[string]interface{}, error) {
	return ParseJSONBody(strings.NewReader(r.Body))
}

// BodyForm parses the Body string as a form
func (r *BasicRequest) BodyForm() (url.Values, error) {
	return url.ParseQuery(r.Body)
}

// ParseJSONBody attempts to parse a JSON body from the provided io.Reader
func ParseJSONBody(body io.Reader) (map[string]interface{}, error) {
	if body == nil {
		return nil, nil
	}

	decoder := json.NewDecoder(body)

	parsed := make(map[string]interface{})

	if err := decoder.Decode(&parsed); err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return parsed, nil
}
