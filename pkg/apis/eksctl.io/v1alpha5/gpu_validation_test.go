package v1alpha5_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
)

var _ = Describe("GPU instance support", func() {

	type gpuInstanceEntry struct {
		gpuInstanceType string
		amiFamily       string
		customErr       string

		expectUnsupportedErr bool
	}

	assertValidationError := func(e gpuInstanceEntry, err error) {
		if e.expectUnsupportedErr {
			Expect(err).To(HaveOccurred())
			if e.customErr != "" {
				Expect(err).To(MatchError(e.customErr))
			} else if instanceutils.IsNvidiaInstanceType(e.gpuInstanceType) {
				Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("GPU instance types are not supported for %s", e.amiFamily))))
			} else {
				Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("Inferentia instance types are not supported for %s", e.amiFamily))))
			}
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
		Entry("Bottlerocket", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyBottlerocket,
			gpuInstanceType:      "g4dn.xlarge",
			expectUnsupportedErr: true,
			customErr:            fmt.Sprintf("NVIDIA GPU instance types are not supported for managed nodegroups with AMIFamily %s", api.NodeImageFamilyBottlerocket),
		}),
		Entry("Ubuntu2004", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyUbuntu2004,
			gpuInstanceType:      "g4dn.xlarge",
			expectUnsupportedErr: true,
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
		assertValidationError(e, api.ValidateNodeGroup(0, ng))

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
		Entry("Bottlerocket", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyBottlerocket,
			gpuInstanceType:      "inf1.xlarge",
			expectUnsupportedErr: true,
		}),
		Entry("Ubuntu2004", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyUbuntu2004,
			gpuInstanceType:      "g4dn.xlarge",
			expectUnsupportedErr: true,
		}),
		Entry("Windows2004", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyWindowsServer2004CoreContainer,
			gpuInstanceType:      "p2.xlarge",
			expectUnsupportedErr: true,
		}),
		Entry("Windows2019Core", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyWindowsServer2019CoreContainer,
			gpuInstanceType:      "g3.8xlarge",
			expectUnsupportedErr: true,
		}),
		Entry("Windows2019Full", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyWindowsServer2019FullContainer,
			gpuInstanceType:      "p3.2xlarge",
			expectUnsupportedErr: true,
		}),
		Entry("Windows20H2Core", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyWindowsServer20H2CoreContainer,
			gpuInstanceType:      "g4dn.xlarge",
			expectUnsupportedErr: true,
		}),
	)
})
