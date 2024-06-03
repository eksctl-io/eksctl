//go:build integration
// +build integration

//revive:disable Not changing package name
package dry_run

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/aws/aws-sdk-go-v2/aws"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"

	_ "embed"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	// No cleanup required for dry-run clusters
	params = tests.NewParams("dry-run")
	if err := api.Register(); err != nil {
		panic(fmt.Errorf("unexpected error registering API scheme: %w", err))
	}
}

func TestDryRun(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const eksVersion = api.LatestVersion

const defaultClusterConfig = `
apiVersion: eksctl.io/v1alpha5
availabilityZones:
- us-west-2a
- us-west-2b
cloudWatch:
  clusterLogging: {}
iam:
  vpcResourceControllerPolicy: true
  withOIDC: false
kind: ClusterConfig
metadata:
  name: %[1]s
  region: us-west-2
kubernetesNetworkConfig:
  ipFamily: IPv4
accessConfig:
  authenticationMode: API_AND_CONFIG_MAP
addonsConfig: {}
nodeGroups:
- amiFamily: AmazonLinux2
  containerRuntime: containerd
  disableIMDSv1: true
  disablePodIMDS: false
  instanceSelector: {}
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: false
      appMeshPreview: false
      autoScaler: false
      awsLoadBalancerController: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  labels:
    alpha.eksctl.io/cluster-name: %[1]s
    alpha.eksctl.io/nodegroup-name: ng-default
  name: ng-default
  privateNetworking: false
  securityGroups:
    withLocal: true
    withShared: true
  ssh:
    allow: false
  volumeIOPS: 3000
  volumeSize: 80
  volumeThroughput: 125
  volumeType: gp3

managedNodeGroups:
- amiFamily: AmazonLinux2
  desiredCapacity: 2
  disableIMDSv1: true
  disablePodIMDS: false
  instanceSelector: {}
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: false
      appMeshPreview: false
      autoScaler: false
      awsLoadBalancerController: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  labels:
    alpha.eksctl.io/cluster-name: %[1]s
    alpha.eksctl.io/nodegroup-name: ng-default
  maxSize: 2
  minSize: 2
  name: ng-default
  privateNetworking: false
  securityGroups:
    withLocal: null
    withShared: null
  ssh:
    allow: false
    publicKeyPath: ""
  tags:
    alpha.eksctl.io/nodegroup-name: ng-default
    alpha.eksctl.io/nodegroup-type: managed
  volumeIOPS: 3000
  volumeSize: 80
  volumeThroughput: 125
  volumeType: gp3

privateCluster:
  enabled: false
vpc:
  autoAllocateIPv6: false
  cidr: 192.168.0.0/16
  clusterEndpoints:
    privateAccess: false
    publicAccess: true
  manageSharedNodeSecurityGroupRules: true
  nat:
    gateway: Single
`

//go:embed assets/cloudformation-vpc.yaml
var cfnVPCTemplate string

var _ = Describe("(Integration) [Dry-Run test]", func() {
	parseOutput := func(output []byte) (*api.ClusterConfig, *api.ClusterConfig) {
		actual, err := eks.ParseConfig(output)
		Expect(err).NotTo(HaveOccurred())
		defaultConfig, err := eks.ParseConfig([]byte(fmt.Sprintf(defaultClusterConfig, params.ClusterName)))
		Expect(err).NotTo(HaveOccurred())
		return actual, defaultConfig
	}

	assertDryRun := func(output []byte, updateConfig func(defaultConfig, actual *api.ClusterConfig)) {
		actual, expected := parseOutput(output)
		updateConfig(expected, actual)
		Expect(actual).To(Equal(expected))
	}

	DescribeTable("`create cluster` with --dry-run", func(updateDefaultConfig func(*api.ClusterConfig), createArgs ...string) {
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--dry-run",
				"--version", eksVersion,
				"--name",
				params.ClusterName,
				"--zones", "us-west-2a,us-west-2b",
			).
			WithArgs(createArgs...)

		session := cmd.Run()
		Expect(session.ExitCode()).To(Equal(0))

		output := session.Buffer().Contents()
		assertDryRun(output, func(c, _ *api.ClusterConfig) {
			c.Metadata.Version = eksVersion
			updateDefaultConfig(c)
		})
	},
		Entry("default values", func(c *api.ClusterConfig) {
			c.NodeGroups = nil
		}, "--nodegroup-name=ng-default"),

		Entry("managed nodegroup defaults", func(c *api.ClusterConfig) {
			c.NodeGroups = nil
			ng := c.ManagedNodeGroups[0]
			ng.VolumeSize = aws.Int(101)
		}, "--nodegroup-name=ng-default", "--node-volume-size=101"),

		Entry("override nodegroup defaults", func(c *api.ClusterConfig) {
			c.NodeGroups = nil
			ng := c.ManagedNodeGroups[0]
			ng.IAM.WithAddonPolicies.ExternalDNS = aws.Bool(true)
			ng.VolumeSize = aws.Int(42)
			ng.PrivateNetworking = true
		}, "--nodegroup-name=ng-default", "--node-volume-size=42",
			"--external-dns-access", "--node-private-networking"),

		Entry("override cluster-wide defaults", func(c *api.ClusterConfig) {
			c.ManagedNodeGroups = nil
			c.NodeGroups = nil
			cidr, err := ipnet.ParseCIDR("192.168.0.0/24")
			ExpectWithOffset(1, err).NotTo(HaveOccurred(), "unexpected error parsing CIDR")
			c.VPC.CIDR = cidr
			c.VPC.NAT.Gateway = aws.String("HighlyAvailable")
			c.IAM.WithOIDC = aws.Bool(true)
		}, "--vpc-cidr=192.168.0.0/24", "--without-nodegroup", "--vpc-nat-mode=HighlyAvailable", "--with-oidc"),
	)

	DescribeTable("Flags incompatible with dry-run", func(flag string) {
		cmd := params.EksctlCreateCmd.
			WithArgs("cluster", "--dry-run").
			WithArgs(flag)

		// TODO consider using a custom matcher
		session := cmd.Run()
		Expect(session.ExitCode()).NotTo(Equal(0))
		output := string(session.Err.Contents())
		Expect(output).To(ContainSubstring(fmt.Sprintf("cannot use %s with --dry-run", strings.Split(flag, "=")[0])))

	},
		Entry("install NVIDIA plugin", "--install-nvidia-plugin"),
		Entry("install Neuron plugin", "--install-neuron-plugin"),
		Entry("set kubectl config", "--set-kubeconfig-context"),
		Entry("write kubectl config", "--write-kubeconfig"),
		Entry("set CFN flag", "--cfn-disable-rollback"),
		Entry("set profile option", "--profile=aws"),
	)

	DescribeTable("create cluster with instance selector options", func(setValues func(actual, expected *api.ClusterConfig), createArgs ...string) {
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--name",
				params.ClusterName,
				"--version", eksVersion,
				"--dry-run",
				"--zones", "us-west-2a,us-west-2b",
				"--nodegroup-name=ng-default",
			).
			WithArgs(createArgs...)

		session := cmd.Run()
		Expect(session.ExitCode()).To(Equal(0))
		output := session.Buffer().Contents()

		actual, expected := parseOutput(output)
		expected.Metadata.Version = eksVersion
		setValues(actual, expected)
		Expect(actual).To(Equal(expected))

	}, Entry("instance selector options with unmanaged nodegroup", func(actual, expected *api.ClusterConfig) {
		// This does not do an exact match because instance types matching the instance selector criteria may
		// change over time as EC2 adds more instance types
		Expect(actual.NodeGroups[0].InstancesDistribution.InstanceTypes).NotTo(BeEmpty())
		actual.NodeGroups[0].InstancesDistribution.InstanceTypes = nil

		expected.ManagedNodeGroups = nil
		ng := expected.NodeGroups[0]
		ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
			InstanceTypes: nil,
		}
		ng.InstanceSelector = &api.InstanceSelector{
			VCPUs:  2,
			Memory: "4",
		}

	}, "--managed=false", "--instance-selector-vcpus=2", "--instance-selector-memory=4"),

		Entry("instance selector options with managed nodegroup", func(actual, expected *api.ClusterConfig) {
			Expect(actual.ManagedNodeGroups[0].InstanceTypes).NotTo(BeEmpty())
			actual.ManagedNodeGroups[0].InstanceTypes = nil

			expected.NodeGroups = nil
			ng := expected.ManagedNodeGroups[0]
			ng.InstanceTypes = nil
			ng.InstanceType = ""
			ng.InstanceSelector = &api.InstanceSelector{
				VCPUs:           2,
				Memory:          "4",
				CPUArchitecture: "x86_64",
			}
		}, "--managed", "--instance-selector-vcpus=2", "--instance-selector-memory=4", "--instance-selector-cpu-architecture=x86_64"),
	)

	Describe("create cluster and nodegroups from the output of dry-run", func() {
		setClusterLabel := func(np api.NodePool) {
			np.BaseNodeGroup().Labels["alpha.eksctl.io/cluster-name"] = params.ClusterName
		}
		It("create cluster and nodegroups with the output of dry-run", func() {
			By("generating a ClusterConfig using dry-run")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--dry-run",
					"--version", eksVersion,
					"--name="+params.ClusterName,
					"--zones", "us-west-2a,us-west-2b",
					"--nodegroup-name=ng-default",
				)
			session := cmd.Run()
			Expect(session.ExitCode()).To(Equal(0))
			output := session.Buffer().Contents()
			assertDryRun(output, func(c, _ *api.ClusterConfig) {
				c.Metadata.Name = params.ClusterName
				c.Metadata.Version = eksVersion
				c.NodeGroups = nil
				setClusterLabel(c.ManagedNodeGroups[0])
			})

			By("creating a new cluster from the output of dry-run")
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file=-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(output))

			Expect(cmd).To(RunSuccessfully())

			By("generating a nodegroup config using dry-run")
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--cluster="+params.ClusterName,
					"--name=private-ng",
					"--node-private-networking",
					"--node-volume-size=82",
					"--instance-selector-vcpus=2",
					"--instance-selector-memory=4",
					"--dry-run",
				).
				WithoutArg("--region", params.Region)

			session = cmd.Run()
			Expect(session.ExitCode()).To(Equal(0))
			output = session.Buffer().Contents()
			assertDryRun(output, func(c, actual *api.ClusterConfig) {
				c.Metadata.Name = params.ClusterName
				c.Metadata.Version = eksVersion
				c.VPC = nil
				c.IAM = nil
				c.CloudWatch = nil
				c.PrivateCluster = nil
				c.NodeGroups = nil
				c.AvailabilityZones = nil
				c.KubernetesNetworkConfig = nil
				c.AccessConfig = nil

				ng := c.ManagedNodeGroups[0]
				actualNG := actual.ManagedNodeGroups[0]
				Expect(actualNG.InstanceTypes).NotTo(BeEmpty())
				actualNG.InstanceTypes = nil
				ng.InstanceType = ""
				ng.InstanceSelector = &api.InstanceSelector{
					VCPUs:  2,
					Memory: "4",
				}
				ng.Name = "private-ng"
				ng.PrivateNetworking = true
				ng.VolumeSize = aws.Int(82)
				setClusterLabel(ng)
				setNodeNameKey := func(values map[string]string) {
					values["alpha.eksctl.io/nodegroup-name"] = "private-ng"
				}
				setNodeNameKey(ng.Labels)
				setNodeNameKey(ng.Tags)
			})

			By("creating a new nodegroup from the output of dry-run")
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file=-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(output))
			Expect(cmd).To(RunSuccessfully())
		})
	})

	Describe("`create cluster` with --dry-run and custom subnets", func() {
		stackName := aws.String(params.ClusterName)

		var vpcConfig struct {
			vpcID          string
			publicSubnets  []string
			privateSubnets []string
		}

		BeforeEach(func() {
			By("creating a VPC and two public and private subnets")
			cfn := cloudformation.NewFromConfig(matchers.NewConfig(params.Region))
			ctx := context.Background()
			_, err := cfn.CreateStack(ctx, &cloudformation.CreateStackInput{
				StackName:    aws.String(params.ClusterName),
				TemplateBody: aws.String(cfnVPCTemplate),
				Parameters: []cfntypes.Parameter{
					{
						ParameterKey:   aws.String("EnvironmentName"),
						ParameterValue: aws.String(params.ClusterName),
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() error {
				_, err := cfn.DeleteStack(ctx, &cloudformation.DeleteStackInput{
					StackName: stackName,
				})
				if err != nil {
					return fmt.Errorf("unexpected error deleting stack: %w", err)
				}
				return nil
			})

			By("waiting for the stack to be created successfully")
			waiter := cloudformation.NewStackCreateCompleteWaiter(cfn)
			output, err := waiter.WaitForOutput(ctx, &cloudformation.DescribeStacksInput{
				StackName: stackName,
			}, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(output.Stacks).To(HaveLen(1))

			for _, output := range output.Stacks[0].Outputs {
				switch *output.OutputKey {
				case "VpcId":
					vpcConfig.vpcID = *output.OutputValue
				case "PublicSubnets":
					vpcConfig.publicSubnets = strings.Split(*output.OutputValue, ",")
				case "PrivateSubnets":
					vpcConfig.privateSubnets = strings.Split(*output.OutputValue, ",")
				}
			}
			Expect(vpcConfig.vpcID != "" && len(vpcConfig.publicSubnets) == 2 && len(vpcConfig.privateSubnets) == 2).To(BeTrue(), "expected to find output values for VPC, and public and private subnets")
		})

		It("should generate config that works with `create cluster`", func() {
			By("creating a cluster with --dry-run and custom subnets")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--dry-run",
					"--version", eksVersion,
					"--name",
					params.ClusterName,
					"--nodegroup-name", "ng-default",
					"--vpc-private-subnets", strings.Join(vpcConfig.privateSubnets, ","),
					"--vpc-public-subnets", strings.Join(vpcConfig.publicSubnets, ","),
				)

			session := cmd.Run()
			Expect(session.ExitCode()).To(Equal(0))

			output := session.Buffer().Contents()
			assertDryRun(output, func(c, _ *api.ClusterConfig) {
				c.Metadata.Version = eksVersion
				c.NodeGroups = nil
				c.AvailabilityZones = nil
				c.VPC.NAT = &api.ClusterNAT{
					Gateway: aws.String("Disable"),
				}
				c.VPC.CIDR = ipnet.MustParseCIDR("10.192.0.0/16")
				c.VPC.ID = vpcConfig.vpcID
				c.VPC.Subnets = &api.ClusterSubnets{
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   vpcConfig.privateSubnets[0],
							CIDR: ipnet.MustParseCIDR("10.192.20.0/24"),
							AZ:   "us-west-2a",
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   vpcConfig.privateSubnets[1],
							CIDR: ipnet.MustParseCIDR("10.192.21.0/24"),
							AZ:   "us-west-2b",
						},
					},
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   vpcConfig.publicSubnets[0],
							CIDR: ipnet.MustParseCIDR("10.192.10.0/24"),
							AZ:   "us-west-2a",
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   vpcConfig.publicSubnets[1],
							CIDR: ipnet.MustParseCIDR("10.192.11.0/24"),
							AZ:   "us-west-2b",
						},
					},
				}
			})

			By("ensuring that passing the generated config to dry-run works")
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file=-",
					"--dry-run",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(output))

			Expect(cmd).To(RunSuccessfully())
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
