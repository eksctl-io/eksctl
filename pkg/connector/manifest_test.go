package connector

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest template", func() {
	Context("GetManifestTemplate", func() {
		It("should fetch the template", func() {
			template, err := GetManifestTemplate()
			Expect(err).ToNot(HaveOccurred())

			Expect(template.Connector).ToNot(BeEmpty())
			Expect(template.RoleBinding).ToNot(BeEmpty())
		})
	})
})
