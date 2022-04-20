//go:build integration
// +build integration

package anywhere

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/pkg/actions/anywhere"

	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo/v2"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("anywhere")
}

func TestAnywhere(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [EKS Anywhere]", func() {
	Context("--help", func() {
		It("shows EKS anywhere in the help text", func() {
			cmd := params.EksctlHelpCmd.WithArgs("")
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("eksctl anywhere")),
				ContainElement(ContainSubstring("EKS anywhere")),
			))
		})
	})

	Context("eksctl anywhere", func() {
		var (
			tmpDir string
		)

		BeforeEach(func() {
			var err error
			tmpDir, err = os.MkdirTemp("", "anywhere-command")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		When("the binary exists in the path", func() {
			var originalPath string

			BeforeEach(func() {
				err := os.WriteFile(filepath.Join(tmpDir, anywhere.BinaryFileName), []byte(`#!/usr/bin/env sh
echo "you called?"
exit 0`), 0777)
				Expect(err).NotTo(HaveOccurred())

				originalPath = os.Getenv("PATH")
				Expect(os.Setenv("PATH", fmt.Sprintf("%s:%s", originalPath, tmpDir))).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Setenv("PATH", originalPath)).To(Succeed())
			})

			It("invokes the binary", func() {
				cmd := params.EksctlAnywhereCmd.WithArgs("")
				Expect(cmd).To(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("you called?"))))
			})
		})

		When("the binary is not on the path", func() {
			It("returns an error", func() {
				cmd := params.EksctlAnywhereCmd.WithArgs("")
				session := cmd.Run()
				Expect(session.ExitCode()).To(Equal(1))
				Expect(string(session.Out.Contents())).To(Equal(fmt.Sprintf("%q plugin was not found on your path\n", anywhere.BinaryFileName)))
			})
		})
	})
})
