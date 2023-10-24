//go:build integration
// +build integration

//revive:disable Not changing package name
package unowned_clusters

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	. "github.com/weaveworks/eksctl/integration/matchers"
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
	params = tests.NewParams("unowned")
}

func TestE2E(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	stackName, ng1, mng1, mng2 string
	version                    string
	upgradeVersion             string
	ctl                        api.ClusterProvider
	cfg                        *api.ClusterConfig
)

var _ = BeforeSuite(func() {
	ng1 = "ng-1"
	mng1 = "mng-1"
	mng2 = "mng-2"
	stackName = fmt.Sprintf("eksctl-%s", params.ClusterName)

	version, upgradeVersion = clusterutils.GetCurrentAndNextVersionsForUpgrade(params.Version)

	cfg = &api.ClusterConfig{
		TypeMeta: api.ClusterConfigTypeMeta(),
		Metadata: &api.ClusterMeta{
			Version: version,
			Name:    params.ClusterName,
			Region:  params.Region,
		},
	}

	if !params.SkipCreate {
		clusterProvider, err := eks.New(context.Background(), &api.ProviderConfig{Region: params.Region}, cfg)
		Expect(err).NotTo(HaveOccurred())
		ctl = clusterProvider.AWSProvider
		cfg.VPC = createClusterWithNodeGroup(context.Background(), params.ClusterName, stackName, mng1, version, ctl)
	}
})

