package goji_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"goji.io/pat"

	goji "goji.io"

	. "github.com/OwnLocal/rv/goji"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request", func() {
	var mux *goji.Mux
	var req *http.Request
	var ctx context.Context
	var res *httptest.ResponseRecorder
	var rvReq Request

	handler := func(w http.ResponseWriter, r *http.Request) {
		rvReq = Request{Request: r}
		ctx = r.Context()
	}

	mux = goji.NewMux()
	mux.HandleFunc(pat.Get("/foo"), handler)
	mux.HandleFunc(pat.Get("/foo/:foo"), handler)

	BeforeEach(func() {
		ctx = nil
		req = httptest.NewRequest("PATCH", "/NOT-SET-SET-ME", nil)
		res = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		mux.ServeHTTP(res, req.WithContext(context.Background()))
	})

	Context("With a valid URL querystring", func() {
		Describe("QueryArgs", func() {
			BeforeEach(func() {
				req = httptest.NewRequest("GET", "/foo?a=b&c=d", nil)
			})
			It("returns the parsed query args", func() {
				Expect(rvReq.QueryArgs()).To(Equal(url.Values{"a": []string{"b"}, "c": []string{"d"}}))
			})
		})
	})

	Context("With no PathParams", func() {
		Describe("PathArgs", func() {
			It("returns an empty map[string]string", func() {
				Expect(rvReq.PathArgs()).To(Equal(map[string]string(nil)))
			})
		})
	})

	Context("With PathParams", func() {
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "/foo/bar", nil)
		})

		Describe("PathArgs", func() {
			It("returns the specified params", func() {
				Expect(rvReq.PathArgs()).To(Equal(map[string]string{"foo": "bar"}))
			})
		})
	})

	Context("With an nil body", func() {
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "/foo", nil)
		})
		Describe("BodyJson", func() {
			It("returns a nil map", func() {
				Expect(rvReq.BodyJson()).To(Equal(map[string]interface{}(nil)))
			})
		})

		Describe("BodyForm", func() {
			It("returns an empty url.Values", func() {
				Expect(rvReq.BodyForm()).To(Equal(url.Values{}))
			})
		})
	})

	Context("With a valid JSON body", func() {
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "/foo", ioutil.NopCloser(strings.NewReader(`{"foo": "bar"}`)))
		})

		Describe("BodyJson", func() {
			It("returns the parsed JSON", func() {
				Expect(rvReq.BodyJson()).To(Equal(map[string]interface{}{"foo": "bar"}))
			})
		})

		Describe("BodyForm", func() {
			It("does not return an error", func() {
				// The result will be meaningless, but it won't be an error...
				_, err := rvReq.BodyForm()
				Expect(err).To(Succeed())
			})
		})
	})

	Context("With a valid form body", func() {
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "/foo", ioutil.NopCloser(strings.NewReader(`foo=bar&one=two`)))
		})

		Describe("BodyJson", func() {
			It("returns the a JSON parse error", func() {
				_, err := rvReq.BodyJson()
				Expect(err).To(MatchError("invalid character 'o' in literal false (expecting 'a')"))
			})
		})

		Describe("BodyForm", func() {
			It("returns the parsed form", func() {
				Expect(rvReq.BodyForm()).To(Equal(url.Values{"foo": []string{"bar"}, "one": []string{"two"}}))
			})
		})
	})

	Context("When the body has already been read by the BodyJson method", func() {
		BeforeEach(func() {
			rvReq.BodyJson()
		})

		Describe("BodyJson", func() {
			It("returns an body-already-read error", func() {
				json, err := rvReq.BodyJson()
				Expect(json).To(Equal(map[string]interface{}(nil)))
				Expect(err).To(MatchError("body already read"))
			})
		})

		Describe("BodyForm", func() {
			It("returns an body-already-read error", func() {
				form, err := rvReq.BodyForm()
				Expect(form).To(Equal(url.Values(nil)))
				Expect(err).To(MatchError("body already read"))
			})
		})
	})

})
