// +build integration

package backwards_compat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("bwardscomp")
}

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const (
	goBackVersions = 2
)

var _ = Describe("(Integration) [Backwards compatibility test]", func() {

	var (
		initialNgName = "ng-1"
		newNgName     = "ng-2"
	)

	It("should support clusters created with a previous version of eksctl", func() {
		By("downloading a previous release")
		eksctlDir, err := ioutil.TempDir(os.TempDir(), "eksctl")
		Expect(err).ToNot(HaveOccurred())

		defer func() {
			Expect(os.RemoveAll(eksctlDir)).ToNot(HaveOccurred())
		}()

		downloadRelease(eksctlDir)

		eksctlPath := path.Join(eksctlDir, "eksctl")

		version, err := getVersion(eksctlPath)
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("creating a cluster with release %q", version))
		cmd := runner.NewCmd(eksctlPath).
			WithArgs(
				"create",
				"cluster",
				"--name", params.ClusterName,
				"--nodegroup-name", initialNgName,
				"-v4",
				"--region", params.Region,
			).
			WithTimeout(20 * time.Minute)

		Expect(cmd).To(RunSuccessfully())

		By("fetching the new cluster")
		cmd = params.EksctlGetCmd.WithArgs(
			"cluster",
			params.ClusterName,
			"--output", "json",
		)

		Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(params.ClusterName)))

		By("adding a nodegroup")
		cmd = params.EksctlCreateCmd.WithArgs(
			"nodegroup",
			"--cluster", params.ClusterName,
			"--nodes", "2",
			newNgName,
		)
		Expect(cmd).To(RunSuccessfully())

		By("scaling the initial nodegroup")
		cmd = params.EksctlScaleNodeGroupCmd.WithArgs(
			"--cluster", params.ClusterName,
			"--nodes", "3",
			"--name", initialNgName,
		)
		Expect(cmd).To(RunSuccessfully())

		By("deleting the new nodegroup")
		cmd = params.EksctlDeleteCmd.WithArgs(
			"nodegroup",
			"--verbose", "4",
			"--cluster", params.ClusterName,
			newNgName,
		)
		Expect(cmd).To(RunSuccessfully())

		By("deleting the initial nodegroup")
		cmd = params.EksctlDeleteCmd.WithArgs(
			"nodegroup",
			"--verbose", "4",
			"--cluster", params.ClusterName,
			initialNgName,
		)
		Expect(cmd).To(RunSuccessfully())
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})

func downloadRelease(dir string) {
	cmd := runner.NewCmd("../../scripts/download-previous-release.sh").
		WithEnv(
			fmt.Sprintf("GO_BACK_VERSIONS=%d", goBackVersions),
			fmt.Sprintf("DOWNLOAD_DIR=%s", dir),
		).
		WithTimeout(30 * time.Second)

	ExpectWithOffset(1, cmd).To(RunSuccessfully())
}

func getVersion(eksctlPath string) (string, error) {
	cmd := runner.NewCmd(eksctlPath).WithArgs("version")
	session := cmd.Run()
	if session.ExitCode() != 0 {
		return "", errors.New(string(session.Err.Contents()))
	}
	return string(session.Buffer().Contents()), nil
}
