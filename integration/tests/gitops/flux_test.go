//go:build integration
// +build integration

package integration_test

import (
	"fmt"
	"testing"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/git"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	"github.com/kubicorn/kubicorn/pkg/namer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	params *tests.Params
)

const (
	repository = "my-gitops-repo"
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("flux")
}

func TestQuickstartProfiles(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	if !params.SkipCreate {
		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--name", params.ClusterName,
			"--verbose", "4",
			"--region", params.Region,
			"--kubeconfig", params.KubeconfigPath,
		)
		Expect(cmd).To(RunSuccessfully())
	}
})

var _ = Describe("Enable GitOps", func() {
	var (
		branch        string
		clusterConfig *api.ClusterConfig
		localRepo     string
	)

	BeforeEach(func() {
		if branch == "" {
			branch = namer.RandomName()
		}

		clusterConfig = &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Version: api.DefaultVersion,
				Name:    params.ClusterName,
				Region:  params.Region,
			},
			GitOps: &api.GitOps{
				Flux: &api.Flux{
					GitProvider: "github",
					Flags: api.FluxFlags{
						"owner":      params.GitopsOwner,
						"branch":     branch,
						"repository": repository,
					},
				},
			},
		}
	})

	AfterEach(func() {
		if err := git.CleanupBranchAndRepo(branch, localRepo); err != nil {
			fmt.Fprintf(GinkgoWriter, "error cleaning up branch and repo: %v", err)
		}
	})

	Context("enable flux", func() {
		It("should deploy Flux v2 components to the cluster", func() {
			AssertFluxPodsAbsentInKubernetes(params.KubeconfigPath, "flux-system")
			cmd := params.EksctlEnableCmd.WithArgs(
				"flux",
				"--config-file", "-",
			).WithStdin(clusterutils.Reader(clusterConfig))

			Expect(cmd).To(RunSuccessfully())
			AssertFlux2PodsPresentInKubernetes(params.KubeconfigPath)
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
