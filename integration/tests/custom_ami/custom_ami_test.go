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
	customAMIAL2           string
	customAMIAL2023        string
	customAMIBottlerocket  string
	customAMIUbuntuPro2204 string
)

var _ = BeforeSuite(func() {
	cfg := NewConfig(params.Region)
	ssm := awsssm.NewFromConfig(cfg)

	// retrieve AL2 AMI
	input := &awsssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id", params.Version)),
	}
	output, err := ssm.GetParameter(context.Background(), input)
	Expect(err).NotTo(HaveOccurred())
	customAMIAL2 = *output.Parameter.Value

	// retrieve AL2023 AMI
	input = &awsssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2023/x86_64/standard/recommended/image_id", params.Version)),
	}
	output, err = ssm.GetParameter(context.Background(), input)
	Expect(err).NotTo(HaveOccurred())
	customAMIAL2023 = *output.Parameter.Value

	// retrieve Bottlerocket AMI
	input = &awsssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/x86_64/latest/image_id", params.Version)),
	}
	output, err = ssm.GetParameter(context.Background(), input)
	Expect(err).NotTo(HaveOccurred())
	customAMIBottlerocket = *output.Parameter.Value

	// retrieve Ubuntu Pro 22.04 AMI
	input = &awsssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/canonical/ubuntu/eks-pro/22.04/%s/stable/current/amd64/hvm/ebs-gp2/ami-id", params.Version)),
	}
	output, err = ssm.GetParameter(context.Background(), input)
	Expect(err).NotTo(HaveOccurred())
	customAMIUbuntuPro2204 = *output.Parameter.Value

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

	Context("al2023 managed and un-managed nodegroups", func() {
		It("can create working nodegroups which can join the cluster", func() {
			By(fmt.Sprintf("using the following EKS optimised AMI: %s", customAMIAL2023))
			content, err := os.ReadFile(filepath.Join("testdata/al2023.yaml"))
			Expect(err).NotTo(HaveOccurred())
			content = bytes.ReplaceAll(content, []byte("<generated>"), []byte(params.ClusterName))
			content = bytes.ReplaceAll(content, []byte("<generated-region>"), []byte(params.Region))
			content = bytes.ReplaceAll(content, []byte("<generated-ami>"), []byte(customAMIAL2023))
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

	Context("ubuntu-pro-2204 un-managed nodegroups", func() {

		It("can create a working nodegroup which can join the cluster", func() {
			By(fmt.Sprintf("using the following EKS optimised AMI: %s", customAMIUbuntuPro2204))
			content, err := os.ReadFile(filepath.Join("testdata/ubuntu-pro-2204.yaml"))
			Expect(err).NotTo(HaveOccurred())
			content = bytes.ReplaceAll(content, []byte("<generated>"), []byte(params.ClusterName))
			content = bytes.ReplaceAll(content, []byte("<generated-region>"), []byte(params.Region))
			content = bytes.ReplaceAll(content, []byte("<generated-ami>"), []byte(customAMIUbuntuPro2204))
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
