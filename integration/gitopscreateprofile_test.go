// +build integration

package integration_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

)

var _ = Describe("(Integration) generate profile", func() {

	Describe("when generating a profile", func() {
		It("should write the processed repo files in the supplied directory", func() {

			if clusterName == "" {
				clusterName = cmdutils.ClusterName("", "")
			}

			cmd := eksctlExperimentalCmd.WithArgs(
				"generate", "profile",
				"--verbose", "4",
				"--cluster", clusterName,
				"--git-url", "git@github.com:eksctl-bot/eksctl-profile-integration-tests.git",
				"--profile-path", testDirectory,
			)
			Expect(cmd).To(RunSuccessfully())

			fs := afero.Afero{
				Fs: afero.NewOsFs(),
			}

			contents, err := fs.ReadFile(filepath.Join(testDirectory, "workloads/namespace.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(MatchYAML(fmt.Sprintf(
				`---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    name: %s-%s
  name: %s
`, clusterName, region, clusterName)))

			contents, err = fs.ReadFile(filepath.Join(testDirectory, "workloads/services/service.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(MatchYAML(fmt.Sprintf(
				`---
apiVersion: v1
kind: Service
metadata:
  name: %s-service1
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
`, clusterName)))

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
