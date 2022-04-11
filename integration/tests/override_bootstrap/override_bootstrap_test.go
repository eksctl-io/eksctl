//go:build integration

//revive:disable Not changing package name
package override_bootstrap

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("override")
}

func TestOverrideBootstrap(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [Test OverrideBootstrapCommand]", func() {
	var (
		customAMI string
	)

	BeforeSuite(func() {
		awsSession := NewSession(params.Region)
		ssm := awsssm.New(awsSession)
		input := &awsssm.GetParameterInput{
			Name: aws.String("/aws/service/eks/optimized-ami/1.21/amazon-linux-2/recommended/image_id"),
		}
		output, err := ssm.GetParameter(input)
		Expect(err).NotTo(HaveOccurred())
		customAMI = *output.Parameter.Value
		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", params.ClusterName,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--version", params.Version,
			"--kubeconfig", params.KubeconfigPath,
		)
		Expect(cmd).To(RunSuccessfully())
	})

	AfterSuite(func() {
		params.DeleteClusters()
		gexec.KillAndWait()
		if params.KubeconfigTemp {
			Expect(os.Remove(params.KubeconfigPath)).To(Succeed())
		}
		Expect(os.RemoveAll(params.TestDirectory)).To(Succeed())
	})

	Context("override bootstrap command for managed and un-managed nodegroups", func() {

		It("can create a working nodegroups which can join the cluster", func() {
			By(fmt.Sprintf("using the following EKS optimised AMI: %s", customAMI))
			content, err := os.ReadFile(filepath.Join("testdata/override-bootstrap.yaml"))
			Expect(err).NotTo(HaveOccurred())
			content = bytes.ReplaceAll(content, []byte("<generated>"), []byte(params.ClusterName))
			content = bytes.ReplaceAll(content, []byte("<generated-region>"), []byte(params.Region))
			content = bytes.ReplaceAll(content, []byte("<generated-ami>"), []byte(customAMI))
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(content))
			Expect(cmd).To(RunSuccessfully())
		})

	})
})
