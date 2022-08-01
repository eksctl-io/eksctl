package create

import (
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	. "github.com/onsi/ginkgo/v2"
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
		Entry("with optional flags", "--cluster", "clusterName", "--name", "serviceAccountName", "--attach-policy-arn", "dummyPolicyArn", "--override-existing-serviceaccounts", "--role-name", "custom-role-name"),
	)

	DescribeTable("invalid flags or arguments",
		func(c invalidParamsCase) {
			cmd := newDefaultCmd(c.args...)
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(c.error)))
		},
		Entry("without cluster name", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--name", "serviceAccountName"},
			error: "--cluster must be set",
		}),
		Entry("with iamserviceaccount name as argument and flag", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--cluster", "clusterName", "--name", "serviceAccountName", "serviceAccountName"},
			error: "--name=serviceAccountName and argument serviceAccountName cannot be used at the same time",
		}),
		Entry("without required flags --attach-policy-arn or --attach-policy-role", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--cluster", "clusterName", "serviceAccountName"},
			error: "--attach-policy-arn or --attach-role-arn must be set",
		}),
		Entry("with --attach-role-arn and --role-name", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--cluster", "clusterName", "serviceAccountName", "--role-name", "foo", "--attach-role-arn", "123"},
			error: "cannot provide --role-name or --role-only when --attach-role-arn is configured",
		}),
		Entry("with --attach-policy-role and --role-only", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--cluster", "clusterName", "serviceAccountName", "--role-only", "--attach-role-arn", "123"},
			error: "cannot provide --role-name or --role-only when --attach-role-arn is configured",
		}),
		Entry("with --attach-role-arn and --attach-policy-arns", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--cluster", "clusterName", "serviceAccountName", "--attach-policy-arn", "123", "--attach-role-arn", "123"},
			error: "cannot provide --attach-role-arn and specify polices to attach",
		}),
		Entry("with invalid flags", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--invalid", "dummy"},
			error: "unknown flag: --invalid",
		}),
	)
})
