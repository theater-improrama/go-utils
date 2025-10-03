package slice_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/theater-improrama/go-utils/slice"
)

var _ = Describe("Unique", func() {
	It("should return a unique string array", func() {
		arr := []string{"test", "test"}

		arr = slice.Unique(arr)

		Expect(arr).To(Equal([]string{"test"}))
	})
})
