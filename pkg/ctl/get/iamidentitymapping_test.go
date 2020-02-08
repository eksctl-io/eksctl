package get

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get", func() {
	Describe("iamidentitymapping", func() {
		It("missing required flag --cluster", func() {
			cmd := newMockCmd( "iamidentitymapping")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("invalid flag --dummy", func() {
			cmd := newMockCmd("iamidentitymapping", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
