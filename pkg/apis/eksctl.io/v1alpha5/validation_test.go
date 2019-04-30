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
})

func checkItDetectsError(SSHConfig *NodeGroupSSH) {
	err := validateNodeGroupSSH(SSHConfig)
	Expect(err).To(HaveOccurred())
}
