package gocraft

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

// BindMiddleware creates an rv.RequestHandler for the specified field
// type and returns a middleware which finds a field of that type on
// the context and binds the values to that field via the
// RequestHandler.
func BindMiddleware(field interface{}, errorWriter ...func(web.ResponseWriter, error, map[string]rv.Field)) func(
	interface{}, web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {

	if len(errorWriter) < 1 {
		errorWriter = append(errorWriter, ErrorWriter)
	}

	argHandler, err := rv.NewRequestHandler(field)
	if err != nil {
		panic("Unable to create RequestHandler: " + err.Error())
	}

	return func(ctx interface{}, rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		err, fieldErrors := argHandler.Bind(&Request{Request: r}, ctx)
		if err != nil || len(fieldErrors) > 0 {
			errorWriter[0](rw, err, fieldErrors)
		} else {
			next(rw, r)
		}
	}
}

func ErrorWriter(rw web.ResponseWriter, argErr error, fieldErrors map[string]rv.Field) {
	if argErr != nil {
		panic(argErr)
	}

	rw.WriteHeader(http.StatusBadRequest)
	for name, field := range fieldErrors {
		for _, err := range field.Errors {
			fmt.Fprintln(rw, name, err)
		}
	}
}
