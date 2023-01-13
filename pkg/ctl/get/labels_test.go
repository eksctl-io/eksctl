package get

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("get", func() {
	Describe("labels", func() {
		It("fails when no flags set", func() {
			cmd := newMockCmd("labels")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --cluster must be set"))
		})

		It("fails when --cluster flag not set", func() {
			cmd := newMockCmd("labels", "--nodegroup", "dummyNodeGroup")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --cluster must be set"))
		})

		It("fails when --nodegroup flag not set", func() {
			cmd := newMockCmd("labels", "--cluster", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --nodegroup must be set"))
		})

		It("fails when name argument is used", func() {
			cmd := newMockCmd("labels", "--cluster", "dummy", "--nodegroup", "dummyNodeGroup", "dummyName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: name argument is not supported"))
		})

		It("setting --cluster and --config-file at the same time", func() {
			f, err := os.CreateTemp("", "configfile")
			Expect(err).NotTo(HaveOccurred())
			_, err = f.WriteString(labelsGetConfigFile)
			Expect(err).NotTo(HaveOccurred())
			cmd := newMockCmd("labels", "--cluster", "name", "--nodegroup", "name", "--config-file", f.Name())
			_, err = cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("Error: cannot use --cluster when --config-file/-f is set")))
		})
	})
})

var labelsGetConfigFile = `apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test-nodegroup-cluster-config
  region: us-west-2
`
