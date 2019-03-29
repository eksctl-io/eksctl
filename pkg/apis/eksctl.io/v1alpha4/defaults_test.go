package v1alpha4

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
				SSH: SSHConfig{
					Allow:         false,
					PublicKeyPath: &testKeyPath,
				},
			}

			SetNodeGroupDefaults(0, &testNodeGroup)

			Expect(testNodeGroup.SSH.Allow).To(BeTrue())
		})

		It("Enabling SSH without a key uses default key", func() {
			testNodeGroup := NodeGroup{
				SSH: SSHConfig{
					Allow: true,
				},
			}

			SetNodeGroupDefaults(0, &testNodeGroup)

			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo("~/.ssh/id_rsa.pub"))
		})
	})
})
