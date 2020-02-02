package generate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("generate", func() {
	Describe("profile", func() {
		It("with required flag", func() {
			count := 0
			cmd := newMockEmptyCmd("profile", "--git-url", "git://dummy-repo")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				generateProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd, o options) error {
					Expect(o.GitOptions.URL).To(Equal("git://dummy-repo"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("without required flag", func() {
			cmd := newMockDefaultCmd("profile")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"git-url\" not set"))
		})

		It("with all flags", func() {
			count := 0
			cmd := newMockEmptyCmd("profile", "--git-url", "git://dummy-repo", "--git-branch", "master", "--profile-path", "/")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				generateProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd, o options) error {
					Expect(o.GitOptions.URL).To(Equal("git://dummy-repo"))
					Expect(o.GitOptions.Branch).To(Equal("master"))
					Expect(o.ProfilePath).To(Equal("/"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with deprecated cluster flag", func() {
			count := 0
			cmd := newMockEmptyCmd("profile", "--git-url", "git://dummy-repo", "--name", "clusterName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				generateProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd, o options) error {
					Expect(o.GitOptions.URL).To(Equal("git://dummy-repo"))
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					count++
					return nil
				})
			})
			out, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
			Expect(out).To(ContainSubstring("Flag --name has been deprecated, use --cluster"))
		})

		It("invalid flag --dummy", func() {
			cmd := newMockDefaultCmd("profile", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
