package update

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("update nodegroup", func() {
	It("returns error if cluster is not set", func() {
		cmd := newMockCmd("nodegroup")
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("--cluster must be set")))
	})

	It("returns error if nodegroup name is not set", func() {
		cmd := newMockCmd("nodegroup", "--cluster", "cluster-name")
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("--name must be set")))
	})

	It("returns error if config file is not set", func() {
		cmd := newMockCmd("nodegroup", "--cluster", "cluster-name", "--name", "ng-name")
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("--config-file must be set")))
	})
})
