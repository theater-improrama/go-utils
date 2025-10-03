package slice_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/theater-improrama/go-utils/slice"
)

var _ = Describe("Partition", func() {
	It("should partition an array into equally sized groups, except for the last group", func() {
		s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		p := slice.PartitionSplit(s, 5)

		Expect(len(p)).To(Equal(3))
		Expect(p[0]).To(Equal([]int{1, 2, 3, 4, 5}))
		Expect(p[1]).To(Equal([]int{6, 7, 8, 9, 10}))
		Expect(p[2]).To(Equal([]int{11, 12}))
	})

	It("should partition an array into equally sized groups", func() {
		s := []int{1, 2, 3, 4}

		p := slice.PartitionSplit(s, 2)
		Expect(len(p)).To(Equal(2))
		Expect(p[0]).To(Equal([]int{1, 2}))
		Expect(p[1]).To(Equal([]int{3, 4}))
	})

	It("should return an empty partition", func() {
		var s []int

		p := slice.PartitionSplit(s, 1)

		Expect(len(p)).To(Equal(0))
	})
})
