package enable

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
)

var _ = Describe("enable", func() {
	Describe("repo", func() {
		It("with required flag", func() {
			count := 0
			cmd := newMockEmptyCmd("repo", "--git-url", "git://dummy-repo", "--git-email", "test@test.com")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				enableRepoWithRunFunc(cmd, func(cmd *cmdutils.Cmd, opts *flux.InstallOpts) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("missing required flag --git-email", func() {
			cmd := newMockDefaultCmd("repo", "--git-url", "git://dummy-repo")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"git-email\" not set"))
		})

		It("missing required flag --git-url", func() {
			cmd := newMockDefaultCmd("repo", "--git-email", "test@test.com")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"git-url\" not set"))
		})

		It("missing all required flags", func() {
			cmd := newMockDefaultCmd("repo")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"git-email\", \"git-url\" not set"))
		})

		It("invalid flag --dummy", func() {
			cmd := newMockDefaultCmd("repo", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
