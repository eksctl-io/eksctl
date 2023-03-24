//go:build integration

//revive:disable Not changing package name
package override_bootstrap

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"

	. "github.com/onsi/ginkgo/v2"
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

var (
	customAMIAL2          string
	customAMIBottlerocket string
)

var _ = BeforeSuite(func() {
	cfg := NewConfig(params.Region)
	ssm := awsssm.NewFromConfig(cfg)

	// retrieve AL2 AMI
	input := &awsssm.GetParameterInput{
		Name: aws.String("/aws/service/eks/optimized-ami/1.22/amazon-linux-2/recommended/image_id"),
	}
	output, err := ssm.GetParameter(context.Background(), input)
	Expect(err).NotTo(HaveOccurred())
	customAMIAL2 = *output.Parameter.Value

	// retrieve Bottlerocket AMI
	input = &awsssm.GetParameterInput{
		Name: aws.String("/aws/service/bottlerocket/aws-k8s-1.25/x86_64/latest/image_id"),
	}
	output, err = ssm.GetParameter(context.Background(), input)
	Expect(err).NotTo(HaveOccurred())
	customAMIBottlerocket = *output.Parameter.Value

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

var _ = Describe("(Integration) [Test Custom AMI]", func() {
	params.LogStacksEventsOnFailure()

	Context("override bootstrap command for managed and un-managed nodegroups", func() {

		It("can create a working nodegroup which can join the cluster", func() {
			By(fmt.Sprintf("using the following EKS optimised AMI: %s", customAMIAL2))
			content, err := os.ReadFile(filepath.Join("testdata/override-bootstrap.yaml"))
			Expect(err).NotTo(HaveOccurred())
			content = bytes.ReplaceAll(content, []byte("<generated>"), []byte(params.ClusterName))
			content = bytes.ReplaceAll(content, []byte("<generated-region>"), []byte(params.Region))
			content = bytes.ReplaceAll(content, []byte("<generated-ami>"), []byte(customAMIAL2))
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

	Context("bottlerocket un-managed nodegroups", func() {

		It("can create a working nodegroup which can join the cluster", func() {
			By(fmt.Sprintf("using the following EKS optimised AMI: %s", customAMIBottlerocket))
			content, err := os.ReadFile(filepath.Join("testdata/bottlerocket-settings.yaml"))
			Expect(err).NotTo(HaveOccurred())
			content = bytes.ReplaceAll(content, []byte("<generated>"), []byte(params.ClusterName))
			content = bytes.ReplaceAll(content, []byte("<generated-region>"), []byte(params.Region))
			content = bytes.ReplaceAll(content, []byte("<generated-ami>"), []byte(customAMIBottlerocket))
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

var _ = AfterSuite(func() {
	params.DeleteClusters()
	gexec.KillAndWait()
	if params.KubeconfigTemp {
		Expect(os.Remove(params.KubeconfigPath)).To(Succeed())
	}
	Expect(os.RemoveAll(params.TestDirectory)).To(Succeed())
})
