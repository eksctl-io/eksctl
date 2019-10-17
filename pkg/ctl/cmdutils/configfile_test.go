package cmdutils_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
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
				Expect(err).ToNot(HaveOccurred())
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
				Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile(`name argument "foo-3"`).Error()))

				fs := cmd.CobraCommand.Flags()

				fs.StringVar(&cfg.Metadata.Name, "cluster", "", "")
				cmd.CobraCommand.Flag("cluster").Changed = true

				Expect(cmd.CobraCommand.Flag("cluster").Changed).To(BeTrue())

				err = NewMetadataLoader(cmd).Load()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile("--cluster").Error()))
			}
		})

		It("load all of example file", func() {
			examples, err := filepath.Glob(examplesDir + "*.yaml")
			Expect(err).ToNot(HaveOccurred())

			Expect(examples).To(HaveLen(14))
			for _, example := range examples {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: example,
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    &api.ProviderConfig{},
				}

				err := NewMetadataLoader(cmd).Load()

				cfg := cmd.ClusterConfig
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.Metadata.Name).ToNot(BeEmpty())
				Expect(cfg.Metadata.Region).ToNot(BeEmpty())
				Expect(cfg.Metadata.Region).To(Equal(cmd.ProviderConfig.Region))
				Expect(cfg.Metadata.Version).To(BeEmpty())
			}
		})

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
					ProviderConfig:    &api.ProviderConfig{},
				}

				Expect(NewCreateClusterLoader(cmd, NewNodeGroupFilter(), nil, true).Load()).To(Succeed())
				cfg := cmd.ClusterConfig
				Expect(cfg.VPC.NAT.Gateway).To(Not(BeNil()))
				Expect(*cfg.VPC.NAT.Gateway).To(Equal(natTest.expectedGateway))
			}
		})

		It("loader should handle named and unnamed nodegroups without config file", func() {
			unnamedNG := api.NewNodeGroup()

			namedNG := api.NewNodeGroup()
			namedNG.Name = "ng-1"

			loaderParams := []struct {
				ng               *api.NodeGroup
				withoutNodeGroup bool
			}{
				{unnamedNG, false},
				{unnamedNG, true},
				{namedNG, false},
				{namedNG, false},
			}

			for _, loaderTest := range loaderParams {
				cmd := &Cmd{
					CobraCommand:   newCmd(),
					ClusterConfig:  api.NewClusterConfig(),
					ProviderConfig: &api.ProviderConfig{},
				}

				ngFilter := NewNodeGroupFilter()

				Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(0))

				Expect(NewCreateClusterLoader(cmd, ngFilter, loaderTest.ng, loaderTest.withoutNodeGroup).Load()).To(Succeed())

				Expect(ngFilter.ExcludeAll).To(Equal(loaderTest.withoutNodeGroup))

				if loaderTest.withoutNodeGroup {
					Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(0))
				}
				if !loaderTest.withoutNodeGroup {
					Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(1))
					Expect(cmd.ClusterConfig.NodeGroups[0]).To(Equal(loaderTest.ng))
				}

				_, err := cmd.NewCtl()
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("loader should handle nodegroup exclusion with config file", func() {
			loaderParams := []struct {
				configFile       string
				nodeGroupCount   int
				withoutNodeGroup bool
			}{
				{"01-simple-cluster.yaml", 1, true},
				{"01-simple-cluster.yaml", 1, false},
				{"02-custom-vpc-cidr-no-nodes.yaml", 0, true},
				{"02-custom-vpc-cidr-no-nodes.yaml", 0, true},
				{"03-two-nodegroups.yaml", 2, true},
				{"03-two-nodegroups.yaml", 2, false},
				{"05-advanced-nodegroups.yaml", 3, true},
				{"05-advanced-nodegroups.yaml", 3, false},
				{"07-ssh-keys.yaml", 5, true},
				{"07-ssh-keys.yaml", 5, false},
			}

			for _, loaderTest := range loaderParams {
				cmd := &Cmd{
					CobraCommand:      newCmd(),
					ClusterConfigFile: filepath.Join(examplesDir, loaderTest.configFile),
					ClusterConfig:     api.NewClusterConfig(),
					ProviderConfig:    &api.ProviderConfig{},
				}

				ngFilter := NewNodeGroupFilter()

				Expect(NewCreateClusterLoader(cmd, ngFilter, nil, loaderTest.withoutNodeGroup).Load()).To(Succeed())

				Expect(ngFilter.ExcludeAll).To(Equal(loaderTest.withoutNodeGroup))

				Expect(cmd.ClusterConfig.NodeGroups).To(HaveLen(loaderTest.nodeGroupCount))

				_, err := cmd.NewCtl()
				Expect(err).ToNot(HaveOccurred())
			}

		})
	})
})
