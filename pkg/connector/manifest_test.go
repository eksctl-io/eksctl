package connector_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/connector"
)

var _ = Describe("Manifest template", func() {
	Context("GetManifestTemplate", func() {
		It("should fetch the template", func() {
			template, err := connector.GetManifestTemplate()
			Expect(err).NotTo(HaveOccurred())

			assertManifestFile := func(m connector.ManifestFile, expectedFilename string) {
				Expect(m.Data).NotTo(BeEmpty())
				Expect(m.Filename).To(Equal(expectedFilename))
			}
			assertManifestFile(template.Connector, "eks-connector.yaml")
			assertManifestFile(template.ClusterRole, "eks-connector-clusterrole.yaml")
			assertManifestFile(template.ConsoleAccess, "eks-connector-console-dashboard-full-access-group.yaml")
		})
	})
})
