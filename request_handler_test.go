package rv_test

import (
	"github.com/ownlocal/rv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SimpleRequest struct {
	SimpleArg string `rv:"query.arg"`
}

type JsonRequest struct {
	Type  []string `rv:"path.type options=posts,pages,categories"`
	Size  int      `rv:"query.size range=1,50 default=10"`
	From  int      `rv:"query.from range=0,1000 default=0"`
	Query string   `rv:"json.q default=*"`
}

type FormRequest struct {
	Type   []string `rv:"path.type options=posts,pages,categories"`
	Size   int      `rv:"query.size range=1,50 default=10"`
	From   int      `rv:"query.from range=0,1000 default=0"`
	Filter []string `rv:"form.filter options=all,some,none default=all,some"`
}

var _ = Describe("RequestHandler", func() {

	Describe("NewRequestHandler", func() {

		It("generates an empty RequestHandler for structs with no rv tags", func() {
			rh, _ := rv.NewRequestHandler(struct{ Foo string }{})
			Expect(rh.Fields).To(BeEmpty())
		})

		It("generates a source and type handler for fields with just a source specified", func() {
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

		It("generates a default handler for tags that specify a default", func() {
			rh, err := rv.NewRequestHandler(struct {
				Foo string `rv:"query.foo default=bar"`
			}{})
			expected := map[string]rv.FieldHandlers{
				"Foo": rv.FieldHandlers{
					rv.SourceFieldHandler{Source: rv.QUERY, Field: "foo"},
					rv.DefaultHandler{Default: "bar"},
					rv.TypeHandler{Type: "string"},
				},
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(rh.Fields).To(Equal(expected))
		})

		It("generates a range handler for tags that specify a range", func() {
			rh, err := rv.NewRequestHandler(struct {
				Foo string `rv:"query.foo range=1,10"`
			}{})
			expected := map[string]rv.FieldHandlers{
				"Foo": rv.FieldHandlers{
					rv.SourceFieldHandler{Source: rv.QUERY, Field: "foo"},
					rv.TypeHandler{Type: "string"},
					rv.RangeHandler{Start: "1", End: "10"},
				},
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(rh.Fields).To(Equal(expected))
		})

	})

})
