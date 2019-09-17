package dir_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/dir"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("dir.IsEmpty", func() {
	It("should return true when passed an empty directory", func() {
		d, err := ioutil.TempDir("", "test_dir_isempty")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(d) // clean up.

		isEmpty, err := dir.IsEmpty(d)
		Expect(err).NotTo(HaveOccurred())
		Expect(isEmpty).To(BeTrue())
	})

	It("should return false when passed a directory containing another directory", func() {
		d, err := ioutil.TempDir("", "test_dir_isempty")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(d) // clean up.
		err = os.Mkdir(filepath.Join(d, "subdir"), 0755)
		Expect(err).NotTo(HaveOccurred())

		isEmpty, err := dir.IsEmpty(d)
		Expect(err).NotTo(HaveOccurred())
		Expect(isEmpty).To(BeFalse())
	})

	It("should return false when passed a directory containing a file", func() {
		d, err := ioutil.TempDir("", "test_dir_isempty")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(d) // clean up.
		f, err := os.OpenFile(filepath.Join(d, "file.tmp"), os.O_CREATE, 0755)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(f.Name()) // clean up.

		isEmpty, err := dir.IsEmpty(d)
		Expect(err).NotTo(HaveOccurred())
		Expect(isEmpty).To(BeFalse())
	})
})
