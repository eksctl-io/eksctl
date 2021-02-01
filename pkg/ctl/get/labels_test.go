package get

import (
	. "github.com/onsi/ginkgo"
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
	})
})
