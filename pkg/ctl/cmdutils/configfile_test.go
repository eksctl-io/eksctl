package cmdutils

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
)

var _ = Describe("cmdutils configfile", func() {

	newCmd := func() *cobra.Command {
		return &cobra.Command{
			Use: "test",
			Run: func(_ *cobra.Command, _ []string) {},
		}
	}

	const examplesDir = "../../../examples/"

	Context("load configfiles", func() {

		It("should handle name argument", func() {
			cfg := api.NewClusterConfig()

			{
				cmd := &Cmd{
					ClusterConfig: cfg,
					NameArg:       "foo-1",
					CobraCommand:  newCmd(),
				}

				err := NewMetadataLoader(cmd).Load()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Metadata.Name).To(Equal("foo-1"))
			}

			{
				cmd := &Cmd{
					ClusterConfig: cfg,
					NameArg:       "foo-2",
					CobraCommand:  newCmd(),
				}

				err := NewMetadataLoader(cmd).Load()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("--cluster=foo-1 and argument foo-2 cannot be used at the same time"))
			}

			{
				cmd := &Cmd{
					ClusterConfig:     cfg,
					NameArg:           "foo-3",
					CobraCommand:      newCmd(),
					ClusterConfigFile: examplesDir + "01-simple-cluster.yaml",
				}

				err := NewMetadataLoader(cmd).Load()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile(`name argument`).Error()))

				fs := cmd.CobraCommand.Flags()

				fs.StringVar(&cfg.Metadata.Name, "cluster", "", "")
				cmd.CobraCommand.Flag("cluster").Changed = true

				Expect(cmd.CobraCommand.Flag("cluster").Changed).To(BeTrue())

				err = NewMetadataLoader(cmd).Load()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile("--cluster").Error()))
			}
		})

		Describe("name argument", func() {
			When("given as --name switch", func() {
				It("succeeds", func() {
					cfg := api.NewClusterConfig()
					cobraCmd := newCmd()
					name := "foo-2"
					cobraCmd.SetArgs([]string{"--name", name})

					cmd := &Cmd{
						ClusterConfig:     cfg,
						NameArg:           name,
						CobraCommand:      cobraCmd,
						ClusterConfigFile: examplesDir + "01-simple-cluster.yaml",
						ProviderConfig:    api.ProviderConfig{},
					}
					l := newCommonClusterConfigLoader(cmd)
					l.flagsIncompatibleWithConfigFile.Delete("name")

					err := l.Load()
					Expect(err).NotTo(HaveOccurred())
				})
			})
			When("given as positional argument", func() {
				It("succeeds", func() {
					cfg := api.NewClusterConfig()
					cmd := &Cmd{
						ClusterConfig:     cfg,
						NameArg:           "foo-2",
						CobraCommand:      newCmd(),
						ClusterConfigFile: examplesDir + "01-simple-cluster.yaml",
						ProviderConfig:    api.ProviderConfig{},
					}
					l := newCommonClusterConfigLoader(cmd)
					l.flagsIncompatibleWithConfigFile.Delete("name")

					err := l.Load()
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Describe("managed nodegroup with container runtime", func() {
			When("container runtime is set for a managed nodegroup", func() {
				It("fails validation", func() {
					cfg := api.NewClusterConfig()
					cobraCmd := newCmd()
					name := "foo-2"
					cobraCmd.SetArgs([]string{"--name", name})

					cmd := &Cmd{
						ClusterConfig:     cfg,
						NameArg:           name,
						CobraCommand:      cobraCmd,
						ClusterConfigFile: "test_data/managed-nodegroup-with-container-runtime.yaml",
						ProviderConfig:    api.ProviderConfig{},
					}
					l := newCommonClusterConfigLoader(cmd)
					l.flagsIncompatibleWithConfigFile.Delete("name")

					err := l.Load()
					Expect(err).To(MatchError(ContainSubstring("error unmarshaling JSON: while decoding JSON: json: unknown field \"containerRuntime\"")))
				})
			})
		})

		It("load all of example file", func() {
			examples, err := filepath.Glob(examplesDir + "*.yaml")
			Expect(err).NotTo(HaveOccurred())

			Expect(examples).NotTo(BeEmpty())
			for _, example := range examples {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: example,
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    api.ProviderConfig{},
				}

				err := NewMetadataLoader(cmd).Load()
				Expect(err).NotTo(HaveOccurred())

				cfg := cmd.ClusterConfig
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Metadata.Name).NotTo(BeEmpty())
				Expect(cfg.Metadata.Region).NotTo(BeEmpty())
				Expect(cfg.Metadata.Region).To(Equal(cmd.ProviderConfig.Region))
			}
		})
	})

	Describe("CreateClusterLoader", func() {
		It("should set VPC.NAT.Gateway with the correct value", func() {
			natTests := []struct {
				configFile      string
				expectedGateway string
			}{
				// No VPC set
				{"01-simple-cluster.yaml", api.ClusterSingleNAT},
				// VPC set but not NAT
				{"02-custom-vpc-cidr-no-nodes.yaml", api.ClusterSingleNAT},
				// VPC and subnets set but not NAT
				{"04-existing-vpc.yaml", api.ClusterSingleNAT},
				// NAT set
				{"09-nat-gateways.yaml", api.ClusterHighlyAvailableNAT},
			}

			for _, natTest := range natTests {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: filepath.Join(examplesDir, natTest.configFile),
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    api.ProviderConfig{},
				}

				params := &CreateClusterCmdParams{WithoutNodeGroup: true, CreateManagedNGOptions: CreateManagedNGOptions{
					Managed: false,
				}}
				Expect(NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()).To(Succeed())
				cfg := cmd.ClusterConfig
				Expect(cfg.VPC.NAT.Gateway).To(Not(BeNil()))
				Expect(*cfg.VPC.NAT.Gateway).To(Equal(natTest.expectedGateway))
			}
		})

		It("should return an error for unset nodegroups", func() {
			cmd := &Cmd{
				CobraCommand:      newCmd(),
				ClusterConfigFile: filepath.Join("test_data", "unset-nodegroups.yaml"),
				ClusterConfig:     api.NewClusterConfig(),
				ProviderConfig:    api.ProviderConfig{},
			}
			params := &CreateClusterCmdParams{}
			Expect(NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()).To(MatchError("invalid ClusterConfig: nodeGroups[1] is not set"))
		})

		When("using zones and node-zones together", func() {
			var (
				flagSet *pflag.FlagSet
				testCmd *cobra.Command
			)
			BeforeEach(func() {
				flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
				flagSet.StringSlice("zones", []string{}, "Zones")
				flagSet.StringSlice("node-zones", []string{}, "Node zones")
				testCmd = newCmd()
				testCmd.Flags().AddFlagSet(flagSet)
			})
			When("node-zones is a subset of zones", func() {
				It("does not error", func() {
					Expect(testCmd.ParseFlags([]string{"--zones", "zone1,zone2", "--node-zones", "zone1"})).To(Succeed())
					cmd := &Cmd{
						CobraCommand:   testCmd,
						ClusterConfig:  api.NewClusterConfig(),
						ProviderConfig: api.ProviderConfig{},
					}
					params := &CreateClusterCmdParams{WithoutNodeGroup: true, CreateManagedNGOptions: CreateManagedNGOptions{
						Managed: false,
					}}
					Expect(NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()).To(Succeed())
				})
			})
			When("node-zones is not a subset of zones", func() {
				It("errors", func() {
					Expect(testCmd.ParseFlags([]string{"--zones", "zone1,zone2", "--node-zones", "zone3"})).To(Succeed())
					cmd := &Cmd{
						CobraCommand:   testCmd,
						ClusterConfig:  api.NewClusterConfig(),
						ProviderConfig: api.ProviderConfig{},
					}

					params := &CreateClusterCmdParams{WithoutNodeGroup: true, CreateManagedNGOptions: CreateManagedNGOptions{
						Managed: false,
					}}
					err := NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()
					Expect(err).To(MatchError(ContainSubstring("node-zones [zone3] must be a subset of zones [zone1 zone2]; \"zone3\" was not found in zones")))
				})
			})
			When("node-zones is defined but zones isn't", func() {
				It("errors", func() {
					Expect(testCmd.ParseFlags([]string{"--node-zones", "zone3"})).To(Succeed())
					cmd := &Cmd{
						CobraCommand:   testCmd,
						ClusterConfig:  api.NewClusterConfig(),
						ProviderConfig: api.ProviderConfig{},
					}

					params := &CreateClusterCmdParams{WithoutNodeGroup: true, CreateManagedNGOptions: CreateManagedNGOptions{
						Managed: false,
					}}
					err := NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()
					Expect(err).To(MatchError(ContainSubstring("zones must be defined if node-zones is provided and must be a superset of node-zones")))
				})
			})
		})

		When("using ipv6", func() {
			It("should default VPC.NAT to nil", func() {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: filepath.Join(examplesDir, "30-vpc-with-ip-family.yaml"),
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    api.ProviderConfig{},
				}
				params := &CreateClusterCmdParams{WithoutNodeGroup: true, CreateManagedNGOptions: CreateManagedNGOptions{
					Managed: false,
				}}
				Expect(NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()).To(Succeed())
				Expect(cmd.ClusterConfig.VPC.NAT).To(BeNil())
			})

			When("using pre-existing VPC", func() {
				It("should default VPC.NAT to nil", func() {
					cmd := &Cmd{
						CobraCommand:      newCmd(),
						ClusterConfigFile: filepath.Join("test_data", "ipv6-existing-vpc.yaml"),
						ClusterConfig:     api.NewClusterConfig(),
						ProviderConfig:    api.ProviderConfig{},
					}
					params := &CreateClusterCmdParams{WithoutNodeGroup: true, CreateManagedNGOptions: CreateManagedNGOptions{
						Managed: false,
					}}
					Expect(NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()).To(Succeed())
					Expect(cmd.ClusterConfig.VPC.NAT).To(BeNil())
				})
			})
		})

		It("loader should handle named and unnamed nodegroups without config file", func() {
			unnamedNG := api.NewNodeGroup()

			namedNG := api.NewNodeGroup()
			namedNG.Name = "ng-1"

			loaderParams := []struct {
				ng               *api.NodeGroup
				withoutNodeGroup bool
				managed          bool
			}{
				{unnamedNG, false, false},
				{unnamedNG, true, false},
				{namedNG, false, false},
				{namedNG, true, false},
				{unnamedNG, false, true},
				{unnamedNG, true, true},
				{namedNG, false, true},
			}

			assertMatchesNg := func(ng *api.NodeGroup, mNg *api.ManagedNodeGroup) {
				Expect(ng.Name).To(Equal(mNg.Name))
				Expect(ng.IAM).To(Equal(mNg.IAM))
				Expect(ng.VolumeSize).To(Equal(mNg.VolumeSize))
				Expect(ng.MinSize).To(Equal(mNg.MinSize))
				Expect(ng.MaxSize).To(Equal(mNg.MaxSize))
				Expect(ng.DesiredCapacity).To(Equal(mNg.DesiredCapacity))
				Expect(ng.AMIFamily).To(Equal(mNg.AMIFamily))
				Expect(ng.InstanceType).To(Equal(mNg.InstanceType))
				Expect(ng.Tags).To(Equal(mNg.Tags))
				Expect(ng.Labels).To(Equal(mNg.Labels))
				Expect(ng.AvailabilityZones).To(Equal(mNg.AvailabilityZones))
				Expect(ng.SSH).To(Equal(mNg.SSH))
			}

			for _, loaderTest := range loaderParams {
				cmd := &Cmd{
					CobraCommand:   newCmd(),
					ClusterConfig:  api.NewClusterConfig(),
					ProviderConfig: api.ProviderConfig{},
				}

				ngFilter := filter.NewNodeGroupFilter()

				Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(0))

				params := &CreateClusterCmdParams{
					WithoutNodeGroup: loaderTest.withoutNodeGroup,
					CreateManagedNGOptions: CreateManagedNGOptions{
						Managed: loaderTest.managed,
					},
				}
				Expect(NewCreateClusterLoader(cmd, ngFilter, loaderTest.ng, params).Load()).To(Succeed())

				Expect(ngFilter.GetExcludeAll()).To(Equal(loaderTest.withoutNodeGroup))

				if loaderTest.withoutNodeGroup {
					Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(0))
				} else {
					if loaderTest.managed {
						Expect(cmd.ClusterConfig.NodeGroups).To(BeEmpty())
						Expect(cmd.ClusterConfig.ManagedNodeGroups).To(HaveLen(1))
						assertMatchesNg(loaderTest.ng, cmd.ClusterConfig.ManagedNodeGroups[0])
						Expect(cmd.ClusterConfig.ManagedNodeGroups[0].Name).To(Equal(loaderTest.ng.Name))
					} else {
						Expect(cmd.ClusterConfig.ManagedNodeGroups).To(BeEmpty())
						Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(1))
						Expect(cmd.ClusterConfig.NodeGroups[0]).To(Equal(loaderTest.ng))
					}
				}
			}
		})

		It("loader should handle nodegroup exclusion with config file", func() {
			loaderParams := []struct {
				configFile       string
				nodeGroupCount   int
				withoutNodeGroup bool
				// determines whether this ClusterConfig contains managed nodegroups and not the managedFlag argument
				// which is not used when using a config file
				managed bool
			}{
				{"01-simple-cluster.yaml", 1, true, false},
				{"01-simple-cluster.yaml", 1, false, false},
				{"02-custom-vpc-cidr-no-nodes.yaml", 0, true, false},
				{"02-custom-vpc-cidr-no-nodes.yaml", 0, true, false},
				{"03-two-nodegroups.yaml", 2, true, false},
				{"03-two-nodegroups.yaml", 2, false, false},
				{"05-advanced-nodegroups.yaml", 3, true, false},
				{"05-advanced-nodegroups.yaml", 3, false, false},
				{"07-ssh-keys.yaml", 7, true, false},
				{"07-ssh-keys.yaml", 7, false, false},
				{"15-managed-nodes.yaml", 4, true, true},
				{"15-managed-nodes.yaml", 4, false, true},
				{"20-bottlerocket.yaml", 2, false, false},
			}

			for _, loaderTest := range loaderParams {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: filepath.Join(examplesDir, loaderTest.configFile),
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    api.ProviderConfig{},
				}

				ngFilter := filter.NewNodeGroupFilter()

				params := &CreateClusterCmdParams{
					WithoutNodeGroup: loaderTest.withoutNodeGroup,
					CreateManagedNGOptions: CreateManagedNGOptions{
						Managed: loaderTest.managed,
					},
				}
				Expect(NewCreateClusterLoader(cmd, ngFilter, nil, params).Load()).To(Succeed())

				Expect(ngFilter.GetExcludeAll()).To(Equal(loaderTest.withoutNodeGroup))

				if loaderTest.managed {
					Expect(cmd.ClusterConfig.ManagedNodeGroups).To(HaveLen(loaderTest.nodeGroupCount))
					Expect(cmd.ClusterConfig.NodeGroups).To(BeEmpty())
				} else {
					Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(loaderTest.nodeGroupCount))
					Expect(cmd.ClusterConfig.ManagedNodeGroups).To(BeEmpty())
				}
			}

		})

		Describe("should set defaults for cluster endpoint access", func() {

			testClusterEndpointAccessDefaults := func(configFilePath string, expectedPrivAccess, expectedPubAccess bool) {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: configFilePath,
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    api.ProviderConfig{},
				}

				params := &CreateClusterCmdParams{
					WithoutNodeGroup: true,
					CreateManagedNGOptions: CreateManagedNGOptions{
						Managed: false,
					},
				}

				Expect(NewCreateClusterLoader(cmd, filter.NewNodeGroupFilter(), nil, params).Load()).To(Succeed())
				cfg := cmd.ClusterConfig
				assertValidClusterEndpoint(cfg.VPC.ClusterEndpoints, expectedPrivAccess, expectedPubAccess)
			}

			It("when VPC is imported and no access is defined", func() {
				testClusterEndpointAccessDefaults("test_data/cluster-with-vpc.yaml", false, true)
			})

			It("when VPC is created by eksctl and no access is defined", func() {
				testClusterEndpointAccessDefaults("test_data/cluster-without-vpc.yaml", false, true)
			})

			It("when VPC is created by eksctl and private endpoint is enabled", func() {
				testClusterEndpointAccessDefaults("test_data/cluster-without-vpc-private-access.yaml", true, true)
			})

			It("when VPC is imported and private endpoint is enabled", func() {
				testClusterEndpointAccessDefaults("test_data/cluster-with-vpc-private-access.yaml", true, true)
			})
		})
	})

	Describe("SetLabelLoader", func() {
		It("should load the right data", func() {
			cmd := &Cmd{
				CobraCommand:      newCmd(),
				ClusterConfigFile: filepath.Join("test_data", "cluster-with-labels.yaml"),
				ClusterConfig:     api.NewClusterConfig(),
				ProviderConfig:    api.ProviderConfig{},
			}

			Expect(NewSetLabelLoader(cmd, "ng-1", nil).Load()).To(Succeed())
			cfg := cmd.ClusterConfig
			Expect(cfg.Metadata.Name).To(Equal("test-labels-1"))
			Expect(cfg.ManagedNodeGroups[0].Name).To(Equal("ng-1"))
		})
		When("no config file is provided and no labels", func() {
			It("errors", func() {
				cmd := &Cmd{
					CobraCommand:   newCmd(),
					ClusterConfig:  api.NewClusterConfig(),
					ProviderConfig: api.ProviderConfig{},
				}

				err := NewSetLabelLoader(cmd, "", nil).Load()
				Expect(err).To(MatchError(ContainSubstring("--labels must be set")))
			})
		})
		When("config file with no nodegroup name is provided", func() {
			It("should load all nodegroups", func() {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: filepath.Join("test_data", "cluster-with-labels.yaml"),
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    api.ProviderConfig{},
				}

				Expect(NewSetLabelLoader(cmd, "", nil).Load()).To(Succeed())
				cfg := cmd.ClusterConfig
				Expect(cfg.Metadata.Name).To(Equal("test-labels-1"))
				Expect(cfg.ManagedNodeGroups[0].Name).To(Equal("ng-1"))
				Expect(cfg.ManagedNodeGroups[1].Name).To(Equal("ng-2"))
			})
		})
		When("config file with missing nodegroup name is provided", func() {
			It("errors", func() {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: filepath.Join("test_data", "cluster-with-labels.yaml"),
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    api.ProviderConfig{},
				}

				err := NewSetLabelLoader(cmd, "invalid", nil).Load()
				Expect(err).To(MatchError(ContainSubstring("nodegroup with name invalid not found in the config file")))
			})
		})
	})
})

func assertValidClusterEndpoint(endpoints *api.ClusterEndpoints, privateAccess, publicAccess bool) {
	Expect(endpoints).To(Not(BeNil()))
	Expect(endpoints.PrivateAccess).To(Not(BeNil()))
	Expect(*endpoints.PrivateAccess).To(Equal(privateAccess))
	Expect(endpoints.PublicAccess).To(Not(BeNil()))
	Expect(*endpoints.PublicAccess).To(Equal(publicAccess))
}
