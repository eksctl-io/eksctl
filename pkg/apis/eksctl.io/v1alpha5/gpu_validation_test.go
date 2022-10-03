package v1alpha5_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("GPU instance support", func() {

	type gpuInstanceEntry struct {
		gpuInstanceType string
		amiFamily       string

		expectUnsupportedErr bool
	}

	assertValidationError := func(e gpuInstanceEntry, err error) {
		if e.expectUnsupportedErr {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("Inferentia instance types are not supported for %s", e.amiFamily))))
			return
		}
		Expect(err).NotTo(HaveOccurred())
	}

	DescribeTable("managed nodegroups", func(e gpuInstanceEntry) {
		mng := api.NewManagedNodeGroup()
		mng.InstanceType = e.gpuInstanceType
		mng.AMIFamily = e.amiFamily
		mng.InstanceSelector = &api.InstanceSelector{}
		assertValidationError(e, api.ValidateManagedNodeGroup(0, mng))
	},
		Entry("AL2", gpuInstanceEntry{
			gpuInstanceType: "asdf",
			amiFamily:       api.NodeImageFamilyAmazonLinux2,
		}),
		Entry("AL2", gpuInstanceEntry{
			gpuInstanceType: "g5.12xlarge",
			amiFamily:       api.NodeImageFamilyAmazonLinux2,
		}),
		Entry("Ubuntu2004", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyUbuntu2004,
			gpuInstanceType: "g4dn.xlarge",
		}),
		Entry("Bottlerocket", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyBottlerocket,
			gpuInstanceType:      "inf1.xlarge",
			expectUnsupportedErr: true,
		}),
	)

	DescribeTable("unmanaged nodegroups", func(e gpuInstanceEntry) {
		ng := api.NewNodeGroup()
		ng.InstanceType = e.gpuInstanceType
		ng.AMIFamily = e.amiFamily
		assertValidationError(e, api.ValidateNodeGroup(0, ng, api.NewClusterConfig()))

	},
		Entry("AL2", gpuInstanceEntry{
			gpuInstanceType: "g4dn.xlarge",
			amiFamily:       api.NodeImageFamilyAmazonLinux2,
		}),
		Entry("AL2", gpuInstanceEntry{
			gpuInstanceType: "g5.12xlarge",
			amiFamily:       api.NodeImageFamilyAmazonLinux2,
		}),
		Entry("AL2", gpuInstanceEntry{
			gpuInstanceType: "inf1.xlarge",
			amiFamily:       api.NodeImageFamilyAmazonLinux2,
		}),
		Entry("AMI unset", gpuInstanceEntry{
			gpuInstanceType: "g4dn.xlarge",
		}),
		Entry("Bottlerocket", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyBottlerocket,
			gpuInstanceType: "g4dn.xlarge",
		}),
		Entry("Bottlerocket infra", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyBottlerocket,
			gpuInstanceType:      "inf1.xlarge",
			expectUnsupportedErr: true,
		}),
		Entry("Bottlerocket nvidia", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyBottlerocket,
			gpuInstanceType: "g4dn.xlarge",
		}),
		Entry("Ubuntu2004", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyUbuntu2004,
			gpuInstanceType: "g4dn.xlarge",
		}),
		Entry("Windows2019Core", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyWindowsServer2019CoreContainer,
			gpuInstanceType: "g3.8xlarge",
		}),
		Entry("Windows2019Full", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyWindowsServer2019FullContainer,
			gpuInstanceType: "p3.2xlarge",
		}),
		Entry("Windows2022Core", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyWindowsServer2022CoreContainer,
			gpuInstanceType: "g3.8xlarge",
		}),
		Entry("Windows2022Full", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyWindowsServer2022FullContainer,
			gpuInstanceType: "p3.2xlarge",
		}),
	)

	DescribeTable("ARM-based GPU instance type support", func(amiFamily string, expectErr bool) {
		ng := api.NewNodeGroup()
		ng.InstanceType = "g5g.medium"
		ng.AMIFamily = amiFamily
		err := api.ValidateNodeGroup(0, ng, api.NewClusterConfig())
		if expectErr {
			Expect(err).To(MatchError(fmt.Sprintf("ARM GPU instance types are not supported for unmanaged nodegroups with AMIFamily %s", amiFamily)))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
	},
		Entry("AmazonLinux2", api.NodeImageFamilyAmazonLinux2, true),
		Entry("Ubuntu2004", api.NodeImageFamilyUbuntu2004, true),
		Entry("Ubuntu1804", api.NodeImageFamilyUbuntu1804, true),
		Entry("Windows2019Full", api.NodeImageFamilyWindowsServer2019FullContainer, true),
		Entry("Windows2019Core", api.NodeImageFamilyWindowsServer2019CoreContainer, true),
		Entry("Bottlerocket", api.NodeImageFamilyBottlerocket, false),
	)
})
