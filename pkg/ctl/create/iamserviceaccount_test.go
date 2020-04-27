package create

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("create iamserviceaccount", func() {
	DescribeTable("create service account successfully",
		func(args ...string) {
			commandArgs := append([]string{"iamserviceaccount"}, args...)
			cmd := newMockEmptyCmd(commandArgs...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				createIAMServiceAccountCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, overrideExistingServiceAccounts bool) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(cmd.ClusterConfig.IAM.ServiceAccounts[0].Name).To(Equal("serviceAccountName"))
					Expect(cmd.ClusterConfig.IAM.ServiceAccounts[0].AttachPolicyARNs).To(ContainElement("dummyPolicyArn"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		},
		Entry("with all required flags", "--cluster", "clusterName", "--name", "serviceAccountName", "--attach-policy-arn", "dummyPolicyArn"),
		Entry("with override flags", "--cluster", "clusterName", "--name", "serviceAccountName", "--attach-policy-arn", "dummyPolicyArn", "--override-existing-serviceaccounts"),
	)

	DescribeTable("invalid flags or arguments",
		func(c invalidParamsCase) {
			cmd := newDefaultCmd(c.args...)
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(c.error.Error()))
		},
		Entry("without cluster name", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--name", "serviceAccountName"},
			error: fmt.Errorf("--cluster must be set"),
		}),
		Entry("with iamserviceaccount name as argument and flag", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--cluster", "clusterName", "--name", "serviceAccountName", "serviceAccountName"},
			error: fmt.Errorf("--name=serviceAccountName and argument serviceAccountName cannot be used at the same time"),
		}),
		Entry("without required flag --attach-policy-arn", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--cluster", "clusterName", "serviceAccountName"},
			error: fmt.Errorf("--attach-policy-arn must be set"),
		}),
		Entry("with invalid flags", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--invalid", "dummy"},
			error: fmt.Errorf("unknown flag: --invalid"),
		}),
	)
})
