package rv_test

import (
	"github.com/ownlocal/rv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestHandler", func() {

	Describe("NewRequestHandler", func() {

		It("generates an empty RequestHandler for structs with no rv tags", func() {
			rh, _ := rv.NewRequestHandler(struct{ Foo string }{})
			Expect(rh.Fields).To(BeEmpty())
		})

		It("generates a source and type handlers for fields with just a source specified", func() {
			rh, err := rv.NewRequestHandler(struct {
				Foo string `rv:"query.foo"`
			}{})
			expected := map[string]rv.FieldHandlers{
				"Foo": rv.FieldHandlers{
					rv.SourceFieldHandler{Source: rv.QUERY, Field: "foo"},
					rv.TypeHandler{Type: "string"},
				},
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(rh.Fields).To(Equal(expected))
		})

	})

})
