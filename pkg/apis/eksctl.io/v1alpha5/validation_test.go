package v1alpha5

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigFile ssh flags validation", func() {
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

func checkItDetectsError(SSHConfig *NodeGroupSSH) {
	err := validateNodeGroupSSH(SSHConfig)
	Expect(err).To(HaveOccurred())
}

func newInt(value int) *int {
	v := value
	return &v
}
