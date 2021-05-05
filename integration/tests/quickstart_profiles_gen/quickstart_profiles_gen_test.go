// +build integration

package quickstart_profiles_gen

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/unowned"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var (
	params         *tests.Params
	unownedCluster *unowned.Cluster
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("qstartgen")
}

func TestQuickstartProfilesGen(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	if params.UnownedCluster {
		unownedCluster = unowned.NewCluster(&api.ClusterConfig{
			Metadata: &api.ClusterMeta{
				Name:    params.ClusterName,
				Region:  params.Region,
				Version: params.Version,
			},
		})
	} else {
		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", params.ClusterName,
			"--region", params.Region,
		)
		Expect(cmd).To(RunSuccessfully())
	}
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
	os.RemoveAll(params.TestDirectory)
	if params.UnownedCluster {
		unownedCluster.DeleteStack()
	}
})

var _ = Describe("(Integration) generate profile", func() {

	Describe("when generating a profile", func() {
		It("should write the processed repo files in the supplied directory", func() {
			cmd := params.EksctlCmd.WithArgs(
				"generate", "profile",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--profile-source", "git@github.com:eksctl-bot/eksctl-profile-integration-tests.git",
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
