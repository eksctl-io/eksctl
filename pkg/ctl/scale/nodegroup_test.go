package scale

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("scale", func() {
	Describe("scale nodegroup", func() {
		DescribeTable("scales  a nodegroup successfully",
			func(args ...string) {
				cmd := newMockEmptyCmd(args...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					scaleNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *v1alpha5.NodeGroupBase) error {
						if len(ng.Name) != 0 {
							Expect(ng.Name).To(Or(Equal("nodeGroup"), Equal("")))
						} else {
							Expect(cmd.NameArg).To(Or(Equal("nodeGroup"), Equal("")))
						}
						if len(cmd.ClusterConfig.Metadata.Name) != 0 {
							Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
							Expect(*ng.DesiredCapacity).To(Equal(2))
						} else {
							Expect(cmd.ClusterConfigFile).To(Equal("dummyConfigFile.yaml"))
						}
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).To(Not(HaveOccurred()))
				Expect(count).To(Equal(1))
			},
			Entry("with all the valid flags", "nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--nodes", "2", "--nodes-max", "3", "--nodes-min", "1", "--wait", "--timeout", "25m"),
			Entry("with config file and name without wait flag", "nodegroup", "nodeGroup", "-f", "dummyConfigFile.yaml"),
			Entry("with name flag and config file without wait flag", "nodegroup", "--name", "nodeGroup", "-f", "dummyConfigFile.yaml"),
			Entry("without --nodes-min and wait flags", "nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--nodes", "2", "--nodes-max", "3"),
			Entry("without --nodes-max and wait flags", "nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--nodes", "2", "--nodes-min", "1"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				cmd := newDefaultCmd(c.args...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(c.error.Error()))
			},
			Entry("missing required flag --cluster", invalidParamsCase{
				args:  []string{"nodegroup", "ng"},
				error: fmt.Errorf("Error: --cluster must be set"),
			}),
			Entry("setting --name and argument at the same time", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--name", "ng"},
				error: fmt.Errorf("Error: --name=ng and argument ng cannot be used at the same time"),
			}),
			Entry("missing required nodes flag", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy"},
				error: fmt.Errorf("Error: at least one of minimum, maximum and desired nodes must be set"),
			}),
			Entry("invalid flag", invalidParamsCase{
				args:  []string{"nodegroup", "--invalid", "dummy"},
				error: fmt.Errorf("Error: unknown flag: --invalid"),
			}),
			Entry("desired node fewer than min nodes", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--nodes", "2", "--nodes-min", "3"},
				error: fmt.Errorf("Error: minimum number of nodes must be fewer than or equal to number of nodes"),
			}),

			Entry("desired node greater than max nodes", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--nodes", "2", "--nodes-max", "1"},
				error: fmt.Errorf("Error: maximum number of nodes must be greater than or equal to number of nodes"),
			}),
			Entry("desired node fewer than min nodes", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--nodes", "2", "--nodes-min", "3"},
				error: fmt.Errorf("Error: minimum number of nodes must be fewer than or equal to number of nodes"),
			}),
			Entry("with config file and nodes flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--nodes", "2"},
				error: fmt.Errorf("Error: cannot use --nodes when --config-file/-f is set"),
			}),
			Entry("with config file and nodes-max flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--nodes-max", "2"},
				error: fmt.Errorf("Error: cannot use --nodes-max when --config-file/-f is set"),
			}),
			Entry("with config file and nodes-min flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--nodes-min", "2"},
				error: fmt.Errorf("Error: cannot use --nodes-min when --config-file/-f is set"),
			}),
			Entry("with config file and cluster flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--cluster", "dummyCluster"},
				error: fmt.Errorf("Error: cannot use --cluster when --config-file/-f is set"),
			}),
		)
	})
})
