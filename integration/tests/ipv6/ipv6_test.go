//go:build integration
// +build integration

package ipv6

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("IPv6")
}

func TestIPv6(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [EKS IPv6 test]", func() {

	Context("Creating a cluster with IPv6", func() {
		clusterName := params.NewClusterName("ipv6")

		BeforeSuite(func() {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = clusterName
			clusterConfig.Metadata.Version = "latest"
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.VPC.IPFamily = aws.String("IPv6")
			clusterConfig.IAM.WithOIDC = api.Enabled()
			clusterConfig.Addons = []*api.Addon{
				{
					Name: "vpc-cni",
				},
				{
					Name: "kube-proxy",
				},
				{
					Name: "coredns",
				},
			}

			data, err := json.Marshal(clusterConfig)
			Expect(err).ToNot(HaveOccurred())

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())
		})

		AfterSuite(func() {
			cmd := params.EksctlDeleteCmd.WithArgs(
				"cluster", clusterName,
				"--verbose", "2",
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should support ipv6", func() {
			By("Asserting that the VPC that is created has a IPv6 CIDR")
			awsSession := NewSession(params.Region)
			cfnSession := cfn.New(awsSession)

			var describeStackOut *cfn.DescribeStacksOutput
			describeStackOut, err := cfnSession.DescribeStacks(&cfn.DescribeStacksInput{
				StackName: aws.String(fmt.Sprintf("eksctl-%s-cluster", clusterName)),
			})
			Expect(err).NotTo(HaveOccurred())

			var vpcID string
			for _, output := range describeStackOut.Stacks[0].Outputs {
				if *output.OutputKey == "VPC" {
					vpcID = *output.OutputValue
				}
			}

			ec2 := awsec2.New(awsSession)
			output, err := ec2.DescribeVpcs(&awsec2.DescribeVpcsInput{
				VpcIds: aws.StringSlice([]string{vpcID}),
			})
			Expect(err).NotTo(HaveOccurred(), output.GoString())
			Expect(output.Vpcs[0].Ipv6CidrBlockAssociationSet).To(HaveLen(1))
		})
	})
})
