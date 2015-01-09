package gocraft

import (
	"net/url"

	"github.com/gocraft/web"
)

type Request struct {
	*web.Request
}

func (r Request) QueryArgs() (url.Values, error) {
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func (r Request) PathArgs() map[string]string {
}

func (r Request) BodyJson() map[string]interface{} {
}

func (r Request) BodyForm() url.Values {
}
