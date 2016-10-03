package gocraft_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	. "github.com/OwnLocal/rv/gocraft"
	"github.com/gocraft/web"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request", func() {
	var req Request

	BeforeEach(func() {
		req = Request{Request: &web.Request{Request: &http.Request{URL: &url.URL{RawQuery: "a=b&c=d"}}}}
	})

	Context("With a valid URL querystring", func() {
		Describe("QueryArgs", func() {
			It("returns the parsed query args", func() {
				Expect(req.QueryArgs()).To(Equal(url.Values{"a": []string{"b"}, "c": []string{"d"}}))
			})
		})
	})

	Context("With no PathParams", func() {
		Describe("PathArgs", func() {
			It("returns an empty map[string]string", func() {
				Expect(req.PathArgs()).To(Equal(map[string]string(nil)))
			})
		})
	})

	Context("With PathParams", func() {
		BeforeEach(func() {
			req.Request.PathParams = map[string]string{"foo": "bar"}
		})

		Describe("PathArgs", func() {
			It("returns the specified params", func() {
				Expect(req.PathArgs()).To(Equal(map[string]string{"foo": "bar"}))
			})
		})
	})

	Context("With an nil body", func() {
		Describe("BodyJSON", func() {
			It("returns a nil map", func() {
				Expect(req.BodyJSON()).To(Equal(map[string]interface{}(nil)))
			})
		})

		Describe("BodyForm", func() {
			It("returns an empty url.Values", func() {
				Expect(req.BodyForm()).To(Equal(url.Values(nil)))
			})
		})
	})

	Context("With a valid JSON body", func() {
		BeforeEach(func() {
			req.Request.Body = ioutil.NopCloser(strings.NewReader(`{"foo": "bar"}`))
		})

		Describe("BodyJSON", func() {
			It("returns the parsed JSON", func() {
				Expect(req.BodyJSON()).To(Equal(map[string]interface{}{"foo": "bar"}))
			})
		})

		Describe("BodyForm", func() {
			It("does not return an error", func() {
				// The result will be meaningless, but it won't be an error...
				_, err := req.BodyForm()
				Expect(err).To(Succeed())
			})
		})
	})

	Context("With a valid form body", func() {
		BeforeEach(func() {
			req.Request.Body = ioutil.NopCloser(strings.NewReader(`foo=bar&one=two`))
		})

		Describe("BodyJSON", func() {
			It("returns the a JSON parse error", func() {
				_, err := req.BodyJSON()
				Expect(err).To(MatchError("invalid character 'o' in literal false (expecting 'a')"))
			})
		})

		Describe("BodyForm", func() {
			It("returns the parsed form", func() {
				Expect(req.BodyForm()).To(Equal(url.Values{"foo": []string{"bar"}, "one": []string{"two"}}))
			})
		})
	})

	Context("When the body has already been read by the BodyJSON method", func() {
		BeforeEach(func() {
			req.BodyJSON()
		})

		Describe("BodyJSON", func() {
			It("returns an body-already-read error", func() {
				json, err := req.BodyJSON()
				Expect(json).To(Equal(map[string]interface{}(nil)))
				Expect(err).To(MatchError("body already read"))
			})
		})

		Describe("BodyForm", func() {
			It("returns an body-already-read error", func() {
				form, err := req.BodyForm()
				Expect(form).To(Equal(url.Values(nil)))
				Expect(err).To(MatchError("body already read"))
			})
		})
	})

})
