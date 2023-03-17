package filter

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("filter", func() {

	Context("Match", func() {
		var filter Filter
		allItems := []string{
			"a1",
			"a2",
			"b1",
			"b2",
			"banana",
			"apple",
			"pineapple",
			"strawberry",
			"raspberry",
		}

		BeforeEach(func() {
			filter = NewFilter()
		})

		It("should include everything when there are no rules", func() {
			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(9))
			Expect(excluded).To(HaveLen(0))
			Expect(included.HasAll(allItems...)).To(BeTrue())
		})

		It("should exclude everything when ExcludeAll is enabled", func() {
			filter.ExcludeAll = true
			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(0))
			Expect(excluded).To(HaveLen(9))
			Expect(excluded.HasAll(allItems...)).To(BeTrue())
		})

		It("should match include filter", func() {
			filter.AppendIncludeNames("banana")
			err := filter.doAppendIncludeGlobs(allItems, "fruits", "*apple", "*berry")
			Expect(err).NotTo(HaveOccurred())

			Expect(filter.Match("a1")).To(BeFalse())
			Expect(filter.Match("a2")).To(BeFalse())
			Expect(filter.Match("b1")).To(BeFalse())
			Expect(filter.Match("b2")).To(BeFalse())
			Expect(filter.Match("banana")).To(BeTrue())
			Expect(filter.Match("apple")).To(BeTrue())
			Expect(filter.Match("pineapple")).To(BeTrue())
			Expect(filter.Match("strawberry")).To(BeTrue())
			Expect(filter.Match("raspberry")).To(BeTrue())

			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(5))
			Expect(included.HasAll("banana", "apple", "pineapple", "strawberry", "raspberry")).To(BeTrue())
			Expect(excluded).To(HaveLen(4))
			Expect(excluded.HasAll("a1", "a2", "b1", "b2")).To(BeTrue())
		})

		It("should match exclude filter with names and globs", func() {
			filter.AppendExcludeNames("banana")
			err := filter.AppendExcludeGlobs("*apple", "*berry")
			Expect(err).NotTo(HaveOccurred())

			Expect(filter.Match("a1")).To(BeTrue())
			Expect(filter.Match("a2")).To(BeTrue())
			Expect(filter.Match("b1")).To(BeTrue())
			Expect(filter.Match("b2")).To(BeTrue())
			Expect(filter.Match("banana")).To(BeFalse())
			Expect(filter.Match("apple")).To(BeFalse())
			Expect(filter.Match("pineapple")).To(BeFalse())
			Expect(filter.Match("strawberry")).To(BeFalse())
			Expect(filter.Match("raspberry")).To(BeFalse())

			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(4))
			Expect(included.HasAll("a1", "a2", "b1", "b2")).To(BeTrue())
			Expect(excluded).To(HaveLen(5))
			Expect(excluded.HasAll("banana", "apple", "pineapple", "strawberry", "raspberry")).To(BeTrue())
		})

		It("should include an item when it exists as an inclusion name overwrite", func() {
			filter.AppendIncludeNames("raspberry")
			err := filter.doAppendIncludeGlobs(allItems, "fruit", "a?", "b?")
			Expect(err).NotTo(HaveOccurred())

			err = filter.AppendExcludeGlobs("*apple", "*berry")
			Expect(err).NotTo(HaveOccurred())

			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(5))
			Expect(included.HasAll("raspberry", "a1", "a2", "b1", "b2")).To(BeTrue())
			Expect(excluded).To(HaveLen(4))
			Expect(excluded.HasAll("banana", "apple", "pineapple", "strawberry")).To(BeTrue())
		})

		It("should not include an item when it exists as an exclusion name overwrite", func() {
			err := filter.doAppendIncludeGlobs(allItems, "fruit", "a?", "b?")
			Expect(err).NotTo(HaveOccurred())

			filter.AppendExcludeNames("a2")

			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(3))
			Expect(included.HasAll("a1", "b1", "b2")).To(BeTrue())
			Expect(excluded).To(HaveLen(6))
			Expect(excluded.HasAll("a2", "banana", "apple", "pineapple", "strawberry", "raspberry")).To(BeTrue())
		})

		It("when a name is in both inclusion and exclusion overwrites the exclusion takes precedence", func() {
			filter.AppendIncludeNames("raspberry")
			err := filter.doAppendIncludeGlobs(allItems, "fruit", "a?", "b?")
			Expect(err).NotTo(HaveOccurred())

			filter.AppendExcludeNames("raspberry")
			err = filter.AppendExcludeGlobs("*apple", "*berry")
			Expect(err).NotTo(HaveOccurred())

			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(4))
			Expect(included.HasAll("a1", "a2", "b1", "b2")).To(BeTrue())
			Expect(excluded).To(HaveLen(5))
			Expect(excluded.HasAll("raspberry", "banana", "apple", "pineapple", "strawberry")).To(BeTrue())
		})

		It("when an item matches inclusion and exclusion globs exclusion takes precedence", func() {
			err := filter.doAppendIncludeGlobs(allItems, "fruit", "*berry", "*apple")
			Expect(err).NotTo(HaveOccurred())

			err = filter.AppendExcludeGlobs("?aspberry", "a?")
			Expect(err).NotTo(HaveOccurred())

			included, excluded := filter.doMatchAll(allItems)
			Expect(included).To(HaveLen(3))
			Expect(included.HasAll("strawberry", "apple", "pineapple")).To(BeTrue())
			Expect(excluded).To(HaveLen(6))
			Expect(excluded.HasAll("raspberry", "banana", "a1", "a2", "b1", "b2")).To(BeTrue())

		})
	})
})
