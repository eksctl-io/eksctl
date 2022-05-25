package get

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("get", func() {
	Describe("iamserviceaccount", func() {
		It("missing required flag --cluster", func() {
			cmd := newMockCmd("iamserviceaccount")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --cluster must be set"))
		})

		It("invalid flag --dummy", func() {
			cmd := newMockCmd("iamserviceaccount", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: unknown flag: --invalid"))
		})
	})
})
