package goji

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/OwnLocal/rv"
	"goji.io/pattern"
)

type Request struct {
	*http.Request
	bodyRead bool
}

func (r *Request) QueryArgs() (url.Values, error) {
	vals, err := url.ParseQuery(r.Request.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

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

func (r *Request) BodyJson() (map[string]interface{}, error) {
	if r.bodyRead {
		return nil, errors.New("body already read")
	}

	r.bodyRead = true
	return rv.ParseJsonBody(r.Request.Body)
}

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
