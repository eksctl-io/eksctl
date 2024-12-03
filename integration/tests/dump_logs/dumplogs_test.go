//go:build integration
// +build integration

package dumplogs

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/integration/tests"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	testing.Init()
	if err := api.Register(); err != nil {
		panic(fmt.Errorf("unexpected error registering API scheme: %w", err))
	}
	params = tests.NewParams("dump_logs")
}

func TestDumpLogs(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Dump logs on failure", func() {

	os.RemoveAll("logs/")

	cmd := params.EksctlGetCmd.WithArgs(
		"cluster",
		"-d",
		"--name",
		"definitely_not_a_cluster",
	)

	session := cmd.Run()

	It("exits with exit code 1", func() {
		Expect(session.ExitCode()).To(Equal(1))
	})

	It("logs/ directory is created", func() {
		Expect("logs/").Should(BeADirectory())
	})

	It("logs/ directory should have one file", func() {
		d, err := os.ReadDir("logs/")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(d)).To(Equal(1))
		os.RemoveAll("logs/")
	})

	// Re-run the command without the -d opt-in flag
	cmd = params.EksctlGetCmd.WithArgs(
		"cluster",
		"--name",
		"definitely_not_a_cluster",
	)

	session = cmd.Run()

	It("exits with exit code 1", func() {
		Expect(session.ExitCode()).To(Equal(1))
	})

	It("logs are not dumped", func() {
		Expect("logs/").ShouldNot(BeADirectory())
	})
})
