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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

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
			Expect(*testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeFalse())
			Expect(testNodeGroup.Bottlerocket.Settings).ToNot(BeNil())
			settings := map[string]interface{}(*testNodeGroup.Bottlerocket.Settings)
			Expect(settings).To(HaveKey("kubernetes"))
			kube := settings["kubernetes"].(map[string]interface{})
			Expect(kube).To(HaveKey("node-labels"))
			Expect(kube).ToNot(HaveKey("node-taints"))
			Expect(kube).ToNot(HaveKey("max-pods"))
			Expect(kube).ToNot(HaveKey("cluster-dns-ip"))
		})

		It("reflects NodeGroup configuration", func() {
			testNodeGroup := NodeGroup{
				AMIFamily: NodeImageFamilyBottlerocket,
				SSH: &NodeGroupSSH{
					Allow: Enabled(),
				},
				Labels: map[string]string{
					"label": "label-value",
				},
				Taints: map[string]string{
					"taint": "taint-value",
				},
				ClusterDNS: "192.0.2.53",
				MaxPodsPerNode: 32,
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{})

			Expect(*testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeTrue())

			Expect(testNodeGroup.Bottlerocket.Settings).ToNot(BeNil())
			settings := map[string]interface{}(*testNodeGroup.Bottlerocket.Settings)

			Expect(settings).To(HaveKey("kubernetes"))
			kube := settings["kubernetes"].(map[string]interface{})

			Expect(kube).To(HaveKeyWithValue("cluster-dns-ip", BeEquivalentTo(testNodeGroup.ClusterDNS)))
			Expect(kube).To(HaveKeyWithValue("max-pods", BeEquivalentTo(testNodeGroup.MaxPodsPerNode)))
			Expect(kube).To(HaveKey("node-labels"))
			Expect(kube).To(HaveKey("node-taints"))


			labels, ok := kube["node-labels"].(map[string]interface{})
			Expect(ok).To(BeTrue(), "unexpected type for labels")
			for key, value := range testNodeGroup.Labels {
				Expect(labels).To(HaveKeyWithValue(key, value))
			}

			taints, ok := kube["node-taints"].(map[string]interface{})
			Expect(ok).To(BeTrue(), "unexpected type for taints")
			for key, value := range testNodeGroup.Taints {
				Expect(taints).To(HaveKeyWithValue(key, value))
			}
		})
	})

	Context("Cluster NAT settings", func() {

		It("Cluster NAT defaults to single NAT gateway mode", func() {
			testVpc := &ClusterVPC{}
			testVpc.NAT = DefaultClusterNAT()

			Expect(*testVpc.NAT.Gateway).To(BeIdenticalTo(ClusterSingleNAT))

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
