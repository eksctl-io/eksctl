package cmdutils

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type scaleNodeGroupCase struct {
	name    string
	err     error
	minSize *int
}

type scaleNodeGroupCLICase struct {
	name        string
	err         error
	minSize     *int
	maxSize     *int
	desiredSize *int
}

var _ = Describe("scale node group config file loader", func() {
	newCmd := func() *cobra.Command {
		return &cobra.Command{
			Use: "test",
			Run: func(_ *cobra.Command, _ []string) {},
		}
	}

	When("a nodegroup name is passed with --name", func() {
		Describe("for unmanaged nodegroups", func() {
			DescribeTable("scales a single nodegroup successfully via config file",
				func(params scaleNodeGroupCase) {
					cmd := &Cmd{
						CobraCommand:      newCmd(),
						ClusterConfigFile: "test_data/scale-ng-test.yaml",
						ClusterConfig:     api.NewClusterConfig(),
						ProviderConfig:    api.ProviderConfig{},
						NameArg:           params.name,
					}

					ng := api.NewNodeGroup().BaseNodeGroup()
					err := NewScaleNodeGroupLoader(cmd, ng).Load()
					if params.err != nil {
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring(params.err.Error()))
					} else {
						if params.minSize != nil {
							Expect(ng.MinSize).To(Equal(params.minSize))
						}
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("one node group matched", scaleNodeGroupCase{
					name: "ng-all-details",
				}),
				Entry("no node group matched", scaleNodeGroupCase{
					name: "123123",
					err:  fmt.Errorf("nodegroup 123123 not found in config file"),
				}),
				Entry("with no desired capacity", scaleNodeGroupCase{
					name: "ng-no-desired-capacity",
					err:  fmt.Errorf("number of nodes must be 0 or greater"),
				}),
				Entry("with no minSize and no maxSize", scaleNodeGroupCase{
					name: "ng-no-min-max",
				}),
				Entry("ng with minSize", scaleNodeGroupCase{
					name:    "ng-with-min",
					minSize: aws.Int(1),
				}),
				Entry("ng with wrong value for minSize", scaleNodeGroupCase{
					name: "ng-with-wrong-min",
					err:  fmt.Errorf("minimum number of nodes must be less than or equal to number of nodes"),
				}),
				Entry("ng with maxSize", scaleNodeGroupCase{
					name: "ng-with-max",
				}),
				Entry("ng with wrong value for maxSize", scaleNodeGroupCase{
					name: "ng-with-wrong-max",
					err:  fmt.Errorf("maximum number of nodes must be greater than or equal to number of nodes"),
				}),
				Entry("ng with desired nodes outside [minSize, maxSize]", scaleNodeGroupCase{
					name: "ng-with-wrong-desired",
					err:  fmt.Errorf("number of nodes must be within range of min nodes and max nodes"),
				}),
			)

			DescribeTable("scales a single nodegroup successfully via cli flags",
				func(params scaleNodeGroupCLICase) {
					cfg := api.NewClusterConfig()
					cfg.Metadata.Name = "cluster"
					cmd := &Cmd{
						CobraCommand:   newCmd(),
						ProviderConfig: api.ProviderConfig{},
						ClusterConfig:  cfg,
						NameArg:        params.name,
					}

					ng := api.NewNodeGroup().BaseNodeGroup()
					ng.MinSize = params.minSize
					ng.MaxSize = params.maxSize
					ng.DesiredCapacity = params.desiredSize
					err := NewScaleNodeGroupLoader(cmd, ng).Load()
					if params.err != nil {
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring(params.err.Error()))
					} else {
						if params.minSize != nil {
							Expect(ng.MinSize).To(Equal(params.minSize))
						}
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("only specifying min-nodes", scaleNodeGroupCLICase{
					name:    "ng-with-max",
					minSize: aws.Int(1),
				}),
				Entry("only specifying max-nodes", scaleNodeGroupCLICase{
					name:    "ng-with-max",
					maxSize: aws.Int(1),
				}),
				Entry("only specifying nodes", scaleNodeGroupCLICase{
					name:        "ng-with-max",
					desiredSize: aws.Int(1),
				}),
				Entry("minSize 0", scaleNodeGroupCLICase{
					name:    "ng-with-max",
					minSize: aws.Int(-1),
					err:     fmt.Errorf("minimum of nodes must be 0 or greater"),
				}),
				Entry("maxSize 0", scaleNodeGroupCLICase{
					name:    "ng-with-max",
					maxSize: aws.Int(-1),
					err:     fmt.Errorf("maximum of nodes must be 0 or greater"),
				}),
				Entry("desiredSize 0", scaleNodeGroupCLICase{
					name:        "ng-with-max",
					desiredSize: aws.Int(-1),
					err:         fmt.Errorf("number of nodes must be 0 or greater"),
				}),
				Entry("desiredSize greater than max", scaleNodeGroupCLICase{
					name:        "ng-with-max",
					desiredSize: aws.Int(3),
					maxSize:     aws.Int(1),
					err:         fmt.Errorf("maximum number of nodes must be greater than or equal to number of nodes"),
				}),
				Entry("desiredSize fewer than min", scaleNodeGroupCLICase{
					name:        "ng-with-max",
					desiredSize: aws.Int(2),
					minSize:     aws.Int(3),
					err:         fmt.Errorf("minimum number of nodes must be fewer than or equal to number of nodes"),
				}),
				Entry("min greater than max", scaleNodeGroupCLICase{
					name:    "ng-with-max",
					minSize: aws.Int(3),
					maxSize: aws.Int(2),
					err:     fmt.Errorf("maximum number of nodes must be greater than minimum number of nodes"),
				}),
				Entry("not specifying any", scaleNodeGroupCLICase{
					name: "ng-with-max",
					err:  fmt.Errorf("at least one of minimum, maximum and desired nodes must be set"),
				}),
			)
		})

		Describe("for managed nodegroups", func() {
			When("passing a config file", func() {
				var (
					config *api.ClusterConfig
					ng     *api.NodeGroupBase
				)

				BeforeEach(func() {
					config = api.NewClusterConfig()
					config.Metadata.Name = "test-cluster"
					ng = config.NewNodeGroup().BaseNodeGroup()
				})

				When("the nodegroup is in the config file", func() {
					It("finds and scales that individual nodegroup", func() {
						ngName := "mng-ng"
						ng.Name = ngName
						ng.ScalingConfig = &api.ScalingConfig{
							DesiredCapacity: aws.Int(2),
						}

						config.ManagedNodeGroups = []*api.ManagedNodeGroup{
							{
								NodeGroupBase: &api.NodeGroupBase{
									Name: ngName,
									ScalingConfig: &api.ScalingConfig{
										DesiredCapacity: aws.Int(2),
									},
								},
							},
						}
						cmd := &Cmd{
							CobraCommand:   newCmd(),
							ClusterConfig:  config,
							ProviderConfig: api.ProviderConfig{},
						}

						err := NewScaleNodeGroupLoader(cmd, ng).Load()
						Expect(err).NotTo(HaveOccurred())
					})
				})

				When("the nodegroup is not in the config file", func() {
					It("fails to find and scale that nodegroup", func() {
						cmd := &Cmd{
							CobraCommand:   newCmd(),
							ClusterConfig:  config,
							ProviderConfig: api.ProviderConfig{},
						}

						err := NewScaleNodeGroupLoader(cmd, ng).Load()
						Expect(err).NotTo(BeNil())
					})
				})
			})
		})
	})

	When("a nodegroup name is not passed with --name", func() {
		When("using a config file with nodegroups", func() {
			It("finds and scales all nodegroups", func() {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: "test_data/scale-all-ng-test.yaml",
				}

				err := NewScaleAllNodeGroupLoader(cmd).Load()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(cmd.ClusterConfig.ManagedNodeGroups)).To(Equal(1))
				Expect(*cmd.ClusterConfig.ManagedNodeGroups[0].BaseNodeGroup().ScalingConfig.DesiredCapacity).To(Equal(3))
				Expect(*cmd.ClusterConfig.ManagedNodeGroups[0].BaseNodeGroup().ScalingConfig.MinSize).To(Equal(3))
				Expect(*cmd.ClusterConfig.ManagedNodeGroups[0].BaseNodeGroup().ScalingConfig.MaxSize).To(Equal(5))

				Expect(len(cmd.ClusterConfig.NodeGroups)).To(Equal(1))
				Expect(*cmd.ClusterConfig.NodeGroups[0].BaseNodeGroup().ScalingConfig.DesiredCapacity).To(Equal(3))
				Expect(*cmd.ClusterConfig.NodeGroups[0].BaseNodeGroup().ScalingConfig.MinSize).To(Equal(3))
				Expect(*cmd.ClusterConfig.NodeGroups[0].BaseNodeGroup().ScalingConfig.MaxSize).To(Equal(5))
			})
		})

		When("using a config file with no nodegroups", func() {
			It("fails to find find any nodegroup", func() {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: "../../../examples/02-custom-vpc-cidr-no-nodes.yaml",
				}

				err := NewScaleAllNodeGroupLoader(cmd).Load()
				Expect(err).To(MatchError("no nodegroups found in config file"))
			})
		})

		When("a config file is not specified", func() {
			It("fails to find find any nodegroups", func() {
				cmd := &Cmd{
					CobraCommand: newCmd(),
				}

				err := NewScaleAllNodeGroupLoader(cmd).Load()
				Expect(err).To(MatchError("a config file is required when --name is not set or when scaling multiple nodegroups"))
			})
		})
	})
})
