// +build integration

package quickstart_profiles_gen

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("qstartgen")
}

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--verbose", "4",
		"--name", params.ClusterName,
		"--region", params.Region,
	)
	Expect(cmd).To(RunSuccessfully())
})

var _ = Describe("(Integration) generate profile", func() {

	Describe("when generating a profile", func() {
		It("should write the processed repo files in the supplied directory", func() {
			cmd := params.EksctlExperimentalCmd.WithArgs(
				"generate", "profile",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--git-url", "git@github.com:eksctl-bot/eksctl-profile-integration-tests.git",
				"--profile-path", params.TestDirectory,
			)
			Expect(cmd).To(RunSuccessfully())

			fs := afero.Afero{
				Fs: afero.NewOsFs(),
			}

			contents, err := fs.ReadFile(filepath.Join(params.TestDirectory, "workloads/namespace.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(MatchYAML(fmt.Sprintf(
				`---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    name: %s-%s
  name: %s
`, params.ClusterName, params.Region, params.ClusterName)))

			contents, err = fs.ReadFile(filepath.Join(params.TestDirectory, "workloads/services/service.yaml"))
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
`, params.ClusterName)))

			contents, err = fs.ReadFile(filepath.Join(params.TestDirectory, "metadata.yaml"))
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

var _ = AfterSuite(func() {
	params.DeleteClusters()
	os.RemoveAll(params.TestDirectory)
})
