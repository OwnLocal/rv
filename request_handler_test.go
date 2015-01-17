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

		It("generates a options handler for tags that specify options", func() {
			rh, err := rv.NewRequestHandler(struct {
				Foo string `rv:"query.foo options=one,two,three"`
			}{})
			y := struct{}{}
			expected := map[string]rv.FieldHandlers{
				"Foo": rv.FieldHandlers{
					rv.SourceFieldHandler{Source: rv.QUERY, Field: "foo"},
					rv.TypeHandler{Type: "string"},
					rv.OptionsHandler{Options: map[string]struct{}{"one": y, "two": y, "three": y}},
				},
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(rh.Fields).To(Equal(expected))
		})

		It("generates a required handler for tags that specify required", func() {
			rh, err := rv.NewRequestHandler(struct {
				Foo int `rv:"query.foo required=true range=1,2"`
			}{})
			expected := map[string]rv.FieldHandlers{
				"Foo": rv.FieldHandlers{
					rv.SourceFieldHandler{Source: rv.QUERY, Field: "foo"},
					rv.TypeHandler{Type: "int"},
					rv.RangeHandler{Start: "1", End: "2"},
					rv.RequiredHandler{Required: true},
				},
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(rh.Fields).To(Equal(expected))
		})

		It("generates a list handler for tags on list types", func() {
			rh, err := rv.NewRequestHandler(struct {
				Foo []string `rv:"query.foo options=one,two,three default=one,two"`
			}{})
			y := struct{}{}
			expected := map[string]rv.FieldHandlers{
				"Foo": rv.FieldHandlers{
					rv.SourceFieldHandler{Source: rv.QUERY, Field: "foo"},
					rv.DefaultHandler{Default: []string{"one", "two"}},
					rv.ListHandler{SubHandlers: rv.FieldHandlers{
						rv.TypeHandler{Type: "string"},
						rv.OptionsHandler{Options: map[string]struct{}{"one": y, "two": y, "three": y}},
					}},
				},
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(rh.Fields).To(Equal(expected))
		})
	})

	Describe("Run", func() {
		type testStruct struct {
			Foo []string `rv:"query.foo options=one,two,three default=one"`
		}

		var (
			rh  *rv.RequestHandler
			err error
			req *rv.BasicRequest
		)

		BeforeEach(func() {
			rh, err = rv.NewRequestHandler(testStruct{})
			req = &rv.BasicRequest{}
		})

		It("returns an error if not supplied a pointer to the type of struct it expects", func() {
			ts := testStruct{}
			err, fieldErrs := rh.Run(req, ts)

			Expect(fieldErrs).To(BeEmpty())
			Expect(err).To(MatchError("Expected *rv_test.testStruct, got rv_test.testStruct"))
		})

		It("fills in the struct values if there are no errors", func() {
			ts := testStruct{}
			err, fieldErrs := rh.Run(req, &ts)
			Expect(err).NotTo(HaveOccurred())
			Expect(fieldErrs).To(BeEmpty())
			Expect(ts.Foo).To(Equal([]string{"one"}))
		})
	})

	Describe("Bind", func() {
		type testStruct struct {
			Foo int `rv:"query.i"`
		}
		type testContext struct {
			Foo  string
			Args testStruct
		}
		type badContext struct {
			Foo  string
			Args map[int]string
		}
		var (
			rh  *rv.RequestHandler
			tc  *testContext
			bc  *badContext
			req *rv.BasicRequest
			err error
		)

		BeforeEach(func() {
			tc = &testContext{}
			bc = &badContext{}
			rh, err = rv.NewRequestHandler(testStruct{})
			req = &rv.BasicRequest{Query: "i=42"}
		})

		Context("When the provided struct contains an attribute of the correct type", func() {

			It("finds the matching struct and fills it in", func() {
				err, fieldErrs := rh.Bind(req, tc)
				Expect(err).NotTo(HaveOccurred())
				Expect(fieldErrs).To(BeEmpty())
				Expect(tc.Args.Foo).To(Equal(42))
				// Call again to make sure caching doesn't explode
				req.Query = "i=1"
				err, fieldErrs = rh.Bind(req, tc)
				Expect(err).NotTo(HaveOccurred())
				Expect(fieldErrs).To(BeEmpty())
				Expect(tc.Args.Foo).To(Equal(1))
			})

			It("returns an error when the container doesn't have the expected struct", func() {
				err, fieldErrs := rh.Bind(req, bc)
				Expect(err).To(HaveOccurred())
				Expect(fieldErrs).To(BeEmpty())
			})

		})

	})
})
