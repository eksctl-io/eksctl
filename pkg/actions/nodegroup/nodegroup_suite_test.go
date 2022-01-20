package nodegroup_test

import (
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNodegroup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nodegroup Suite")
}

var (
	al2Template, al2UpdatedTemplate string
)

var _ = BeforeSuite(func() {
	al2Template = mustReadFile("testdata/al2-template.json")
	al2UpdatedTemplate = strings.Trim(mustReadFile("testdata/al2-updated-template.json"), "\n")
})

func mustReadFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	return string(bytes)
}
