// +build integration

package matchers

import (
	"fmt"
	"os"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/weaveworks/eksctl/integration/utilities/git"
)

// AssertQuickStartComponentsPresentInGit asserts that the expected quickstart
// components are present in Git, under the provided branch.
func AssertQuickStartComponentsPresentInGit(branch, privateSSHKeyPath string) {
	dir, err := git.GetBranch(branch, privateSSHKeyPath)
	Expect(err).ShouldNot(HaveOccurred())
	defer os.RemoveAll(dir)
	FS := afero.Afero{Fs: afero.NewOsFs()}
	allFiles := make([]string, 0)
	err = FS.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		allFiles = append(allFiles, path)
		return nil
	})
	Expect(err).ToNot(HaveOccurred())
	fmt.Fprintf(ginkgo.GinkgoWriter, "\n all files:\n%v", allFiles)
}
