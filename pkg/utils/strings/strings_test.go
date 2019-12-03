package strings_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("strings", func() {
	Describe("Pointer", func() {
		It("returns a pointer to the provided string", func() {
			p := strings.Pointer("test")
			Expect(*p).To(Equal("test"))
		})
	})

	Describe("NilIfEmpty", func() {
		It("returns a pointer to the provided string if non-empty", func() {
			p := strings.NilIfEmpty("test")
			Expect(*p).To(Equal("test"))
		})

		It("returns nil if the provided string is empty", func() {
			p := strings.NilIfEmpty("")
			Expect(p).To(BeNil())
		})
	})

	Describe("EmptyIfNil", func() {
		It("returns the value pointed to by the provided pointer if non-nil", func() {
			in := "test"
			out := strings.EmptyIfNil(&in)
			Expect(out).To(Equal(in))
		})

		It("returns an empty string if the provided pointer is nil", func() {
			out := strings.EmptyIfNil(nil)
			Expect(out).To(Equal(""))
		})
	})

	Describe("ToPointersMap", func() {
		It("converts the provided map[string]string to a map[string]*string", func() {
			valuesMap := map[string]string{"a": "A", "b": "B"}
			pointersMap := strings.ToPointersMap(valuesMap)
			Expect(pointersMap).To(HaveLen(len(valuesMap))) // that is...
			Expect(pointersMap).To(HaveLen(2))
			Expect(pointersMap).To(HaveKey("a"))
			Expect(*pointersMap["a"]).To(Equal("A"))
			Expect(pointersMap).To(HaveKey("b"))
			Expect(*pointersMap["b"]).To(Equal("B"))
		})
	})

	Describe("NilPointersMapIfEmpty", func() {
		It("returns the provided map if non-empty", func() {
			a := "A"
			in := map[string]*string{"a": &a}
			out := strings.NilPointersMapIfEmpty(in)
			Expect(out).To(Equal(in))
		})

		It("returns nil if the provided map is empty", func() {
			in := map[string]*string{}
			out := strings.NilPointersMapIfEmpty(in)
			Expect(out).To(BeNil())
		})
	})

	Describe("ToValuesMap", func() {
		It("converts the provided map[string]*string to a map[string]string", func() {
			a := "A"
			b := "B"
			pointersMap := map[string]*string{"a": &a, "b": &b}
			valuesMap := strings.ToValuesMap(pointersMap)
			Expect(valuesMap).To(HaveLen(len(pointersMap))) // that is...
			Expect(valuesMap).To(HaveLen(2))
			Expect(valuesMap).To(HaveKeyWithValue("a", "A"))
			Expect(valuesMap).To(HaveKeyWithValue("b", "B"))
		})
	})

	Describe("ToPointersArray", func() {
		It("converts the provided []string to a []*string", func() {
			valuesArray := []string{"a", "b"}
			pointersArray := strings.ToPointersArray(valuesArray)
			Expect(pointersArray).To(HaveLen(len(valuesArray))) // that is...
			Expect(pointersArray).To(HaveLen(2))
			Expect(*pointersArray[0]).To(Equal("a"))
			Expect(*pointersArray[1]).To(Equal("b"))
		})
	})

	Describe("NilPointersArrayIfEmpty", func() {
		It("returns the provided array if non-empty", func() {
			a := "a"
			in := []*string{&a}
			out := strings.NilPointersArrayIfEmpty(in)
			Expect(out).To(Equal(in))
		})

		It("returns nil if the provided array is empty", func() {
			in := []*string{}
			out := strings.NilPointersArrayIfEmpty(in)
			Expect(out).To(BeNil())
		})
	})

	Describe("ToValuesArray", func() {
		It("converts the provided []*string to a []string", func() {
			a := "a"
			b := "b"
			pointersArray := []*string{&a, &b}
			valuesArray := strings.ToValuesArray(pointersArray)
			Expect(valuesArray).To(HaveLen(len(pointersArray))) // that is...
			Expect(valuesArray).To(HaveLen(2))
			Expect(valuesArray[0]).To(Equal("a"))
			Expect(valuesArray[1]).To(Equal("b"))
		})
	})
})
