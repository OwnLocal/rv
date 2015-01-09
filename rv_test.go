package rv_test

import (
	"net/url"

	. "github.com/pib/rv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BasicRequest", func() {
	var req BasicRequest

	BeforeEach(func() {
		req = BasicRequest{}
	})

	Context("With an empty request", func() {

		Describe("QueryArgs", func() {
			It("returns an empty url.Values", func() {
				Expect(req.QueryArgs()).To(Equal(url.Values{}))
			})
		})

		Describe("PathArgs", func() {
			It("returns an empty map[string]string", func() {
				Expect(req.PathArgs()).To(Equal(map[string]string(nil)))
			})
		})

		Describe("BodyJson", func() {
			It("returns an empty map[string]interface{}", func() {
				Expect(req.BodyJson()).To(Equal(map[string]interface{}(nil)))
			})
		})

		Describe("BodyForm", func() {
			It("returns an empty url.Values", func() {
				Expect(req.BodyForm()).To(Equal(url.Values{}))
			})
		})

	})

	Context("With Query specified", func() {

		BeforeEach(func() {
			req.Query = "foo=bar"
		})

		Describe("QueryArgs", func() {
			It("returns the parsed query values", func() {
				Expect(req.QueryArgs()).To(Equal(url.Values{"foo": []string{"bar"}}))
			})
		})

	})

	Context("With PathArgs specified", func() {

		BeforeEach(func() {
			req.Path = map[string]string{"foo": "bar"}
		})

		Describe("PathArgs", func() {
			It("returns the specified path arguments", func() {
				Expect(req.PathArgs()).To(Equal(map[string]string{"foo": "bar"}))
			})
		})

	})

	Context("With a valid JSON body", func() {
		BeforeEach(func() {
			req.Body = `{"foo": "bar"}`
		})

		Describe("BodyJson", func() {
			It("returns the parsed JSON", func() {
				Expect(req.BodyJson()).To(Equal(map[string]interface{}{"foo": "bar"}))
			})
		})

		Describe("BodyForm", func() {
			It("does not return an error", func() {
				// This is a bit weird because it will return a url.Values with a key of `{"foo": "bar"}`
				_, err := req.BodyForm()
				Expect(err).To(Succeed())
			})
		})
	})

	Context("With a valid FORM body", func() {
		BeforeEach(func() {
			req.Body = `one=two`
		})

		Describe("BodyJson", func() {
			It("return a JSON parse error", func() {
				_, err := req.BodyJson()
				Expect(err).To(MatchError("invalid character 'o' looking for beginning of value"))
			})
		})

		Describe("BodyForm", func() {
			It("does not return an error", func() {
				Expect(req.BodyForm()).To(Equal(url.Values{"one": []string{"two"}}))
			})
		})
	})

	Context("With an invalid FORM body", func() {
		BeforeEach(func() {
			req.Body = `%ZZ`
		})

		Describe("BodyForm", func() {
			It("returns an EscapeError", func() {
				_, err := req.BodyForm()
				Expect(err).To(Equal(url.EscapeError("%ZZ")))
			})
		})

	})

})
