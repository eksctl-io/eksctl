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

			assertManifestFile := func(m ManifestFile, expectedFilename string) {
				Expect(m.Data).ToNot(BeEmpty())
				Expect(m.Filename).To(Equal(expectedFilename))
			}
			assertManifestFile(template.Connector, "eks-connector.yaml")
			assertManifestFile(template.ClusterRole, "eks-connector-clusterrole.yaml")
			assertManifestFile(template.ConsoleAccess, "eks-connector-console-dashboard-full-access-group.yaml")
		})
	})
})
