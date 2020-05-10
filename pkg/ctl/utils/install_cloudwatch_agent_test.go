package utils

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("utils", func() {
	BeforeEach(func() {
		_ = os.Setenv("EKSCTL_EXPERIMENTAL", "true")
	})

	Describe("install-cloudwatch-agent", func() {
		DescribeTable("install cloudwatch agent successfully",
			func(args ...string) {
				allArgs := []string{"install-cloudwatch-agent"}
				allArgs = append(allArgs, args...)
				cmd := newMockEmptyCmd(allArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					installCloudwatchAgentWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).To(Not(HaveOccurred()))
				Expect(count).To(Equal(1))
			},
			Entry("with cluster flag", "--cluster", "clusterName"),
			Entry("with cluster and region flag", "--cluster", "clusterName", "--region", "dummyRegion"),
			Entry("with config file flag", "--config-file", "dummyConfigFile"),
			Entry("with config file flag and approve flag", "--config-file", "dummyConfigFile", "--approve"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				args := []string{"install-cloudwatch-agent"}
				args = append(args, c.args...)
				cmd := newDefaultCmd(args...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(c.error.Error()))
			},
			Entry("missing required flag --cluster", invalidParamsCase{
				args:  []string{"install-cloudwatch-agent"},
				error: fmt.Errorf("--cluster must be set"),
			}),
		)
	})
})
