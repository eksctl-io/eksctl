package nodegroup_test

import (
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNodegroup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nodegroup Suite")
}

var (
	al2WithoutForceTemplate, al2ForceFalseTemplate, al2FullyUpdatedTemplate string
	brForceFalseTemplate, brForceTrueTemplate, brFulllyUpdatedTemplate      string
)

var _ = BeforeSuite(func() {
	//Does not have ForceUpdateEnabled specified
	al2WithoutForceTemplate = mustReadFile("testdata/al2-no-force-template.json")
	//ForceUpdateEnabled set to false
	al2ForceFalseTemplate = mustReadFile("testdata/al2-force-false-template.json")
	//ForceUpdateEnabled set to false and ReleaseVersion set to 1.20-20201212
	al2FullyUpdatedTemplate = mustReadFile("testdata/al2-updated-template.json")
	//ForceUpdateEnabled set to false
	brForceFalseTemplate = mustReadFile("testdata/br-force-false-template.json")
	//ForceUpdateEnabled set to true
	brForceTrueTemplate = mustReadFile("testdata/br-force-true-template.json")
	//ForceUpdateEnabled set to true and Version set to 1.21
	brFulllyUpdatedTemplate = mustReadFile("testdata/br-updated-template.json")
})

func mustReadFile(path string) string {
	bytes, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	return strings.Trim(string(bytes), "\n")
}
