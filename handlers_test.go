package rv_test

import (
	"fmt"
	"math"
	"strings"

	"github.com/OwnLocal/rv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockHandler struct {
	called int
	values []interface{}
	errs   []error
}

func (h *mockHandler) Run(req rv.Request, field *rv.Field) {
	h.values = append(h.values, field.Value)
	h.called++
	for _, err := range h.errs {
		field.Errors = append(field.Errors, err)
	}
}

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

	Describe("TypeHandler", func() {

		Describe("NewTypeHandler", func() {

			for _, typeName := range strings.Split("bool int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64 float32 float64 string", " ") {
				It(fmt.Sprintf("accepts %s type", typeName), func() {
					Expect(rv.NewTypeHandler([]string{typeName})).To(Equal(rv.TypeHandler{Type: typeName}))
				})
			}

			for _, typeName := range strings.Split("uintptr complex64 complex128 array func interface map ptr slice struct unsafe.Pointer", " ") {
				It(fmt.Sprintf("doesn't accept %s type yet", typeName), func() {
					_, err := rv.NewTypeHandler([]string{typeName})
					Expect(err).To(HaveOccurred())
				})
			}

		})

		Describe("Run", func() {
			type tm struct {
				ttype string
				from  interface{}
				to    interface{}
			}

			for _, tc := range []tm{
				tm{"bool", true, true},
				tm{"bool", "1", true},
				tm{"bool", "t", true},
				tm{"bool", "T", true},
				tm{"bool", "TRUE", true},
				tm{"bool", "true", true},
				tm{"bool", "True", true},
				tm{"bool", false, false},
				tm{"bool", "0", false},
				tm{"bool", "f", false},
				tm{"bool", "F", false},
				tm{"bool", "FALSE", false},
				tm{"bool", "false", false},
				tm{"bool", "False", false},

				tm{"bool", 1, true},
				tm{"bool", 0, false},
				tm{"bool", uint(1), true},
				tm{"bool", uint(0), false},

				tm{"int", "42", 42},
				tm{"int", 42, 42},
				tm{"int8", int64(42), int8(42)},
				tm{"int16", int8(42), int16(42)},
				tm{"int32", int(42), int32(42)},
				tm{"int64", int32(42), int64(42)},

				tm{"int", "-42", -42},
				tm{"int", -42, -42},
				tm{"int8", int64(-42), int8(-42)},
				tm{"int16", int8(-42), int16(-42)},
				tm{"int32", int(-42), int32(-42)},
				tm{"int64", int32(-42), int64(-42)},

				tm{"uint", "42", uint(42)},
				tm{"uint", 42, uint(42)},
				tm{"uint8", int64(42), uint8(42)},
				tm{"uint16", int8(42), uint16(42)},
				tm{"uint32", int(42), uint32(42)},
				tm{"uint64", int32(42), uint64(42)},

				tm{"uint8", uint64(42), uint8(42)},
				tm{"uint16", uint8(42), uint16(42)},
				tm{"uint32", uint(42), uint32(42)},
				tm{"uint64", uint32(42), uint64(42)},

				tm{"float32", 42.0, float32(42.0)},

				tm{"string", "yarp", "yarp"},
				tm{"string", false, "false"},
				tm{"string", true, "true"},
				tm{"string", int64(-4), "-4"},
				tm{"string", int32(-3), "-3"},
				tm{"string", int16(-2), "-2"},
				tm{"string", int8(-1), "-1"},
				tm{"string", 0, "0"},
				tm{"string", uint8(1), "1"},
				tm{"string", uint16(2), "2"},
				tm{"string", uint32(3), "3"},
				tm{"string", uint64(4), "4"},
				tm{"string", float32(5.0), "5"},
				tm{"string", float32(6.1), "6.1"},
				tm{"string", float64(7.0), "7"},
				tm{"string", float64(8.1), "8.1"},
			} {
				ttype, from, to := tc.ttype, tc.from, tc.to
				It(fmt.Sprintf("coerces %T(%#v) to %s(%#v)", from, from, ttype, to), func() {
					field.Value = from
					rv.TypeHandler{Type: ttype}.Run(req, field)
					Expect(field.Errors).To(BeEmpty())
					Expect(field.Value).To(Equal(to))
				})
			}

			for _, tc := range []tm{
				tm{"string", nil, nil},
				tm{"bool", nil, nil},
				tm{"int", nil, nil},
				tm{"uint", nil, nil},
				tm{"float32", nil, nil},
			} {
				ttype, from := tc.ttype, tc.from
				It("does not set value if incoming value is nil", func() {
					field.Value = from
					rv.TypeHandler{Type: ttype}.Run(req, field)
					Expect(field.Errors).To(BeEmpty())
					Expect(field.Value).To(BeNil())
				})
			}

			for _, tc := range []tm{
				tm{"bool", "foobly", nil},
				tm{"bool", "42", nil},
				tm{"bool", 42, nil},
				tm{"bool", uint(42), nil},

				tm{"int", "arrr", nil},
				tm{"int8", uint64(42), nil},
				tm{"int16", uint8(42), nil},
				tm{"int32", uint(42), nil},
				tm{"int64", uint32(42), nil},
				tm{"int", uint64(math.MaxUint64), nil},
				tm{"int8", int64(256), nil},
				tm{"int32", uint64(math.MaxUint64), nil},
				tm{"int64", uint64(math.MaxUint64), nil},

				tm{"int8", int64(math.MinInt64), nil},
				tm{"int16", int64(math.MinInt64), nil},
				tm{"int32", int64(math.MinInt64), nil},

				tm{"uint", "-42", nil},
				tm{"uint", -42, nil},
				tm{"uint8", int64(-42), nil},
				tm{"uint16", int8(-42), nil},
				tm{"uint32", int(-42), nil},
				tm{"uint64", int32(-42), nil},

				tm{"float32", "blar", nil},
				tm{"float64", "blar", nil},
			} {
				ttype, from := tc.ttype, tc.from
				It(fmt.Sprintf("cannot coerce %T(%v) to %s", from, from, ttype), func() {
					field.Value = from
					rv.TypeHandler{Type: ttype}.Run(req, field)
					Expect(field.Value).To(Equal(from))
					Expect(field.Errors).ToNot(BeEmpty())
					Expect(field.Errors[0]).To(HaveOccurred())
				})
			}

		})
	})

	Describe("DefaultHandler", func() {
		Describe("NewDefaultHandler", func() {
			It("has a single value if it gets a single argument", func() {
				h, err := rv.NewDefaultHandler([]string{"one"})
				Expect(err).NotTo(HaveOccurred())
				Expect(h.(rv.DefaultHandler).Default).To(Equal("one"))
			})

			It("has a slice of values if it gets a multiple arguments", func() {
				h, err := rv.NewDefaultHandler([]string{"one", "two", "three"})
				Expect(err).NotTo(HaveOccurred())
				Expect(h.(rv.DefaultHandler).Default).To(Equal([]string{"one", "two", "three"}))
			})
		})

		Describe("Run", func() {

			It("does nothing if there is already a value set", func() {
				field.Value = "already set"
				rv.DefaultHandler{Default: "not set"}.Run(req, field)
				Expect(field.Value).To(Equal("already set"))
			})

			It("sets the value to the default there is no value set", func() {
				rv.DefaultHandler{Default: "not set"}.Run(req, field)
				Expect(field.Value).To(Equal("not set"))
			})

		})
	})

	Describe("RangeHandler", func() {
		Describe("Run", func() {

			It("returns no error if value is in range", func() {
				field.Value = 5
				rv.RangeHandler{Start: "1", End: "10"}.Run(req, field)
				Expect(field.Errors).To(BeEmpty())
			})

			It("works with uint values", func() {
				field.Value = uint32(5)
				rv.RangeHandler{Start: "1", End: "10"}.Run(req, field)
				Expect(field.Errors).To(BeEmpty())
			})

			It("works with floating point values", func() {
				field.Value = 5.5
				rv.RangeHandler{Start: "1", End: "10"}.Run(req, field)
				Expect(field.Errors).To(BeEmpty())
			})

			It("works with string values", func() {
				field.Value = "abc"
				rv.RangeHandler{Start: "aaa", End: "ddd"}.Run(req, field)
				Expect(field.Errors).To(BeEmpty())
			})

			It("returns an error for out-of-range string values", func() {
				field.Value = "zzz"
				rv.RangeHandler{Start: "aaa", End: "ddd"}.Run(req, field)
				Expect(field.Errors).ToNot(BeEmpty())
				Expect(field.Errors[0]).To(HaveOccurred())
			})

			It("returns an error if there is no value", func() {
				rv.RangeHandler{Start: "1", End: "10"}.Run(req, field)
				Expect(field.Errors).ToNot(BeEmpty())
				Expect(field.Errors[0]).To(HaveOccurred())
			})

			It("returns an error if the value is out of range", func() {
				field.Value = -1
				rv.RangeHandler{Start: "1", End: "10"}.Run(req, field)
				Expect(field.Errors).ToNot(BeEmpty())
				Expect(field.Errors[0]).To(HaveOccurred())
			})

			It("returns an error if the range is not valid", func() {
				field.Value = 5
				rv.RangeHandler{Start: "one", End: "10"}.Run(req, field)
				Expect(field.Errors).ToNot(BeEmpty())
				Expect(field.Errors[0]).To(HaveOccurred())
			})

		})
	})

	Describe("OptionsHandler", func() {
		y := struct{}{}
		Describe("Run", func() {

			It("returns no error if value is in options", func() {
				field.Value = "two"
				rv.OptionsHandler{Options: map[string]struct{}{"one": y, "two": y, "three": y}}.Run(req, field)
				Expect(field.Errors).To(BeEmpty())
			})

			It("returns an error if value is not in options", func() {
				field.Value = "five"
				rv.OptionsHandler{Options: map[string]struct{}{"one": y, "two": y, "three": y}}.Run(req, field)
				Expect(field.Errors).ToNot(BeEmpty())
				Expect(field.Errors[0]).To(HaveOccurred())
			})

			It("works on ints", func() {
				field.Value = 5
				rv.OptionsHandler{Options: map[string]struct{}{"1": y, "2": y, "3": y}}.Run(req, field)
				Expect(field.Errors).ToNot(BeEmpty())
				Expect(field.Errors[0]).To(HaveOccurred())
			})

			It("works on floats", func() {
				field.Value = 2.2
				rv.OptionsHandler{Options: map[string]struct{}{"1.1": y, "2.2": y, "3.3": y}}.Run(req, field)
				Expect(field.Errors).To(BeEmpty())
			})

		})
	})

	Describe("ListHandler", func() {
		Describe("Run", func() {
			var handler rv.ListHandler

			Context("When the supplied field value is a comma-separated string", func() {
				BeforeEach(func() {
					handler = rv.ListHandler{SubHandlers: rv.FieldHandlers{
						&mockHandler{errs: []error{fmt.Errorf("ow")}},
						&mockHandler{},
					}}
					field.Value = "one,two,three"
					handler.Run(req, field)
				})

				It("splits string arguments and calls each sub-handler with each one", func() {
					Expect(handler.SubHandlers[0].(*mockHandler).called).To(Equal(3))
					Expect(handler.SubHandlers[0].(*mockHandler).values).To(Equal([]interface{}{"one", "two", "three"}))
					Expect(handler.SubHandlers[1].(*mockHandler).called).To(Equal(3))
					Expect(handler.SubHandlers[1].(*mockHandler).values).To(Equal([]interface{}{"one", "two", "three"}))
				})

				It("sets the field value to a list of the types specified", func() {
					Expect(field.Value).To(Equal([]string{"one", "two", "three"}))
				})

				It("gathers errors returned by handlers into the top-level field", func() {
					ow := fmt.Errorf("ow")
					Expect(field.Errors).To(Equal([]error{ow, ow, ow}))
				})

			})

			Context("When the supplied field is already a slice of strings", func() {
				BeforeEach(func() {
					handler = rv.ListHandler{SubHandlers: rv.FieldHandlers{
						&mockHandler{errs: []error{fmt.Errorf("ow")}},
					}}
					field.Value = []string{"one", "two", "three"}
					handler.Run(req, field)
				})

				It("still works as expected", func() {
					Expect(handler.SubHandlers[0].(*mockHandler).called).To(Equal(3))
					Expect(handler.SubHandlers[0].(*mockHandler).values).To(Equal([]interface{}{"one", "two", "three"}))
					Expect(field.Value).To(Equal([]string{"one", "two", "three"}))
					ow := fmt.Errorf("ow")
					Expect(field.Errors).To(Equal([]error{ow, ow, ow}))
				})

			})

			Context("When the supplied field is a single int", func() {
				BeforeEach(func() {
					handler = rv.ListHandler{SubHandlers: rv.FieldHandlers{
						&mockHandler{errs: []error{fmt.Errorf("ow")}},
					}}
					field.Value = 42
					handler.Run(req, field)
				})

				It("calls the SubHandlers for that item as if it were in a slice", func() {
					Expect(handler.SubHandlers[0].(*mockHandler).called).To(Equal(1))
					Expect(handler.SubHandlers[0].(*mockHandler).values).To(Equal([]interface{}{42}))
					Expect(field.Value).To(Equal([]int{42}))
					Expect(field.Errors).To(Equal([]error{fmt.Errorf("ow")}))
				})

			})

			Context("When there is a type sub-handler", func() {
				BeforeEach(func() {
					handler = rv.ListHandler{SubHandlers: rv.FieldHandlers{
						rv.TypeHandler{Type: "int"},
					}}
				})

				It("will transform a comma-separated string with ints into a slice of ints", func() {
					field.Value = "1,2,3"
					handler.Run(req, field)
					Expect(field.Value).To(Equal([]int{1, 2, 3}))
				})

				It("will transform a slice of strings into a slice of ints", func() {
					field.Value = []string{"1", "2", "3"}
					handler.Run(req, field)
					Expect(field.Value).To(Equal([]int{1, 2, 3}))
				})

				It("will transform a slice of ints into a slice of strings", func() {
					field.Value = []int{1, 2, 3}
					handler = rv.ListHandler{SubHandlers: rv.FieldHandlers{
						rv.TypeHandler{Type: "string"},
					}}
					handler.Run(req, field)
					Expect(field.Value).To(Equal([]string{"1", "2", "3"}))
				})

			})
		})
	})

	Describe("RequiredHandler", func() {
		Describe("NewRequiredHandler", func() {
			It("accepts string forms of bools as in strconv.ParseBool", func() {
				Expect(rv.NewRequiredHandler([]string{"1"})).To(Equal(rv.RequiredHandler{Required: true}))
				Expect(rv.NewRequiredHandler([]string{"t"})).To(Equal(rv.RequiredHandler{Required: true}))
				Expect(rv.NewRequiredHandler([]string{"T"})).To(Equal(rv.RequiredHandler{Required: true}))
				Expect(rv.NewRequiredHandler([]string{"TRUE"})).To(Equal(rv.RequiredHandler{Required: true}))
				Expect(rv.NewRequiredHandler([]string{"true"})).To(Equal(rv.RequiredHandler{Required: true}))
				Expect(rv.NewRequiredHandler([]string{"True"})).To(Equal(rv.RequiredHandler{Required: true}))

				Expect(rv.NewRequiredHandler([]string{"0"})).To(Equal(rv.RequiredHandler{Required: false}))
				Expect(rv.NewRequiredHandler([]string{"f"})).To(Equal(rv.RequiredHandler{Required: false}))
				Expect(rv.NewRequiredHandler([]string{"F"})).To(Equal(rv.RequiredHandler{Required: false}))
				Expect(rv.NewRequiredHandler([]string{"FALSE"})).To(Equal(rv.RequiredHandler{Required: false}))
				Expect(rv.NewRequiredHandler([]string{"false"})).To(Equal(rv.RequiredHandler{Required: false}))
				Expect(rv.NewRequiredHandler([]string{"False"})).To(Equal(rv.RequiredHandler{Required: false}))
			})

			It("rejects strings that don't represent bools", func() {
				_, err := rv.NewRequiredHandler([]string{"dksjfsdds"})
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("Run", func() {
			It("adds an error if a required field is not present", func() {
				rv.RequiredHandler{Required: true}.Run(req, field)
				Expect(field.Errors).NotTo(BeEmpty())
				Expect(field.Errors[0]).To(HaveOccurred())
			})
		})
	})

})
