package get

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("get", func() {
	Describe("nodegroup", func() {
		It("missing required flag --cluster", func() {
			cmd := newMockCmd("nodegroup")
			_, err := cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: --cluster must be set")))
		})

		It("setting --name and argument at the same time", func() {
			cmd := newMockCmd("nodegroup", "ng", "--cluster", "dummy", "--name", "ng")
			_, err := cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: --name=ng and argument ng cannot be used at the same time")))
		})

		It("setting --name and --config-file at the same time", func() {
			f, err := os.CreateTemp("", "configfile")
			Expect(err).NotTo(HaveOccurred())
			_, err = f.WriteString(nodegroupConfigFile)
			Expect(err).NotTo(HaveOccurred())
			cmd := newMockCmd("nodegroup", "--name", "name", "--config-file", f.Name())
			_, err = cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: cannot use --name when --config-file/-f is set")))
		})

		It("setting --cluster and --config-file at the same time", func() {
			f, err := os.CreateTemp("", "configfile")
			Expect(err).NotTo(HaveOccurred())
			_, err = f.WriteString(nodegroupConfigFile)
			Expect(err).NotTo(HaveOccurred())
			cmd := newMockCmd("nodegroup", "--cluster", "name", "--config-file", f.Name())
			_, err = cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: cannot use --cluster when --config-file/-f is set")))
		})

		It("invalid flag", func() {
			cmd := newMockCmd("nodegroup", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: unknown flag: --invalid")))
		})
	})
})

var nodegroupConfigFile = `apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test-nodegroup-cluster-config
  region: us-west-2

managedNodeGroups:
  - name: managed-ng-1
    minSize: 1
    maxSize: 2
    desiredCapacity: 1
  - name: managed-ng-2
    minSize: 1
    maxSize: 2
    desiredCapacity: 1
`
