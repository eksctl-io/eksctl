package get

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get", func() {
	Describe("cluster", func() {
		It("with invalid flags", func() {
			cmd := newMockCmd("cluster", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: unknown flag: --invalid")))
		})
		It("--name and --config-file together", func() {
			f, err := os.CreateTemp("", "configfile")
			Expect(err).NotTo(HaveOccurred())
			_, err = f.WriteString(getClusterConfigFile)
			Expect(err).NotTo(HaveOccurred())
			cmd := newMockCmd("cluster", "--name", "dummy", "--config-file", f.Name())
			_, err = cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: cannot use --name when --config-file/-f is set")))
		})
	})
})

var getClusterConfigFile = `apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: test-nodegroup-cluster-config
  region: us-west-2
  version: '1.20'
`
