package set_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/theater-improrama/go-utils/set"
)

var _ = Describe("Set", func() {
	It("should contain a unique element once", func() {
		s := set.New[string]()
		s.Set("test")
		s.Set("test")

		Expect(len(s.Values())).To(Equal(1))
		Expect(s.Values()).To(Equal([]string{"test"}))
	})

	It("should delete only the unique element", func() {
		s := set.New[string]()
		s.Set("test1")
		s.Set("test2")
		s.Delete("test2")

		Expect(len(s.Values())).To(Equal(1))
		Expect(s.Values()).To(Equal([]string{"test1"}))
	})

	It("should be possible to check for an elements existence", func() {
		s := set.New[string]()
		s.Set("test1")

		Expect(s.Exists("test1")).To(BeTrue())
		Expect(s.Exists("test2")).To(BeFalse())
	})

	It("should be generated from a slice", func() {
		s := set.From([]string{"test1", "test2"})

		Expect(s.Exists("test1")).To(BeTrue())
		Expect(s.Exists("test2")).To(BeTrue())
	})

	It("should return the correct set intersection", func() {
		s1 := set.From([]int{1, 2, 3})
		s2 := set.From([]int{3})

		Expect(s1.Intersect(s2).Values()).To(Equal([]int{3}))
	})

	It("should return the correct set difference", func() {
		s1 := set.From([]int{1, 2, 3, 4})
		s2 := set.From([]int{2, 4})

		Expect(s1.Difference(s2).Values()).To(ContainElements(1, 3))
		Expect(s2.Difference(s1).Values()).To(Equal([]int{}))
	})
})
