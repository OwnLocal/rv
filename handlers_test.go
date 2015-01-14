package rv_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ownlocal/rv"
)

var _ = Describe("Validators", func() {
	var req *rv.BasicRequest
	var field *rv.Field

	BeforeEach(func() {
		req = &rv.BasicRequest{
			Path:  map[string]string{"foo": "bar"},
			Query: "bar=baz&blah=flah"}
		field = new(rv.Field)
	})

	Describe("SourceFieldHandler", func() {
		Describe("NewSourceFieldHandler", func() {
			It("rejects invalid sources", func() {
				_, err := rv.NewSourceFieldHandler([]string{"foo.bar"})
				Expect(err).To(HaveOccurred())
			})

			It("accepts valid sources", func() {
				Expect(rv.NewSourceFieldHandler([]string{"path.foo"})).To(Equal(rv.SourceFieldHandler{Source: rv.PATH, Field: "foo"}))
				Expect(rv.NewSourceFieldHandler([]string{"query.foo"})).To(Equal(rv.SourceFieldHandler{Source: rv.QUERY, Field: "foo"}))
				Expect(rv.NewSourceFieldHandler([]string{"json.foo"})).To(Equal(rv.SourceFieldHandler{Source: rv.JSON, Field: "foo"}))
				Expect(rv.NewSourceFieldHandler([]string{"form.foo"})).To(Equal(rv.SourceFieldHandler{Source: rv.FORM, Field: "foo"}))
			})
		})

		Describe("Run", func() {

			It("properly pulls fields from path arguments", func() {
				rv.SourceFieldHandler{Source: rv.PATH, Field: "foo"}.Run(req, field)
				Expect(field.Value).To(Equal("bar"))
			})

			It("properly pulls fields from query arguments", func() {
				rv.SourceFieldHandler{Source: rv.QUERY, Field: "blah"}.Run(req, field)
				Expect(field.Value).To(Equal("flah"))
			})

			It("properly pulls fields from JSON body arguments", func() {
				req.Body = `{"one": 2}`
				rv.SourceFieldHandler{Source: rv.JSON, Field: "one"}.Run(req, field)
				Expect(field.Value).To(Equal(2.0))
			})

			It("properly pulls fields from form body arguments", func() {
				req.Body = `one=two&three=four`
				rv.SourceFieldHandler{Source: rv.FORM, Field: "three"}.Run(req, field)
				Expect(field.Value).To(Equal("four"))
			})

		})
	})
})
