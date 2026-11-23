package v1alpha5_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// instanceIDNodeNameConfig returns a nodeadm NodeConfig that sets the InstanceIdNodeName feature gate.
func instanceIDNodeNameConfig(enabled bool) *string {
	return aws.String(fmt.Sprintf(`apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  featureGates:
    InstanceIdNodeName: %t
`, enabled))
}

var _ = Describe("AmazonLinux2027", func() {
	It("is a supported AMI family", func() {
		Expect(api.SupportedAMIFamilies()).To(ContainElement(api.NodeImageFamilyAmazonLinux2027))
	})

	DescribeTable("IsNodeadmImage", func(amiFamily string, expected bool) {
		Expect(api.IsNodeadmImage(amiFamily)).To(Equal(expected))
	},
		Entry(api.NodeImageFamilyAmazonLinux2027, api.NodeImageFamilyAmazonLinux2027, true),
		Entry(api.NodeImageFamilyAmazonLinux2023, api.NodeImageFamilyAmazonLinux2023, true),
		Entry(api.NodeImageFamilyAmazonLinux2, api.NodeImageFamilyAmazonLinux2, false),
		Entry(api.NodeImageFamilyBottlerocket, api.NodeImageFamilyBottlerocket, false),
	)

	DescribeTable("node name and access entry type", func(amiFamily string, overrideBootstrapCommand *string, usesInstanceIDNodeName bool, expectedType api.AccessEntryType) {
		ng := api.NewNodeGroup()
		ng.AMIFamily = amiFamily
		ng.OverrideBootstrapCommand = overrideBootstrapCommand
		Expect(api.UsesInstanceIDNodeName(ng.NodeGroupBase)).To(Equal(usesInstanceIDNodeName))
		Expect(api.GetAccessEntryType(ng)).To(Equal(expectedType))
	},
		Entry("AmazonLinux2027 by default", api.NodeImageFamilyAmazonLinux2027, (*string)(nil), true, api.AccessEntryTypeEC2),
		Entry("AmazonLinux2027 with InstanceIdNodeName off", api.NodeImageFamilyAmazonLinux2027, instanceIDNodeNameConfig(false), false, api.AccessEntryTypeLinux),
		Entry("AmazonLinux2027 with a NodeConfig that doesn't set InstanceIdNodeName", api.NodeImageFamilyAmazonLinux2027,
			aws.String("apiVersion: node.eks.aws/v1alpha1\nkind: NodeConfig\nspec:\n  instance:\n    localStorage:\n      strategy: RAID0\n"), true, api.AccessEntryTypeEC2),
		Entry("AmazonLinux2023 by default", api.NodeImageFamilyAmazonLinux2023, (*string)(nil), false, api.AccessEntryTypeLinux),
		Entry("AmazonLinux2023 with InstanceIdNodeName on", api.NodeImageFamilyAmazonLinux2023, instanceIDNodeNameConfig(true), true, api.AccessEntryTypeEC2),
		Entry("AmazonLinux2 ignores overrideBootstrapCommand", api.NodeImageFamilyAmazonLinux2, aws.String("/etc/eks/bootstrap.sh my-cluster"), false, api.AccessEntryTypeLinux),
		Entry("Windows", api.NodeImageFamilyWindowsServer2022CoreContainer, (*string)(nil), false, api.AccessEntryTypeWindows),
	)

	Describe("self-managed nodegroups", func() {
		var (
			cfg *api.ClusterConfig
			ng  *api.NodeGroup
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			ng = cfg.NewNodeGroup()
			ng.AMIFamily = api.NodeImageFamilyAmazonLinux2027
		})

		DescribeTable("accepts nodegroups without an AMI ID, as eksctl looks up the AMI", func(ami string) {
			ng.AMI = ami
			Expect(api.ValidateNodeGroup(0, ng, cfg)).To(Succeed())
		},
			Entry("ami not set", ""),
			Entry("ami is auto-ssm", api.NodeImageResolverAutoSSM),
		)

		It("accepts an AMI ID without overrideBootstrapCommand", func() {
			ng.AMI = "ami-0123456789abcdef0"
			Expect(api.ValidateNodeGroup(0, ng, cfg)).To(Succeed())
		})

		It("accepts the InstanceIdNodeName feature gate in overrideBootstrapCommand", func() {
			ng.AMI = "ami-0123456789abcdef0"
			ng.OverrideBootstrapCommand = instanceIDNodeNameConfig(false)
			Expect(api.ValidateNodeGroup(0, ng, cfg)).To(Succeed())
		})

		It("rejects the dockerd container runtime", func() {
			ng.AMI = "ami-0123456789abcdef0"
			ng.ContainerRuntime = aws.String(api.ContainerRuntimeDockerD)
			Expect(api.ValidateNodeGroup(0, ng, cfg)).To(MatchError(ContainSubstring(fmt.Sprintf("only %s is supported for container runtime on %s nodes", api.ContainerRuntimeContainerD, api.NodeImageFamilyAmazonLinux2027))))
		})

		It("defaults the container runtime to containerd", func() {
			ng.AMI = "ami-0123456789abcdef0"
			Expect(api.SetNodeGroupDefaults(ng, &api.ClusterMeta{Version: api.Version1_23}, false)).To(Succeed())
			Expect(*ng.ContainerRuntime).To(Equal(api.ContainerRuntimeContainerD))
		})
	})

	It("accepts a managed nodegroup with an AMI ID without overrideBootstrapCommand", func() {
		mng := api.NewManagedNodeGroup()
		mng.AMIFamily = api.NodeImageFamilyAmazonLinux2027
		mng.AMI = "ami-0123456789abcdef0"
		Expect(api.ValidateManagedNodeGroup(0, mng)).To(Succeed())
	})

	It("allows an AmazonLinux2027 managed nodegroup with a launch template, which sets its own node config", func() {
		mng := api.NewManagedNodeGroup()
		mng.AMIFamily = api.NodeImageFamilyAmazonLinux2027
		// A launch template nodegroup must not set these fields; the template provides them.
		mng.VolumeSize = nil
		mng.VolumeType = nil
		mng.LaunchTemplate = &api.LaunchTemplate{ID: "lt-1234"}
		Expect(api.ValidateManagedNodeGroup(0, mng)).To(Succeed())
	})
})
