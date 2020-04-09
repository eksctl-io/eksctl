package get

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get", func() {
	Describe("labels", func() {
		It("missing required flag --cluster", func() {
			cmd := newMockCmd("labels")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("missing required flag --cluster, but with --nodegroup", func() {
			cmd := newMockCmd("labels", "--nodegroup", "dummyNodeGroup")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("setting name argument", func() {
			cmd := newMockCmd("labels", "--cluster", "dummy", "dummyName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name argument is not supported"))
		})
	})
})
