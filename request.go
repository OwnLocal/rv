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
	return ParseJsonBody(strings.NewReader(r.Body))
}
func (r *BasicRequest) BodyForm() (url.Values, error) {
	return url.ParseQuery(r.Body)
}

func ParseJsonBody(body io.Reader) (map[string]interface{}, error) {
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
