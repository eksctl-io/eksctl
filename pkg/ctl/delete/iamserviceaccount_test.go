package delete

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("get", func() {
	Describe("iamserviceaccount", func() {
		It("with cluster flag", func() {
			count := 0
			cmd := newMockEmptyCmd("iamserviceaccount", "--cluster", "dummyName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMServiceAccountWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
					Expect(cmd.NameArg).To(Equal(""))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster argument", func() {
			count := 0
			cmd := newMockEmptyCmd("iamserviceaccount", "dummyName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMServiceAccountWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster and name flag", func() {
			count := 0
			cmd := newMockEmptyCmd("iamserviceaccount", "--name", "serviceAccount")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMServiceAccountWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
					Expect(serviceAccount.Name).To(Equal("serviceAccount"))
					Expect(serviceAccount.Namespace).To(Equal("default"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster and namespace flag", func() {
			count := 0
			cmd := newMockEmptyCmd("iamserviceaccount", "--namespace", "dummyNS")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMServiceAccountWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
					Expect(serviceAccount.Namespace).To(Equal("dummyNS"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("missing required flag --cluster", func() {
			cmd := newMockDefaultCmd("iamserviceaccount")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMServiceAccountWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("invalid flag --dummy", func() {
			cmd := newMockDefaultCmd("iamserviceaccount", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
