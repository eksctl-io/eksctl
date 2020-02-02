package delete

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("delete", func() {
	Describe("iamidentitymapping", func() {
		It("with cluster name argument", func() {
			count := 0
			cmd := newMockEmptyCmd("iamidentitymapping", "--all")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMIdentityMappingWithRunFunc(cmd, func(cmd *cmdutils.Cmd, arn string, all bool) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster name flag (--cluster)", func() {
			count := 0
			cmd := newMockEmptyCmd("iamidentitymapping", "--cluster", "dummyName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMIdentityMappingWithRunFunc(cmd, func(cmd *cmdutils.Cmd, arn string, all bool) error {
					Expect(cmd.NameArg).To(Equal(""))
					count++
					return nil
				})
			})
			out, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
			Expect(out).To(Not(ContainSubstring("Flag --name has been deprecated, use --cluster")))
		})

		It("no cluster name", func() {
			cmd := newMockDefaultCmd("iamidentitymapping")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("invalid flag --dummy", func() {
			cmd := newMockDefaultCmd("iamidentitymapping", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
