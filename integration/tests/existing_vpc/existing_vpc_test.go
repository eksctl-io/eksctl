//go:build integration
// +build integration

package unowned

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/weaveworks/eksctl/pkg/eks"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/aws/aws-sdk-go/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

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
	params = tests.NewParams("exist-vpc")
}

func TestVPC(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [using existing VPC]", func() {
	var (
		stackName  string
		ng1        = "ng-1"
		mng1       = "mng-1"
		ctl        api.ClusterProvider
		configFile *os.File
		cfg        *api.ClusterConfig
	)

	BeforeSuite(func() {
		stackName = fmt.Sprintf("eksctl-%s", params.ClusterName)
		cfg = &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   params.ClusterName,
				Region: params.Region,
			},
		}

		var err error
		configFile, err = ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())

		clusterProvider, err := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
		Expect(err).NotTo(HaveOccurred())
		ctl = clusterProvider.Provider
		cfg.VPC = createVPC(stackName, ctl)

		configData, err := json.Marshal(&cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(configFile.Name(), configData, 0755)).To(Succeed())
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--config-file", configFile.Name(),
				"--verbose", "2",
			).WithoutArg("--region", params.Region)
		Expect(cmd).To(RunSuccessfully())
	})

	AfterSuite(func() {
		cmd := params.EksctlDeleteClusterCmd.
			WithArgs(
				"--config-file", configFile.Name(),
				"--wait",
			).WithoutArg("--region", params.Region)
		Expect(cmd).To(RunSuccessfully())
		deleteStack(stackName, ctl)
		Expect(os.RemoveAll(configFile.Name())).To(Succeed())
	})

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

		By("writing the config file")
		configData, err := json.Marshal(&cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(configFile.Name(), configData, 0755)).To(Succeed())

		By("creating the nodegroups")
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"nodegroup",
				"--config-file", configFile.Name(),
				"--verbose", "2",
			).WithoutArg("--region", params.Region)
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

	return newVPC
}

func createVPCStackAndGetOutputs(stackName string, ctl api.ClusterProvider) ([]string, []string, string, string) {
	templateBody, err := ioutil.ReadFile("cf-template.yaml")
	Expect(err).NotTo(HaveOccurred())
	createStackInput := &cfn.CreateStackInput{
		StackName: &stackName,
	}
	By("creating the stack")
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

	_, err := ctl.CloudFormation().DeleteStack(deleteStackInput)
	Expect(err).NotTo(HaveOccurred())
}
