// +build integration

package integration_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("(Integration) generate profile", func() {

	Describe("when generating a profile", func() {
		It("should write the processed repo files in the supplied directory", func() {

			clusterName = "amazing-testing-gopher"

			eksctlSuccess("generate", "profile",
				"--verbose", "4",
				"--name", clusterName,
				"--region", region,
				"--git-url", "git@github.com:eksctl-bot/eksctl-profile-integration-tests.git",
				"--profile-path", testDirectory,
			)

			fs := afero.Afero{
				Fs: afero.NewOsFs(),
			}

			contents, err := fs.ReadFile(filepath.Join(testDirectory, "workloads/namespace.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(MatchYAML(
				`---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    name: amazing-testing-gopher
  name: amazing-testing-gopher
`))

			contents, err = fs.ReadFile(filepath.Join(testDirectory, "workloads/services/service.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(MatchYAML(
				`---
apiVersion: v1
kind: Service
metadata:
  name: amazing-testing-gopher-service1
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
`))

			contents, err = fs.ReadFile(filepath.Join(testDirectory, "metadata.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(MatchYAML(
				`---
somekey:
  repo: eks-gitops-tests
  thisFile: should not be modified by eksctl generate profile
anotherkey:
  nestedKey: nestedvalue
`))

		})
	})
})
