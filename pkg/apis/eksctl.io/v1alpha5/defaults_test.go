package v1alpha5

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClusterConfig validation", func() {
	Describe("cloudWatch.clusterLogging", func() {
		var (
			cfg *ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		It("should handle unknown types", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"anything"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})

		It("should have no logging by default", func() {
			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(BeEmpty())
		})

		It("should expand `['*']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"*"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})

		It("should expand `['all']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"all"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})

		It("should expand `['api', '*', 'audit']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"api", "*", "audit"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})

		It("should expand `['authenticator', 'controllermanager', 'all']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"authenticator", "controllermanager", "all"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})
	})

	Context("SSH settings", func() {

		It("Providing an SSH key enables SSH when SSH.Allow not set", func() {
			testKeyPath := "some/path/to/file.pub"

			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeSize: &DefaultNodeVolumeSize,
					SSH: &NodeGroupSSH{
						PublicKeyPath: &testKeyPath,
					},
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(*testNodeGroup.SSH.Allow).To(BeTrue())
			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo(testKeyPath))
		})

		It("Enabling SSH without a key uses default key", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeSize: &DefaultNodeVolumeSize,
					SSH: &NodeGroupSSH{
						Allow: Enabled(),
					},
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo("~/.ssh/id_rsa.pub"))
		})

		It("Providing an SSH key and explicitly disabling SSH keeps SSH disabled", func() {
			testKeyPath := "some/path/to/file.pub"

			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeSize: &DefaultNodeVolumeSize,
					SSH: &NodeGroupSSH{
						Allow:         Disabled(),
						PublicKeyPath: &testKeyPath,
					},
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(*testNodeGroup.SSH.Allow).To(BeFalse())
			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo(testKeyPath))
		})
	})

	Context("Bottlerocket Settings", func() {
		It("enables SSH with NodeGroup", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
					SSH: &NodeGroupSSH{
						Allow: Enabled(),
					},
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(*testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeTrue())
		})

		It("has default NodeGroup configuration", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(testNodeGroup.Bottlerocket).ToNot(BeNil())
			Expect(*testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeFalse())
		})
	})

	Context("Cluster NAT settings", func() {

		It("Cluster NAT defaults to single NAT gateway mode", func() {
			testVpc := &ClusterVPC{}
			testVpc.NAT = DefaultClusterNAT()

			Expect(*testVpc.NAT.Gateway).To(BeIdenticalTo(ClusterSingleNAT))

		})

	})

	Describe("Cluster Managed Shared Node Security Group settings", func() {
		var (
			cfg *ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		It("should be enabled by default", func() {
			SetClusterConfigDefaults(cfg)
			Expect(*cfg.VPC.ManageSharedNodeSecurityGroupRules).To(BeTrue())
		})

		It("should fail validation if disabled without a defined shared node security group", func() {
			cfg.VPC.ManageSharedNodeSecurityGroupRules = Disabled()
			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})

		It("should pass validation if disabled with a defined shared node security group", func() {
			cfg.VPC.SharedNodeSecurityGroup = "sg-0123456789"
			cfg.VPC.ManageSharedNodeSecurityGroupRules = Disabled()
			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Describe("ClusterConfig", func() {
		var cfg *ClusterConfig

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		Describe("SetDefaultFargateProfile", func() {
			It("should create a default Fargate profile with two selectors matching default and kube-system w/o any label", func() {
				Expect(cfg.FargateProfiles).To(HaveLen(0))
				cfg.SetDefaultFargateProfile()
				Expect(cfg.FargateProfiles).To(HaveLen(1))
				profile := cfg.FargateProfiles[0]
				Expect(profile.Name).To(Equal("fp-default"))
				Expect(profile.Selectors).To(HaveLen(2))
				Expect(profile.Selectors[0].Namespace).To(Equal("default"))
				Expect(profile.Selectors[0].Labels).To(HaveLen(0))
				Expect(profile.Selectors[1].Namespace).To(Equal("kube-system"))
				Expect(profile.Selectors[1].Labels).To(HaveLen(0))
			})
		})
	})
})
