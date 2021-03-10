package unowned

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/weaveworks/eksctl/pkg/eks"
)

type Cluster struct {
	cfg              *api.ClusterConfig
	ctl              api.ClusterProvider
	clusterName      string
	clusterStackName string
	publicSubnets    []string
	privateSubnets   []string
	clusterRoleARN   string
	nodeRoleARN      string
	VPC              *api.ClusterVPC
}

var timeoutDuration = time.Minute * 30

func NewCluster(cfg *api.ClusterConfig) *Cluster {
	stackName := fmt.Sprintf("eksctl-%s", cfg.Metadata.Name)

	clusterProvider, err := eks.New(&api.ProviderConfig{Region: cfg.Metadata.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())
	ctl := clusterProvider.Provider
	publicSubnets, privateSubnets, clusterRoleARN, nodeRoleARN, vpc := createVPCAndRole(stackName, ctl)

	uc := &Cluster{
		cfg:              cfg,
		ctl:              ctl,
		clusterStackName: stackName,
		clusterName:      cfg.Metadata.Name,
		publicSubnets:    publicSubnets,
		privateSubnets:   privateSubnets,
		clusterRoleARN:   clusterRoleARN,
		nodeRoleARN:      nodeRoleARN,
		VPC:              vpc,
	}

	uc.createCluster()
	return uc
}

func (uc *Cluster) DeleteStack() {
	deleteStackInput := &cfn.DeleteStackInput{
		StackName: &uc.clusterStackName,
	}

	_, err := uc.ctl.CloudFormation().DeleteStack(deleteStackInput)
	Expect(err).NotTo(HaveOccurred())
}

func (uc *Cluster) createCluster() {
	_, err := uc.ctl.EKS().CreateCluster(&awseks.CreateClusterInput{
		Name: &uc.clusterName,
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			SubnetIds: aws.StringSlice(append(uc.publicSubnets, uc.privateSubnets...)),
		},
		RoleArn: &uc.clusterRoleARN,
		Version: &uc.cfg.Metadata.Version,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() string {
		out, err := uc.ctl.EKS().DescribeCluster(&awseks.DescribeClusterInput{
			Name: &uc.clusterName,
		})
		Expect(err).NotTo(HaveOccurred())
		return *out.Cluster.Status
	}, timeoutDuration, time.Minute*4).Should(Equal("ACTIVE"))
}

func (uc Cluster) CreateNodegroups(names ...string) {
	for _, name := range names {
		_, err := uc.ctl.EKS().CreateNodegroup(&awseks.CreateNodegroupInput{
			NodegroupName: &name,
			ClusterName:   &uc.clusterName,
			NodeRole:      &uc.nodeRoleARN,
			Subnets:       aws.StringSlice(uc.publicSubnets),
			ScalingConfig: &awseks.NodegroupScalingConfig{
				MaxSize:     aws.Int64(1),
				DesiredSize: aws.Int64(1),
				MinSize:     aws.Int64(1),
			},
		})
		Expect(err).NotTo(HaveOccurred())
	}

	for _, name := range names {
		Eventually(func() string {
			out, err := uc.ctl.EKS().DescribeNodegroup(&awseks.DescribeNodegroupInput{
				ClusterName:   &uc.clusterName,
				NodegroupName: &name,
			})
			Expect(err).NotTo(HaveOccurred())
			return *out.Nodegroup.Status
		}, timeoutDuration, time.Second*30).Should(Equal("ACTIVE"))
	}
}

func createVPCAndRole(stackName string, ctl api.ClusterProvider) ([]string, []string, string, string, *api.ClusterVPC) {
	templateBody, err := ioutil.ReadFile("../../utilities/unowned/cf-template.yaml")
	Expect(err).NotTo(HaveOccurred())
	createStackInput := &cfn.CreateStackInput{
		StackName: &stackName,
	}
	createStackInput.SetTemplateBody(string(templateBody))
	createStackInput.SetCapabilities(aws.StringSlice([]string{cfn.CapabilityCapabilityIam}))
	createStackInput.SetCapabilities(aws.StringSlice([]string{cfn.CapabilityCapabilityNamedIam}))

	_, err = ctl.CloudFormation().CreateStack(createStackInput)
	Expect(err).NotTo(HaveOccurred())

	var describeStackOut *cfn.DescribeStacksOutput
	Eventually(func() string {
		describeStackOut, err = ctl.CloudFormation().DescribeStacks(&cfn.DescribeStacksInput{
			StackName: &stackName,
		})
		Expect(err).NotTo(HaveOccurred())
		return *describeStackOut.Stacks[0].StackStatus
	}, time.Minute*10, time.Second*15).Should(Equal(cfn.StackStatusCreateComplete))

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

	newVPC := api.NewClusterVPC()
	newVPC.ID = vpcID
	newVPC.SecurityGroup = securityGroups[0]
	newVPC.Subnets = &api.ClusterSubnets{
		Public: api.AZSubnetMapping{
			"public1": api.AZSubnetSpec{
				ID: publicSubnets[0],
			},
			"public2": api.AZSubnetSpec{
				ID: publicSubnets[1],
			},
			"public3": api.AZSubnetSpec{
				ID: publicSubnets[2],
			},
		},
		Private: api.AZSubnetMapping{
			"private4": api.AZSubnetSpec{
				ID: privateSubnets[0],
			},
			"private5": api.AZSubnetSpec{
				ID: privateSubnets[1],
			},
			"private6": api.AZSubnetSpec{
				ID: privateSubnets[2],
			},
		},
	}

	return publicSubnets, privateSubnets, clusterRoleARN, nodeRoleARN, newVPC
}
