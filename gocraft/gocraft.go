package gocraft

import (
	"errors"
	"io/ioutil"
	"net/url"

	"github.com/gocraft/web"
	"github.com/ownlocal/rv"
)

type Request struct {
	Request  *web.Request
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
	return r.Request.PathParams, nil
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
