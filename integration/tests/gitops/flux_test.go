//go:build integration
// +build integration

package integration_test

import (
	"encoding/json"
	"os"
	"testing"

	. "github.com/weaveworks/eksctl/integration/matchers"
	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/git"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	"github.com/kubicorn/kubicorn/pkg/namer"
	. "github.com/onsi/ginkgo"
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

var _ = AfterSuite(func() {
	params.DeleteClusters()
})

var _ = Describe("Enable GitOps", func() {
	var (
		branch     string
		cmd        runner.Cmd
		configFile *os.File
		localRepo  string
	)

	BeforeEach(func() {
		if branch == "" {
			branch = namer.RandomName()
		}

		cfg := &api.ClusterConfig{
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
		configData, err := json.Marshal(&cfg)
		Expect(err).NotTo(HaveOccurred())
		configFile, err = os.CreateTemp("", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile.Name(), configData, 0755)).To(Succeed())
	})

	AfterEach(func() {
		_ = git.CleanupBranchAndRepo(branch, localRepo)
		Expect(os.RemoveAll(configFile.Name())).To(Succeed())
	})

	Context("enable flux", func() {
		It("should deploy Flux v2 components to the cluster", func() {
			AssertFluxPodsAbsentInKubernetes(params.KubeconfigPath, "flux-system")
			cmd = params.EksctlEnableCmd.WithArgs(
				"flux",
				"--config-file", configFile.Name(),
			)
			Expect(cmd).To(RunSuccessfully())
			AssertFlux2PodsPresentInKubernetes(params.KubeconfigPath)
		})
	})
})
