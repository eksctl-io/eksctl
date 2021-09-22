package kubernetes_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	. "github.com/weaveworks/eksctl/pkg/kubernetes"
)

var _ = Describe("Kubernetes client toolkit", func() {
	Describe("can load objects", func() {

		Context("can load and flatten deeply nested lists", func() {
			It("loads all items into flattened list without errors", func() {
				jb, err := ioutil.ReadFile("testdata/misc-sample-nested-list-1.json")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(jb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(6))
			})
		})

		Context("can load and flatten deeply nested lists", func() {
			It("flatten all items into an empty list without errors", func() {
				jb, err := ioutil.ReadFile("testdata/misc-sample-empty-list-1.json")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(jb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(0))
			})
		})

		Context("can combine empty nested lists from a multidoc", func() {
			It("can load without errors", func() {
				yb, err := ioutil.ReadFile("testdata/misc-sample-multidoc-empty-lists-1.yaml")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(yb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(0))
			})
		})

		Context("can combine two empty lists from a multidoc", func() {

			It("can load without errors", func() {
				yb, err := ioutil.ReadFile("testdata/misc-sample-multidoc-empty-lists-2.yaml")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(yb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(0))
			})
		})

		Context("can combine empty and non-empty lists from a multidoc", func() {
			It("can load without errors", func() {
				yb, err := ioutil.ReadFile("testdata/misc-sample-multidoc-nested-lists-1.yaml")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(yb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(4))
			})
		})

		Context("can handle comment nodes", func() {
			It("should be able to parse lists with comment nodes", func() {
				bytes, err := ioutil.ReadFile("testdata/list-with-comment-nodes.yaml")
				Expect(err).ToNot(HaveOccurred())
				list, err := NewList(bytes)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(4))
			})
		})
	})
})

func TestConcatManifests(t *testing.T) {
	a := "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: a\n"
	b := "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: b\n"

	assert.Equal(t, []byte(a), ConcatManifests([][]byte{
		[]byte(a),
	}...))

	assert.Equal(t, []byte(a+"---\n"+b), ConcatManifests([][]byte{
		[]byte(a),
		[]byte(b),
	}...))
}
