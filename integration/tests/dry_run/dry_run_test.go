//go:build integration
// +build integration

package dry_run

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	. "github.com/weaveworks/eksctl/integration/runner"

	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"

	. "github.com/onsi/ginkgo"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	// No cleanup required for dry-run clusters
	params = tests.NewParams("dry-run")
	if err := api.Register(); err != nil {
		panic(errors.Wrap(err, "unexpected error registering API scheme"))
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
nodeGroups:
- amiFamily: AmazonLinux2
  containerRuntime: dockerd
  disableIMDSv1: false
  disablePodIMDS: false
  instanceSelector: {}
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: false
      appMeshPreview: false
      autoScaler: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  instanceType: m5.large
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
  disableIMDSv1: false
  disablePodIMDS: false
  instanceSelector: {}
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: false
      appMeshPreview: false
      autoScaler: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  instanceType: m5.large
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
		ng.InstanceType = "mixed"
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
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
