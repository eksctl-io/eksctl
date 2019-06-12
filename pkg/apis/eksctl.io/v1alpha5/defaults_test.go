package v1alpha5

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default settings", func() {
	var (
		testKeyPath = "some/path/to/file.pub"
	)

	Context("SSH settings", func() {

		It("Providing an SSH key enables SSH", func() {
			testNodeGroup := NodeGroup{
				VolumeSize: &DefaultNodeVolumeSize,
				SSH: &NodeGroupSSH{
					Allow:         Disabled(),
					PublicKeyPath: &testKeyPath,
				},
			}

			SetNodeGroupDefaults(0, &testNodeGroup)

			Expect(*testNodeGroup.SSH.Allow).To(BeTrue())
		})

		It("Enabling SSH without a key uses default key", func() {
			testNodeGroup := NodeGroup{
				VolumeSize: &DefaultNodeVolumeSize,
				SSH: &NodeGroupSSH{
					Allow: Enabled(),
				},
			}

			SetNodeGroupDefaults(0, &testNodeGroup)

			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo("~/.ssh/id_rsa.pub"))
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
