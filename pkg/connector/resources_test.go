package connector_test

import (
	"fmt"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/weaveworks/eksctl/pkg/connector"
)

var _ = Describe("Writing manifests", func() {
	Context("WriteResources", func() {
		It("should write the manifests for EKS Connector", func() {
			fs := afero.NewMemMapFs()
			manifestList := &connector.ManifestList{
				ConnectorResources: connector.ManifestFile{
					Data:     []byte("connector"),
					Filename: "eks-connector.yaml",
				},
				ClusterRoleResources: connector.ManifestFile{
					Data:     []byte("clusterrole"),
					Filename: "eks-connector-clusterrole.yaml",
				},
				ConsoleAccessResources: connector.ManifestFile{
					Data:     []byte("console-dashboard-full-access-group"),
					Filename: "eks-connector-console-dashboard-full-access-group.yaml",
				},
			}
			err := connector.WriteResources(fs, manifestList)
			Expect(err).ToNot(HaveOccurred())

			wd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			files, err := afero.ReadDir(fs, wd)
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(3))

			filenameData := map[string][]byte{}
			for _, manifestFile := range []connector.ManifestFile{manifestList.ConnectorResources, manifestList.ClusterRoleResources, manifestList.ConsoleAccessResources} {
				filenameData[manifestFile.Filename] = manifestFile.Data
			}

			for _, file := range files {
				data, ok := filenameData[file.Name()]
				if !ok {
					Fail(fmt.Sprintf("unexpected filename %q", file.Name()))
				}
				file, err := afero.ReadFile(fs, path.Join(wd, file.Name()))
				Expect(err).ToNot(HaveOccurred())
				Expect(file).To(Equal(data))
			}
		})
	})
})
