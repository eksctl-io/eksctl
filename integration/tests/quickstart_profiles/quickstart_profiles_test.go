// +build integration

package integration_test

import (
	"testing"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/git"
	"github.com/weaveworks/eksctl/integration/utilities/unowned"
	"github.com/weaveworks/eksctl/pkg/testutils"

	"github.com/kubicorn/kubicorn/pkg/namer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var (
	params         *tests.Params
	unownedCluster *unowned.Cluster
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("qstart")
}

func TestQuickstartProfiles(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	if !params.SkipCreate {
		if params.UnownedCluster {
			cfg := api.NewClusterConfig()
			cfg.Metadata = &api.ClusterMeta{
				Version: params.Version,
				Name:    params.ClusterName,
				Region:  params.Region,
			}
			unownedCluster = unowned.NewCluster(cfg)

			unownedCluster.CreateNodegroups("ng-1")

			cmd := params.EksctlUtilsCmd.WithArgs(
				"write-kubeconfig",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		} else {
			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--name", params.ClusterName,
				"--verbose", "4",
				"--region", params.Region,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		}
	}
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
	if params.UnownedCluster {
		unownedCluster.DeleteStack()
	}
})

var _ = Describe("Enable and use GitOps quickstart profiles", func() {
	var (
		branch   string
		cloneDir string
		err      error
	)

	BeforeEach(func() {
		if branch == "" {
			branch = namer.RandomName()
			cloneDir, err = git.CreateBranch(branch)
		}
	})

	Context("enable repo", func() {
		It("should add Flux to the repo and the cluster", func() {
			Expect(err).NotTo(HaveOccurred()) // Creating the branch should have succeeded.
			AssertFluxManifestsAbsentInGit(branch)
			AssertFluxPodsAbsentInKubernetes(params.KubeconfigPath, "flux")

			cmd := params.EksctlCmd.WithArgs(
				"enable", "repo",
				"--git-url", git.Repository,
				"--git-email", git.Email,
				"--git-branch", branch,
				"--cluster", params.ClusterName,
			)
			Expect(cmd).To(RunSuccessfully())

			AssertFluxManifestsPresentInGit(branch)
			AssertFluxPodsPresentInKubernetes(params.KubeconfigPath)
		})
	})

	Context("enable repo", func() {
		It("should not add Flux to the repo and the cluster if there is a flux deployment already", func() {
			Expect(err).NotTo(HaveOccurred()) // Creating the branch should have succeeded.
			AssertFluxPodsPresentInKubernetes(params.KubeconfigPath)

			cmd := params.EksctlCmd.WithArgs(
				"enable", "repo",
				"--git-url", git.Repository,
				"--git-email", git.Email,
				"--git-branch", branch,
				"--cluster", params.ClusterName,
			)
			Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("Skipping installation")))
		})
	})

	Context("enable profile", func() {
		It("should add the configured quickstart profile to the repo and the cluster", func() {
			Expect(err).NotTo(HaveOccurred()) // Creating the branch should have succeeded.
			// Flux should have been installed by the previously run "enable repo" command:
			AssertFluxManifestsPresentInGit(branch)
			AssertFluxPodsPresentInKubernetes(params.KubeconfigPath)

			cmd := params.EksctlCmd.WithArgs(
				"enable", "profile",
				"--git-url", git.Repository,
				"--git-email", git.Email,
				"--git-branch", branch,
				"--cluster", params.ClusterName,
				"app-dev",
			)
			Expect(cmd).To(RunSuccessfully())

			AssertQuickStartComponentsPresentInGit(branch)
			// Flux should still be present:
			AssertFluxManifestsPresentInGit(branch)
			AssertFluxPodsPresentInKubernetes(params.KubeconfigPath)
			// Clean-up:
			err := git.CleanupBranchAndRepo(branch, cloneDir)
			Expect(err).NotTo(HaveOccurred()) // Deleting the branch should have succeeded.
		})
	})
})
