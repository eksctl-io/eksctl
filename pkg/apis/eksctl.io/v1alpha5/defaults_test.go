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

	Context("Cluster NAT settings", func() {

		It("Cluster NAT defaults to single NAT gateway mode", func() {
			testVpc := &ClusterVPC{}
			testVpc.NAT = DefaultClusterNAT()

			Expect(*testVpc.NAT.Gateway).To(BeIdenticalTo(ClusterSingleNAT))

		})

	})

})
