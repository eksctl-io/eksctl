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
	})

	Context("SSH settings", func() {

		It("Providing an SSH key enables SSH when SSH.Allow not set", func() {
			testKeyPath := "some/path/to/file.pub"

			testNodeGroup := NodeGroup{
				VolumeSize: &DefaultNodeVolumeSize,
				SSH: &NodeGroupSSH{
					PublicKeyPath: &testKeyPath,
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(*testNodeGroup.SSH.Allow).To(BeTrue())
			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo(testKeyPath))
		})

		It("Enabling SSH without a key uses default key", func() {
			testNodeGroup := NodeGroup{
				VolumeSize: &DefaultNodeVolumeSize,
				SSH: &NodeGroupSSH{
					Allow: Enabled(),
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Region: "eu-west-2"})

			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo("~/.ssh/id_rsa.pub"))
		})

		It("Providing an SSH key and explicitly disabling SSH keeps SSH disabled", func() {
			testKeyPath := "some/path/to/file.pub"

			testNodeGroup := NodeGroup{
				VolumeSize: &DefaultNodeVolumeSize,
				SSH: &NodeGroupSSH{
					Allow:         Disabled(),
					PublicKeyPath: &testKeyPath,
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Region: "us-east-1"})

			Expect(*testNodeGroup.SSH.Allow).To(BeFalse())
			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo(testKeyPath))
		})
	})

	Context("Bottlerocket Settings", func() {
		It("enables SSH with NodeGroup", func() {
			testNodeGroup := NodeGroup{
				AMIFamily: NodeImageFamilyBottlerocket,
				SSH: &NodeGroupSSH{
					Allow: Enabled(),
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(*testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeTrue())
		})

		It("has default NodeGroup configuration", func() {
			testNodeGroup := NodeGroup{
				AMIFamily: NodeImageFamilyBottlerocket,
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(testNodeGroup.Bottlerocket).ToNot(BeNil())
			Expect(testNodeGroup.AMI).To(Equal(NodeImageResolverAutoSSM))
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

	Context("KubletExtraConfigDefaults", func() {

		It("Should calculate correct memory size and cpu cores for instance types", func() {
			expRes := []struct {
				instType      string
				expMemRes     string
				expCPURes     string
				expEphStorRes string
			}{
				{"t3.nano", "0.125", "100", "1.25"},
				{"t3a.micro", "0.25", "100", "1.25"},
				{"t2.small", "0.5", "60", "1.25"},
				{"t2.medium", "1.0", "100", "1.25"},
				{"m5ad.large", "1.8", "100", "4.69"},
				{"m5ad.xlarge", "2.6", "140", "9.38"},
				{"m5ad.2xlarge", "3.56", "180", "15"},
				{"m5ad.8xlarge", "9.32", "420", "15"},
				{"m5ad.12xlarge", "10.6", "580", "15"},
				{"m5ad.16xlarge", "11.88", "740", "15"},
				{"m5ad.24xlarge", "14.44", "1040", "15"},
			}

			for _, m := range expRes {
				it := m.instType
				expMemRes := m.expMemRes
				expCPURes := m.expCPURes
				expEphStorRes := m.expEphStorRes

				testNodeGroup := NodeGroup{
					VolumeSize:         &DefaultNodeVolumeSize,
					InstanceType:       it,
					KubeletExtraConfig: &InlineDocument{"kubeReserved": make(map[string]interface{})},
				}
				SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Region: DefaultRegion})
				KEC := *testNodeGroup.KubeletExtraConfig
				kubeReserved := KEC["kubeReserved"].(map[string]interface{})
				Expect(kubeReserved["memory"]).To(Equal(expMemRes + "Mi"))
				Expect(kubeReserved["cpu"]).To(Equal(expCPURes + "m"))
				Expect(kubeReserved["ephemeral-storage"]).To(Equal(expEphStorRes + "Gi"))
			}
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
