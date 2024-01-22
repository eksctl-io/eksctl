//go:build integration
// +build integration

package unowned

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("exist-vpc")
}

func TestVPC(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	stackName string
	ng1       = "ng-1"
	mng1      = "mng-1"
	ctl       api.ClusterProvider
	cfg       *api.ClusterConfig
)

var _ = BeforeSuite(func() {
	stackName = fmt.Sprintf("eksctl-%s", params.ClusterName)
	cfg = &api.ClusterConfig{
		TypeMeta: api.ClusterConfigTypeMeta(),
		Metadata: &api.ClusterMeta{
			Name:   params.ClusterName,
			Region: params.Region,
		},
	}

	clusterProvider, err := eks.New(context.Background(), &api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())
	ctl = clusterProvider.AWSProvider
	cfg.VPC = createVPC(stackName, ctl)

	cmd := params.EksctlCreateCmd.
		WithArgs(
			"cluster",
			"--config-file", "-",
			"--verbose", "2",
		).
		WithoutArg("--region", params.Region).
		WithStdin(clusterutils.Reader(cfg))

	Expect(cmd).To(RunSuccessfully())
})

var _ = Describe("(Integration) [using existing VPC]", func() {
	params.LogStacksEventsOnFailure()

	It("supports creating managed and unmanaged nodegroups in the existing VPC", func() {
		cfg.NodeGroups = []*api.NodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name: ng1,
			}},
		}
		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name: mng1,
			}},
		}

		By("creating the nodegroups")
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"nodegroup",
				"--config-file", "-",
				"--verbose", "2",
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(cfg))
		Expect(cmd).To(RunSuccessfully())

		By("checking the cluster is created with the corret VPC/subnets")
		cmd = params.EksctlGetCmd.WithArgs("cluster", "--name", params.ClusterName, "-o", "yaml")
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(cfg.VPC.ID)),
			ContainElement(ContainSubstring(cfg.VPC.Subnets.Public["public1"].ID)),
			ContainElement(ContainSubstring(cfg.VPC.Subnets.Public["public2"].ID)),
			ContainElement(ContainSubstring(cfg.VPC.Subnets.Public["public3"].ID)),
			ContainElement(ContainSubstring(cfg.VPC.Subnets.Public["private4"].ID)),
			ContainElement(ContainSubstring(cfg.VPC.Subnets.Public["private5"].ID)),
			ContainElement(ContainSubstring(cfg.VPC.Subnets.Public["private6"].ID)),
		))
	})
})

func createVPC(stackName string, ctl api.ClusterProvider) *api.ClusterVPC {
	publicSubnets, privateSubnets, vpcID, securityGroup := createVPCStackAndGetOutputs(stackName, ctl)

	By("creating the cluster config from the existing VPC")
	newVPC := api.NewClusterVPC(false)
	newVPC.ID = vpcID
	newVPC.SecurityGroup = securityGroup

	output, err := ctl.EC2().DescribeSubnets(context.Background(), &ec2.DescribeSubnetsInput{
		SubnetIds: append(publicSubnets, privateSubnets...),
	})
	Expect(err).NotTo(HaveOccurred())
	subnetToAZMap := map[string]string{}
	for _, s := range output.Subnets {
		subnetToAZMap[*s.SubnetId] = *s.AvailabilityZone
	}

	newVPC.Subnets = &api.ClusterSubnets{
		Public: api.AZSubnetMapping{
			"public1": api.AZSubnetSpec{
				ID: publicSubnets[0],
				AZ: subnetToAZMap[publicSubnets[0]],
			},
			"public2": api.AZSubnetSpec{
				ID: publicSubnets[1],
				AZ: subnetToAZMap[publicSubnets[1]],
			},
			"public3": api.AZSubnetSpec{
				ID: publicSubnets[2],
				AZ: subnetToAZMap[publicSubnets[2]],
			},
		},
		Private: api.AZSubnetMapping{
			"private4": api.AZSubnetSpec{
				ID: privateSubnets[0],
				AZ: subnetToAZMap[privateSubnets[0]],
			},
			"private5": api.AZSubnetSpec{
				ID: privateSubnets[1],
				AZ: subnetToAZMap[privateSubnets[1]],
			},
			"private6": api.AZSubnetSpec{
				ID: privateSubnets[2],
				AZ: subnetToAZMap[privateSubnets[2]],
			},
		},
	}

	return newVPC
}

func createVPCStackAndGetOutputs(stackName string, ctl api.ClusterProvider) ([]string, []string, string, string) {
	templateBody, err := os.ReadFile("cf-template.yaml")
	Expect(err).NotTo(HaveOccurred())
	createStackInput := &cfn.CreateStackInput{
		StackName: &stackName,
	}
	By("creating the stack")
	createStackInput.TemplateBody = aws.String(string(templateBody))
	createStackInput.Capabilities = []types.Capability{types.CapabilityCapabilityIam, types.CapabilityCapabilityNamedIam}

	_, err = ctl.CloudFormation().CreateStack(context.Background(), createStackInput)
	Expect(err).NotTo(HaveOccurred())

	var describeStackOut *cfn.DescribeStacksOutput
	Eventually(func() types.StackStatus {
		describeStackOut, err = ctl.CloudFormation().DescribeStacks(context.TODO(), &cfn.DescribeStacksInput{
			StackName: &stackName,
		})
		Expect(err).NotTo(HaveOccurred())
		return describeStackOut.Stacks[0].StackStatus
	}, time.Minute*10, time.Second*15).Should(Equal(types.StackStatusCreateComplete))

	By("fetching the outputs of the stack")
	var vpcID string
	var publicSubnets, privateSubnets, securityGroups []string
	for _, output := range describeStackOut.Stacks[0].Outputs {
		switch *output.OutputKey {
		case "PublicSubnetIds":
			publicSubnets = strings.Split(*output.OutputValue, ",")
		case "PrivateSubnetIds":
			privateSubnets = strings.Split(*output.OutputValue, ",")
		case "VpcId":
			vpcID = *output.OutputValue
		case "SecurityGroups":
			securityGroups = strings.Split(*output.OutputValue, ",")
		}
	}

	return publicSubnets, privateSubnets, vpcID, securityGroups[0]
}

func deleteStack(stackName string, ctl api.ClusterProvider) {
	deleteStackInput := &cfn.DeleteStackInput{
		StackName: &stackName,
	}

	_, err := ctl.CloudFormation().DeleteStack(context.Background(), deleteStackInput)
	Expect(err).NotTo(HaveOccurred())
}

var _ = AfterSuite(func() {
	cmd := params.EksctlDeleteClusterCmd.
		WithArgs(
			"--disable-nodegroup-eviction",
			"--config-file", "-",
			"--wait",
		).
		WithoutArg("--region", params.Region).
		WithStdin(clusterutils.Reader(cfg))
	Expect(cmd).To(RunSuccessfully())
	deleteStack(stackName, ctl)
})
