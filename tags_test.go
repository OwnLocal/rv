package rv

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("tags", func() {

	Describe("extractTags", func() {

		It("generates a time entry for a time.Time struct field", func() {
			tagMap, err := extractTags(struct {
				T time.Time `rv:"query.t"`
			}{})
			expected := map[string]map[string][]string{
				"T": map[string][]string{
					"type":   []string{"time"},
					"source": []string{"query.t"},
				},
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(tagMap).To(Equal(expected))
		})

	})
})
