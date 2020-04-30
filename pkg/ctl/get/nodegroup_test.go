package get

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get", func() {
	Describe("nodegroup", func() {
		It("missing required flag --cluster", func() {
			cmd := newMockCmd("nodegroup")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("setting --name and argument at the same time", func() {
			cmd := newMockCmd("nodegroup", "ng", "--cluster", "dummy", "--name", "ng")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--name=ng and argument ng cannot be used at the same time"))
		})

		It("invalid flag", func() {
			cmd := newMockCmd("nodegroup", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
