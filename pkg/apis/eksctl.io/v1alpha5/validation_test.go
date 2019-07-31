package v1alpha5

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigFile validation", func() {
	Describe("ssh flags", func() {
		var (
			testKeyPath = "some/path/to/file.pub"
			testKeyName = "id_rsa.pub"
			testKey     = "THIS IS A KEY"
		)

		It("fails when a key path and a key name are specified", func() {
			SSHConfig := &NodeGroupSSH{
				Allow:         Enabled(),
				PublicKeyPath: &testKeyPath,
				PublicKeyName: &testKeyName,
			}

			checkItDetectsError(SSHConfig)
		})

		It("fails when a key path and a key are specified", func() {
			SSHConfig := &NodeGroupSSH{
				Allow:         Enabled(),
				PublicKeyPath: &testKeyPath,
				PublicKey:     &testKey,
			}

			checkItDetectsError(SSHConfig)
		})

		It("fails when a key name and a key are specified", func() {
			SSHConfig := &NodeGroupSSH{
				Allow:         Enabled(),
				PublicKeyName: &testKeyName,
				PublicKey:     &testKey,
			}

			checkItDetectsError(SSHConfig)
		})

		It("It doesn't panic when SSH is not enabled", func() {
			err := validateNodeGroupSSH(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Instances distribution", func() {

			var ng *NodeGroup
			BeforeEach(func() {
				ng = &NodeGroup{
					InstancesDistribution: &NodeGroupInstancesDistribution{
						InstanceTypes: []string{"t3.medium", "t3.large"},
					},
				}
			})

			It("It doesn't panic when instance distribution is not enabled", func() {
				ng.InstancesDistribution = nil
				err := validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when instance distribution is enabled and instanceType is not empty or \"mixed\"", func() {
				err := validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())

				ng.InstanceType = "t3.small"

				err = validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())
			})

			It("It fails when the instance distribution doesn't have at least 2 different instance types", func() {
				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{"t3.medium", "t3.medium"}

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{"t3.medium", "t3.small"}

				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when the onDemandBaseCapacity is not above 0", func() {
				ng.InstancesDistribution.OnDemandBaseCapacity = newInt(-1)

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.OnDemandBaseCapacity = newInt(1)

				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when the spotInstancePools is not between 1 and 20", func() {
				ng.InstancesDistribution.SpotInstancePools = newInt(0)

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.SpotInstancePools = newInt(21)
				err = validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.SpotInstancePools = newInt(2)
				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when the onDemandPercentageAboveBaseCapacity is not between 0 and 100", func() {
				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(-1)

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(101)
				err = validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(50)
				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("kubelet extra config", func() {
		Context("Instances distribution", func() {

			var ng *NodeGroup
			BeforeEach(func() {
				ng = &NodeGroup{}
			})

			It("Forbids overriding basic fields", func() {
				testKeys := []string{"kind", "apiVersion", "address", "clusterDomain", "authentication",
					"authorization", "serverTLSBootstrap"}

				for _, key := range testKeys {
					ng.KubeletExtraConfig = &NodeGroupKubeletConfig{
						key: "should-not-be-allowed",
					}
					err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig)
					Expect(err).To(HaveOccurred())
				}
			})

			It("Allows other kubelet options", func() {
				ng.KubeletExtraConfig = &NodeGroupKubeletConfig{
					"kubeReserved": map[string]string{
						"cpu":               "300m",
						"memory":            "300Mi",
						"ephemeral-storage": "1Gi",
					},
					"kubeReservedCgroup": "/kube-reserved",
					"cgroupDriver":       "systemd",
					"featureGates": map[string]bool{
						"VolumeScheduling":         false,
						"VolumeSnapshotDataSource": true,
					},
				}
				err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig)
				Expect(err).ToNot(HaveOccurred())
			})

		})
	})

	Describe("ebs encryption", func() {
		var (
			nodegroup = "ng1"
			volSize   = 50
			kmsKeyID  = "36c0b54e-64ed-4f2d-a1c7-96558764311e"
			disabled  = false
			enabled   = true
		)

		Context("Encrypted workers", func() {

			var ng *NodeGroup
			BeforeEach(func() {
				ng = &NodeGroup{}
			})

			It("Forbids setting volumeKmsKeyID without volumeEncrypted", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = nil
				ng.VolumeKmsKeyID = &kmsKeyID
				err := ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
			})

			It("Forbids setting volumeKmsKeyID with volumeEncrypted: false", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = &disabled
				ng.VolumeKmsKeyID = &kmsKeyID
				err := ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
			})

			It("Allows setting volumeKmsKeyID with volumeEncrypted: true", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = &enabled
				ng.VolumeKmsKeyID = &kmsKeyID
				err := ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})

		})
	})

})

func checkItDetectsError(SSHConfig *NodeGroupSSH) {
	err := validateNodeGroupSSH(SSHConfig)
	Expect(err).To(HaveOccurred())
}

func newInt(value int) *int {
	v := value
	return &v
}