var _ = Describe("(Integration) [non-eksctl cluster & nodegroup support]", func() {

	It("supports creating nodegroups", func() {
		cfg.NodeGroups = []*api.NodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name: ng1,
			}},
		}

		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name: mng2,
			}},
		}

		cmd := params.EksctlCreateNodegroupCmd.
			WithArgs(
				"--config-file", "-",
				"--verbose", "2",
			).
			WithStdin(clusterutils.Reader(cfg))
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports getting non-eksctl resources", func() {
		By("getting clusters")
		cmd := params.EksctlGetCmd.
			WithArgs(
				"clusters",
				"--verbose", "2",
			)
		AssertContainsCluster(cmd, GetClusterOutput{
			ClusterName:   params.ClusterName,
			Region:        params.Region,
			EksctlCreated: "False",
		})

		By("getting nodegroups")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"nodegroups",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(ng1)),
		))
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(mng1)),
		))
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(mng2)),
		))
	})

	It("supports labels", func() {
		By("setting labels on a managed nodegroup")
		cmd := params.EksctlSetLabelsCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--labels", "key=value",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("getting labels for a managed nodegroup")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"labels",
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--verbose", "2",
			)
		// It sometimes takes forever for the above set to take effect
		Eventually(func() *gbytes.Buffer { return cmd.Run().Out }, time.Minute*4).Should(gbytes.Say("key=value"))

		By("unsetting labels on a managed nodegroup")
		cmd = params.EksctlUnsetLabelsCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--labels", "key",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports IRSA", func() {
		By("enabling OIDC")
		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"associate-iam-oidc-provider",
				"--cluster", params.ClusterName,
				"--approve",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("creating an IAMServiceAccount")
		cmd = params.EksctlCreateCmd.
			WithArgs(
				"iamserviceaccount",
				"--cluster", params.ClusterName,
				"--name", "test-sa",
				"--namespace", "default",
				"--attach-policy-arn",
				"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
				"--approve",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("getting IAMServiceAccounts")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"iamserviceaccounts",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("test-sa")),
		))
	})

	It("supports cluster upgrades", func() {
		By("upgrading the cluster")
		cmd := params.EksctlUpgradeCmd.
			WithArgs(
				"cluster",
				"--name", params.ClusterName,
				"--version", upgradeVersion,
				"--timeout", "1h30m",
				"--approve",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports addons", func() {
		By("creating an addon")
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"addon",
				"--cluster", params.ClusterName,
				"--name", "vpc-cni",
				"--wait",
				"--force",
				"--version", "latest",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("getting an addon")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"addons",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("vpc-cni")),
		))
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("ACTIVE")),
		))
	})

	It("supports fargate", func() {
		By("creating a fargate profile")
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"fargateprofile",
				"--cluster", params.ClusterName,
				"--name", "fp-test",
				"--namespace", "default",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(SatisfyAll(ContainSubstring("created"), ContainSubstring("fp-test"))),
		))

		By("getting a fargate profile")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"fargateprofile",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("fp-test")),
		))

		By("deleting a fargate profile")
		cmd = params.EksctlDeleteCmd.
			WithArgs(
				"fargateprofile",
				"--cluster", params.ClusterName,
				"--name", "fp-test",
				"--wait",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(SatisfyAll(ContainSubstring("deleted"), ContainSubstring("fp-test"))),
		))
	})

	It("supports managed nodegroup upgrades", func() {
		cmd := params.EksctlUpgradeCmd.
			WithArgs(
				"nodegroup",
				"--name", mng1,
				"--cluster", params.ClusterName,
				"--kubernetes-version", upgradeVersion,
				"--timeout", "1h30m",
				"--wait",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports draining and scaling nodegroups", func() {
		By("scaling a nodegroup")
		cmd := params.EksctlScaleNodeGroupCmd.
			WithArgs(
				"--name", mng1,
				"--nodes", "2",
				"--nodes-max", "3",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("draining a nodegroup")
		cmd = params.EksctlDrainNodeGroupCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--name", mng1,
				"--parallel", "2",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports deleting nodegroups", func() {
		cmd := params.EksctlDeleteCmd.
			WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--name", mng1,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports deleting clusters", func() {
		if params.SkipDelete {
			Skip("params.SkipDelete is true")
		}
		By("deleting the cluster")
		cmd := params.EksctlDeleteCmd.
			WithArgs(
				"cluster",
				"--name", params.ClusterName,
				"--timeout", "1h",
				"--verbose", "3",
			)
		Expect(cmd).To(RunSuccessfully())
	})
})

func createClusterWithNodeGroup(ctx context.Context, clusterName, stackName, ng1, version string, ctl api.ClusterProvider) *api.ClusterVPC {
	timeoutDuration := time.Minute * 30
	publicSubnets, privateSubnets, clusterRoleArn, nodeRoleArn, vpcID, securityGroup := createVPCAndRole(stackName, ctl)

	_, err := ctl.EKS().CreateCluster(ctx, &awseks.CreateClusterInput{
		Name: &clusterName,
		ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
			SubnetIds: append(publicSubnets, privateSubnets...),
		},
		RoleArn: &clusterRoleArn,
		Version: aws.String(version),
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() string {
		out, err := ctl.EKS().DescribeCluster(ctx, &awseks.DescribeClusterInput{
			Name: &clusterName,
		})
		Expect(err).NotTo(HaveOccurred())
		return string(out.Cluster.Status)
	}, timeoutDuration, time.Second*30).Should(Equal("ACTIVE"))

	newVPC := api.NewClusterVPC(false)
	newVPC.ID = vpcID
	newVPC.SecurityGroup = securityGroup

	output, err := ctl.EC2().DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
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

	_, err = ctl.EKS().CreateNodegroup(ctx, &awseks.CreateNodegroupInput{
		NodegroupName: &ng1,
		ClusterName:   &clusterName,
		NodeRole:      &nodeRoleArn,
		Subnets:       publicSubnets,
		ScalingConfig: &ekstypes.NodegroupScalingConfig{
			MaxSize:     aws.Int32(1),
			DesiredSize: aws.Int32(1),
			MinSize:     aws.Int32(1),
		},
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() string {
		out, err := ctl.EKS().DescribeNodegroup(ctx, &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: &ng1,
		})
		Expect(err).NotTo(HaveOccurred())
		return string(out.Nodegroup.Status)
	}, timeoutDuration, time.Second*30).Should(Equal("ACTIVE"))

	return newVPC
}

func createVPCAndRole(stackName string, ctl api.ClusterProvider) ([]string, []string, string, string, string, string) {
	templateBody, err := os.ReadFile("cf-template.yaml")
	Expect(err).NotTo(HaveOccurred())
	createStackInput := &cfn.CreateStackInput{
		StackName: &stackName,
	}
	createStackInput.TemplateBody = aws.String(string(templateBody))
	createStackInput.Capabilities = []types.Capability{types.CapabilityCapabilityIam, types.CapabilityCapabilityNamedIam}

	ctx := context.Background()
	_, err = ctl.CloudFormation().CreateStack(ctx, createStackInput)
	Expect(err).NotTo(HaveOccurred())

	var describeStackOut *cfn.DescribeStacksOutput
	Eventually(func() types.StackStatus {
		describeStackOut, err = ctl.CloudFormation().DescribeStacks(ctx, &cfn.DescribeStacksInput{
			StackName: &stackName,
		})
		Expect(err).NotTo(HaveOccurred())
		return describeStackOut.Stacks[0].StackStatus
	}, time.Minute*10, time.Second*15).Should(Equal(types.StackStatusCreateComplete))

	var clusterRoleARN, nodeRoleARN, vpcID string
	var publicSubnets, privateSubnets, securityGroups []string
	for _, output := range describeStackOut.Stacks[0].Outputs {
		switch *output.OutputKey {
		case "ClusterRoleARN":
			clusterRoleARN = *output.OutputValue
		case "NodeRoleARN":
			nodeRoleARN = *output.OutputValue
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

	return publicSubnets, privateSubnets, clusterRoleARN, nodeRoleARN, vpcID, securityGroups[0]
}

func deleteStack(stackName string, ctl api.ClusterProvider) {
	deleteStackInput := &cfn.DeleteStackInput{
		StackName: &stackName,
	}

	_, err := ctl.CloudFormation().DeleteStack(context.Background(), deleteStackInput)
	Expect(err).NotTo(HaveOccurred())
}

var _ = AfterSuite(func() {
	if !params.SkipCreate && !params.SkipDelete {
		deleteStack(stackName, ctl)
	}
})
