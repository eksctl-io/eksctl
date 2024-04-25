package v1alpha5_test

import (
	"bytes"
	"fmt"

	"github.com/kris-nova/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("GPU instance support", func() {

	type gpuInstanceEntry struct {
		gpuInstanceType  string
		amiFamily        string
		instanceTypeName string
		instanceSelector *api.InstanceSelector

		expectUnsupportedErr bool
		expectWarning        bool
	}

	assertValidationError := func(e gpuInstanceEntry, err error) {
		if e.expectUnsupportedErr {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("%s instance types are not supported for %s", e.instanceTypeName, e.amiFamily))))
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
		Entry("AL2023 INF", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyAmazonLinux2023,
			gpuInstanceType:      "inf1.xlarge",
			expectUnsupportedErr: true,
			instanceTypeName:     "Inferentia",
		}),
		Entry("AL2023 TRN", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyAmazonLinux2023,
			gpuInstanceType:      "trn1.2xlarge",
			expectUnsupportedErr: true,
			instanceTypeName:     "Trainium",
		}),
		Entry("AL2023 NVIDIA", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyAmazonLinux2023,
			gpuInstanceType:      "g4dn.xlarge",
			expectUnsupportedErr: true,
			instanceTypeName:     "GPU",
		}),
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
			instanceTypeName:     "Inferentia",
		}),
		Entry("Bottlerocket", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyBottlerocket,
			gpuInstanceType:      "trn1.2xlarge",
			expectUnsupportedErr: true,
			instanceTypeName:     "Trainium",
		}),
	)

	DescribeTable("unmanaged nodegroups", func(e gpuInstanceEntry) {
		ng := api.NewNodeGroup()
		ng.InstanceType = e.gpuInstanceType
		ng.AMIFamily = e.amiFamily
		assertValidationError(e, api.ValidateNodeGroup(0, ng, api.NewClusterConfig()))

	},
		Entry("AL2023 INF", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyAmazonLinux2023,
			gpuInstanceType:      "inf1.xlarge",
			expectUnsupportedErr: true,
			instanceTypeName:     "Inferentia",
		}),
		Entry("AL2023 TRN", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyAmazonLinux2023,
			gpuInstanceType:      "trn1.2xlarge",
			expectUnsupportedErr: true,
			instanceTypeName:     "Trainium",
		}),
		Entry("AL2023 NVIDIA", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyAmazonLinux2023,
			gpuInstanceType:      "g4dn.xlarge",
			expectUnsupportedErr: true,
			instanceTypeName:     "GPU",
		}),
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
		Entry("AL2", gpuInstanceEntry{
			gpuInstanceType: "trn1.2xlarge",
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
		Entry("Bottlerocket infra", gpuInstanceEntry{
			amiFamily:            api.NodeImageFamilyBottlerocket,
			gpuInstanceType:      "trn1.2xlarge",
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

	DescribeTable("GPU drivers", func(e gpuInstanceEntry) {
		ng := api.NewNodeGroup()
		ng.AMIFamily = e.amiFamily
		ng.InstanceType = e.gpuInstanceType
		ng.InstanceSelector = e.instanceSelector

		mng := api.NewManagedNodeGroup()
		mng.AMIFamily = e.amiFamily
		mng.InstanceType = e.gpuInstanceType
		mng.InstanceSelector = e.instanceSelector
		if mng.InstanceSelector == nil {
			mng.InstanceSelector = &api.InstanceSelector{}
		}

		output := &bytes.Buffer{}
		logger.Writer = output
		Expect(api.ValidateNodeGroup(0, ng, api.NewClusterConfig())).NotTo(HaveOccurred())
		if e.expectWarning {
			Expect(output.String()).To(ContainSubstring(api.GPUDriversWarning(mng.AMIFamily)))
		} else {
			Expect(output.String()).NotTo(ContainSubstring(api.GPUDriversWarning(mng.AMIFamily)))
		}

		output = &bytes.Buffer{}
		logger.Writer = output
		Expect(api.ValidateManagedNodeGroup(0, mng)).NotTo(HaveOccurred())
		if e.expectWarning {
			Expect(output.String()).To(ContainSubstring(api.GPUDriversWarning(mng.AMIFamily)))
		} else {
			Expect(output.String()).NotTo(ContainSubstring(api.GPUDriversWarning(mng.AMIFamily)))
		}
	},
		Entry("Windows without GPU instances", gpuInstanceEntry{
			amiFamily: api.NodeImageFamilyUbuntu2004,
			instanceSelector: &api.InstanceSelector{
				VCPUs: 4,
				GPUs:  newInt(0),
			},
		}),
		Entry("Windows with explicit GPU instance", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyWindowsServer2019FullContainer,
			gpuInstanceType: "g4dn.xlarge",
			expectWarning:   true,
		}),
		Entry("Windows with implicit GPU instance", gpuInstanceEntry{
			amiFamily: api.NodeImageFamilyWindowsServer2022CoreContainer,
			instanceSelector: &api.InstanceSelector{
				VCPUs: 4,
			},
			expectWarning: true,
		}),
		Entry("Ubuntu with explicit GPU instance", gpuInstanceEntry{
			amiFamily:       api.NodeImageFamilyUbuntu1804,
			gpuInstanceType: "g4dn.xlarge",
			expectWarning:   true,
		}),
		Entry("Ubuntu with implicit GPU instance", gpuInstanceEntry{
			amiFamily: api.NodeImageFamilyUbuntu2004,
			instanceSelector: &api.InstanceSelector{
				VCPUs: 4,
				GPUs:  newInt(2),
			},
			expectWarning: true,
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
