package enable

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("enable", func() {
	Describe("profile", func() {
		It("with required flag", func() {
			count := 0
			cmd := newMockEmptyCmd("profile", "--git-url", "git://dummy-repo", "--git-email", "test@test.com")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				enableProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd, opts *ProfileOptions) error {
					Expect(opts.gitOptions.Email).To(Equal("test@test.com"))
					Expect(opts.gitOptions.URL).To(Equal("git://dummy-repo"))
					Expect(opts.profileOptions.Name).To(Equal(""))
					Expect(opts.profileOptions.Revision).To(Equal("master"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("missing required flag --git-email", func() {
			cmd := newMockDefaultCmd("profile", "--git-url", "git://dummy-repo")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"git-email\" not set"))
		})

		It("missing required flag --git-url", func() {
			cmd := newMockDefaultCmd("profile", "--git-email", "test@test.com")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"git-url\" not set"))
		})

		It("missing all required flags", func() {
			cmd := newMockDefaultCmd("profile")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"git-email\", \"git-url\" not set"))
		})

		It("invalid flag --dummy", func() {
			cmd := newMockDefaultCmd("profile", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
