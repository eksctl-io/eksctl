package scale

import (
	"fmt"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("scale", func() {
	Describe("nodegroup", func() {
		DescribeTable("create cluster successfully",
			func(args ...string) {
				cmd := newMockEmptyCmd(args...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					scaleNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *v1alpha5.NodeGroup) error {
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
			Entry("with all the valid flags", "nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--nodes", "2", "--nodes-max", "3", "--nodes-min", "1"),
			Entry("with config file and name", "nodegroup", "nodeGroup", "-f", "dummyConfigFile.yaml"),
			Entry("with config file and name flags", "nodegroup", "--name", "nodeGroup", "-f", "dummyConfigFile.yaml"),
			Entry("without --nodes-min flags", "nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--nodes", "2", "--nodes-max", "3"),
			Entry("without --nodes-max flags", "nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--nodes", "2", "--nodes-min", "1"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				cmd := newDefaultCmd(c.args...)
				_, err := cmd.execute()
				fmt.Println(err)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(c.error.Error()))
			},
			Entry("missing required flag --cluster", invalidParamsCase{
				args:  []string{"nodegroup", "ng"},
				error: fmt.Errorf("--cluster must be set"),
			}),
			Entry("setting --name and argument at the same time", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--name", "ng"},
				error: fmt.Errorf("--name=ng and argument ng cannot be used at the same time"),
			}),
			Entry("missing required nodes flag --nodes", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy"},
				error: fmt.Errorf("number of nodes must be 0 or greater"),
			}),
			Entry("invalid flag", invalidParamsCase{
				args:  []string{"nodegroup", "--invalid", "dummy"},
				error: fmt.Errorf("unknown flag: --invalid"),
			}),
			Entry("desired node less than min nodes", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--nodes", "2", "--nodes-min", "3"},
				error: fmt.Errorf("minimum number of nodes must be less than or equal to number of nodes"),
			}),

			Entry("desired node greater than max nodes", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--nodes", "2", "--nodes-max", "1"},
				error: fmt.Errorf("maximum number of nodes must be greater than or equal to number of nodes"),
			}),
			Entry("desired node less than min nodes", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--nodes", "2", "--nodes-min", "3"},
				error: fmt.Errorf("minimum number of nodes must be less than or equal to number of nodes"),
			}),
			Entry("desired node outside the range [min, max]", invalidParamsCase{
				args:  []string{"nodegroup", "ng", "--cluster", "dummy", "--nodes", "2", "--nodes-min", "1", "--nodes-max", "1"},
				error: fmt.Errorf("number of nodes must be within range of min nodes and max nodes"),
			}),
			Entry("with config file and no name flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml"},
				error: fmt.Errorf("--name must be set"),
			}),
			Entry("with config file and nodes flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--nodes", "2"},
				error: fmt.Errorf("cannot use --nodes when --config-file/-f is set"),
			}),
			Entry("with config file and nodes-max flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--nodes-max", "2"},
				error: fmt.Errorf("cannot use --nodes-max when --config-file/-f is set"),
			}),
			Entry("with config file and nodes-min flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--nodes-min", "2"},
				error: fmt.Errorf("cannot use --nodes-min when --config-file/-f is set"),
			}),
			Entry("with config file and cluster flags", invalidParamsCase{
				args:  []string{"nodegroup", "-f", "../cmdutils/test_data/scale-ng-test.yaml", "--cluster", "dummyCluster"},
				error: fmt.Errorf("cannot use --cluster when --config-file/-f is set"),
			}),
		)
	})
})
