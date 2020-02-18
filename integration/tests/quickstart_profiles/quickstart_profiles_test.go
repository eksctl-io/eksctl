// +build integration

package integration_test

import (
	"testing"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/git"
	"github.com/weaveworks/eksctl/pkg/testutils"

	"github.com/kubicorn/kubicorn/pkg/namer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("qstart")
}

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	if params.DoCreate {
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

var _ = Describe("Enable and use GitOps quickstart profiles", func() {
	var (
		branch   string
		cloneDir string
		err      error
	)

	BeforeEach(func() {
		if branch == "" {
			branch = namer.RandomName()
			cloneDir, err = git.CreateBranch(branch, params.PrivateSSHKeyPath)
		}
	})

	Context("enable repo", func() {
		It("should add Flux to the repo and the cluster", func() {
			Expect(err).NotTo(HaveOccurred()) // Creating the branch should have succeeded.
			AssertFluxManifestsAbsentInGit(branch, params.PrivateSSHKeyPath)
			AssertFluxPodsAbsentInKubernetes(params.KubeconfigPath)

			cmd := params.EksctlExperimentalCmd.WithArgs(
				"enable", "repo",
				"--git-url", git.Repository,
				"--git-email", git.Email,
				"--git-private-ssh-key-path", params.PrivateSSHKeyPath,
				"--git-branch", branch,
				"--cluster", params.ClusterName,
			)
			Expect(cmd).To(RunSuccessfully())

			AssertFluxManifestsPresentInGit(branch, params.PrivateSSHKeyPath)
			AssertFluxPodsPresentInKubernetes(params.KubeconfigPath)
		})
	})

	Context("enable profile", func() {
		It("should add the configured quickstart profile to the repo and the cluster", func() {
			Expect(err).NotTo(HaveOccurred()) // Creating the branch should have succeeded.
			// Flux should have been installed by the previously run "enable repo" command:
			AssertFluxManifestsPresentInGit(branch, params.PrivateSSHKeyPath)
			AssertFluxPodsPresentInKubernetes(params.KubeconfigPath)

			cmd := params.EksctlExperimentalCmd.WithArgs(
				"enable", "profile",
				"--git-url", git.Repository,
				"--git-email", git.Email,
				"--git-branch", branch,
				"--git-private-ssh-key-path", params.PrivateSSHKeyPath,
				"--cluster", params.ClusterName,
				"app-dev",
			)
			Expect(cmd).To(RunSuccessfully())

			AssertQuickStartComponentsPresentInGit(branch, params.PrivateSSHKeyPath)
			// Flux should still be present:
			AssertFluxManifestsPresentInGit(branch, params.PrivateSSHKeyPath)
			AssertFluxPodsPresentInKubernetes(params.KubeconfigPath)
			// Clean-up:
			err := git.DeleteBranch(branch, cloneDir, params.PrivateSSHKeyPath)
			Expect(err).NotTo(HaveOccurred()) // Deleting the branch should have succeeded.
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
