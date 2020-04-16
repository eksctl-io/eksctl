package delete

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type invalidParamsCase struct {
	args  []string
	error error
}

var _ = Describe("delete", func() {
	DescribeTable("drain node group successfully",
		func(args ...string) {
			cmd := newMockEmptyCmd(args...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *v1alpha5.NodeGroup, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(ng.Name).To(Equal("ng"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		},
		Entry("with valid details", "nodegroup", "--cluster", "clusterName", "--name", "ng"),
		Entry("with deprecated flag --only", "nodegroup", "--cluster", "clusterName", "--name", "ng", "--only", "ng"),
	)

	DescribeTable("invalid flags or arguments",
		func(c invalidParamsCase) {
			cmd := newDefaultCmd(c.args...)
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(c.error.Error()))
		},
		Entry("missing required flag --cluster", invalidParamsCase{
			args:  []string{"nodegroup"},
			error: fmt.Errorf("--cluster must be set"),
		}),
		Entry("setting --name and argument at the same time", invalidParamsCase{
			args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--name", "ng"},
			error: fmt.Errorf("--name=ng and argument ng cannot be used at the same time"),
		}),
	)
})
