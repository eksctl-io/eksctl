package delete

import (
	"fmt"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("delete iamserviceaccount", func() {
	DescribeTable("delete service account successfully",
		func(args ...string) {
			commandArgs := append([]string{"iamserviceaccount"}, args...)
			cmd := newMockEmptyCmd(commandArgs...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteIAMServiceAccountCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, onlyMissing bool) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(cmd.ClusterConfig.IAM.ServiceAccounts[0].Name).To(Equal("serviceAccountName"))
					Expect(onlyMissing).To(Equal(strings.Contains(strings.Join(commandArgs, " "), "only-missing")))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		},
		Entry("with all required flags", "--cluster", "clusterName", "--name", "serviceAccountName"),
		Entry("with namespace flag", "--cluster", "clusterName", "--name", "serviceAccountName", "--namespace", "dev"),
		Entry("with only-missing flag", "--cluster", "clusterName", "--name", "serviceAccountName", "--only-missing"),
		Entry("with approve flag", "--cluster", "clusterName", "--name", "serviceAccountName", "--approve"),
	)

	DescribeTable("invalid flags or arguments",
		func(c invalidParamsCase) {
			commandArgs := append([]string{"iamserviceaccount"}, c.args...)
			cmd := newDefaultCmd(commandArgs...)
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(c.error.Error()))
		},
		Entry("without cluster name", invalidParamsCase{
			args:  []string{"--name", "serviceAccountName"},
			error: fmt.Errorf("--cluster must be set"),
		}),
		Entry("with iamserviceaccount name as argument and flag", invalidParamsCase{
			args:  []string{"--cluster", "clusterName", "--name", "serviceAccountName", "serviceAccountName"},
			error: fmt.Errorf("--name=serviceAccountName and argument serviceAccountName cannot be used at the same time"),
		}),
		Entry("with invalid flags", invalidParamsCase{
			args:  []string{"iamserviceaccount", "--invalid", "dummy"},
			error: fmt.Errorf("unknown flag: --invalid"),
		}),
	)
})
