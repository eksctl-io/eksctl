package v1alpha5_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("AmazonLinux2027", func() {
	It("is a supported AMI family", func() {
		Expect(api.SupportedAMIFamilies()).To(ContainElement(api.NodeImageFamilyAmazonLinux2027))
	})

	DescribeTable("AMI family helpers", func(amiFamily string, usesInstanceIDNodeName, isNodeadm bool) {
		Expect(api.UsesInstanceIDNodeName(amiFamily)).To(Equal(usesInstanceIDNodeName))
		Expect(api.IsNodeadmImage(amiFamily)).To(Equal(isNodeadm))
	},
		Entry(api.NodeImageFamilyAmazonLinux2027, api.NodeImageFamilyAmazonLinux2027, true, true),
		Entry(api.NodeImageFamilyAmazonLinux2023, api.NodeImageFamilyAmazonLinux2023, false, true),
		Entry(api.NodeImageFamilyAmazonLinux2, api.NodeImageFamilyAmazonLinux2, false, false),
		Entry(api.NodeImageFamilyBottlerocket, api.NodeImageFamilyBottlerocket, false, false),
	)

	DescribeTable("access entry type", func(amiFamily string, expectedType api.AccessEntryType) {
		ng := api.NewNodeGroup()
		ng.AMIFamily = amiFamily
		Expect(api.GetAccessEntryType(ng)).To(Equal(expectedType))
	},
		Entry(api.NodeImageFamilyAmazonLinux2027, api.NodeImageFamilyAmazonLinux2027, api.AccessEntryTypeEC2),
		Entry(api.NodeImageFamilyAmazonLinux2023, api.NodeImageFamilyAmazonLinux2023, api.AccessEntryTypeLinux),
		Entry(api.NodeImageFamilyWindowsServer2022CoreContainer, api.NodeImageFamilyWindowsServer2022CoreContainer, api.AccessEntryTypeWindows),
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

		DescribeTable("require an AMI ID", func(ami string) {
			ng.AMI = ami
			Expect(api.ValidateNodeGroup(0, ng, cfg)).To(MatchError("nodeGroups[0].ami must be set to an AMI ID when nodeGroups[0].amiFamily is AmazonLinux2027, as eksctl can't look up AmazonLinux2027 AMIs yet"))
		},
			Entry("ami not set", ""),
			Entry("ami is auto", api.NodeImageResolverAuto),
			Entry("ami is auto-ssm", api.NodeImageResolverAutoSSM),
		)

		It("accepts an AMI ID without overrideBootstrapCommand", func() {
			ng.AMI = "ami-0123456789abcdef0"
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

	DescribeTable("managed nodegroups are not supported", func(updateMNG func(*api.ManagedNodeGroup)) {
		mng := api.NewManagedNodeGroup()
		mng.AMIFamily = api.NodeImageFamilyAmazonLinux2027
		updateMNG(mng)
		Expect(api.ValidateManagedNodeGroup(0, mng)).To(MatchError("managedNodeGroups[0].amiFamily AmazonLinux2027 is not supported for managed nodegroups yet, use a self-managed nodegroup (nodeGroups) instead"))
	},
		Entry("without an AMI", func(*api.ManagedNodeGroup) {}),
		Entry("with an AMI ID", func(mng *api.ManagedNodeGroup) {
			mng.AMI = "ami-0123456789abcdef0"
		}),
		Entry("with a launch template", func(mng *api.ManagedNodeGroup) {
			mng.LaunchTemplate = &api.LaunchTemplate{ID: "lt-1234"}
		}),
		Entry("with a lowercase AMI family", func(mng *api.ManagedNodeGroup) {
			mng.AMIFamily = "amazonlinux2027"
			mng.AMI = "ami-0123456789abcdef0"
		}),
	)
})
