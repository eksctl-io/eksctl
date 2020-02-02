package get

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("get", func() {
	Describe("iamidentitymapping", func() {

		It("with cluster name argument", func() {
			count := 0
			cmd := newMockEmptyCmd( "iamidentitymapping", "dummyName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				getIAMIdentityMappingWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, arn string) error {
					Expect(cmd.NameArg).To(Equal("dummyName"))
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
			cmd := newMockEmptyCmd( "iamidentitymapping", "--cluster", "dummyName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				getIAMIdentityMappingWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, arn string) error {
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

		It("with cluster name flag (--name) (deprecated)", func() {
			count := 0
			cmd := newMockEmptyCmd( "iamidentitymapping", "--name", "dummyName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				getIAMIdentityMappingWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, arn string) error {
					count++
					return nil
				})
			})
			out, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
			Expect(out).To(ContainSubstring("Flag --name has been deprecated, use --cluster"))
		})

		It("cluster name argument and --cluster flag", func() {
			cmd := newMockDefaultCmd( "iamidentitymapping", "dummyName", "--cluster", "dummyName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster=dummyName and argument dummyName cannot be used at the same time"))
		})

		It("no cluster name argument or flag", func() {
			cmd := newMockDefaultCmd( "iamidentitymapping")
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
