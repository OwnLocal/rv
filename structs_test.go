package rv_test

import (
	//. "github.com/ownlocal/rv"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
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

var _ = Describe("Args", func() {
	/*
		Describe("ParseTag", func() {
			It("parses tags with just a source field specified", func() {
				Expect(ParseTag(`path.foo`)).To(Equal(Args{Source: PATH, Field: "foo"}))
			})

			It("parses the range argument", func() {
				Expect(ParseTag(`query.foo range=1,5`)).To(Equal(Args{Source: QUERY, Field: "foo", Range: Range{1, 5}}))
			})

			It("parses the options argument", func() {
				Expect(ParseTag(`json.foo options=a,b,c`)).To(Equal(Args{Source: JSON, Field: "foo", Options: []string{"a", "b", "c"}}))
			})

			It("parses the default argument", func() {
				Expect(ParseTag(`form.foo default=a`)).To(Equal(Args{Source: FORM, Field: "foo", Default: "a"}))
			})
		})

	*/
})
