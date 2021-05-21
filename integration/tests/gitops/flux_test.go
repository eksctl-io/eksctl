// +build integration

package integration_test

import (
	"encoding/json"
	"io/ioutil"
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

			var err error
			localRepo, err = git.CreateBranch(branch)
			Expect(err).NotTo(HaveOccurred())
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
					Repository:  repository,
					Branch:      branch,
					GitProvider: "github",
					Owner:       params.GitopsOwner,
					Kubeconfig:  params.KubeconfigPath,
				},
			},
		}
		configData, err := json.Marshal(&cfg)
		Expect(err).NotTo(HaveOccurred())
		configFile, err = ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(configFile.Name(), configData, 0755)).To(Succeed())
	})

	AfterEach(func() {
		_ = git.CleanupBranchAndRepo(branch, localRepo)
		Expect(os.RemoveAll(configFile.Name())).To(Succeed())
	})

	Context("enable flux", func() {
		BeforeEach(func() {
			cmd = params.EksctlEnableCmd.WithArgs(
				"flux",
				"--config-file", configFile.Name(),
			)
		})

		It("should deploy Flux v2 components to the cluster", func() {
			AssertFluxPodsAbsentInKubernetes(params.KubeconfigPath, "flux-system")
			Expect(cmd).To(RunSuccessfully())
			AssertFlux2PodsPresentInKubernetes(params.KubeconfigPath)
		})

		It("should not add Flux to the cluster if there is a flux deployment already", func() {
			Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("skipping installation")))
		})
	})
})
