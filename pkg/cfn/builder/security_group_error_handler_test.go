package builder

import (
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SecurityGroupErrorHandler", func() {
	var (
		handler     *SecurityGroupErrorHandler
		clusterName string
	)

	BeforeEach(func() {
		clusterName = "test-cluster"
		handler = NewSecurityGroupErrorHandler(clusterName)
	})

	Describe("NewSecurityGroupErrorHandler", func() {
		It("should create a new handler with the correct cluster name", func() {
			Expect(handler.clusterName).To(Equal(clusterName))
		})
	})

	Describe("ValidateTaggingPrerequisites", func() {
		Context("when neither Karpenter nor discovery tag are present", func() {
			It("should return nil (no error, just skip tagging)", func() {
				err := handler.ValidateTaggingPrerequisites(false, false, "")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when only Karpenter is enabled but no discovery tag", func() {
			It("should return nil (no error, just skip tagging)", func() {
				err := handler.ValidateTaggingPrerequisites(true, false, "")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when only discovery tag exists but Karpenter is not enabled", func() {
			It("should return nil (no error, just skip tagging)", func() {
				err := handler.ValidateTaggingPrerequisites(false, true, "cluster-name")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when both conditions are met with valid discovery value", func() {
			It("should return nil (validation passes)", func() {
				err := handler.ValidateTaggingPrerequisites(true, true, "valid-cluster-name")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when both conditions are met but discovery value is empty", func() {
			It("should return an error", func() {
				err := handler.ValidateTaggingPrerequisites(true, true, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("karpenter.sh/discovery tag value cannot be empty"))
			})
		})

		Context("when both conditions are met but discovery value exceeds maximum length", func() {
			It("should return an error", func() {
				longValue := make([]byte, 257)
				for i := range longValue {
					longValue[i] = 'a'
				}
				err := handler.ValidateTaggingPrerequisites(true, true, string(longValue))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("karpenter.sh/discovery tag value exceeds maximum length of 256 characters"))
			})
		})

		Context("when both conditions are met with discovery value at maximum length", func() {
			It("should return nil (validation passes)", func() {
				maxLengthValue := make([]byte, 256)
				for i := range maxLengthValue {
					maxLengthValue[i] = 'a'
				}
				err := handler.ValidateTaggingPrerequisites(true, true, string(maxLengthValue))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("WrapConfigurationError", func() {
		Context("when error is nil", func() {
			It("should return nil", func() {
				err := handler.WrapConfigurationError(nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when error contains 'nil'", func() {
			It("should wrap with metadata guidance", func() {
				originalErr := errors.New("metadata is nil")
				wrappedErr := handler.WrapConfigurationError(originalErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("invalid security group configuration for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please ensure your cluster configuration includes valid metadata and tags"))
			})
		})

		Context("when error contains 'empty'", func() {
			It("should wrap with empty tag value guidance", func() {
				originalErr := errors.New("tag value is empty")
				wrappedErr := handler.WrapConfigurationError(originalErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("invalid security group configuration for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("The karpenter.sh/discovery tag value cannot be empty"))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please provide a valid tag value, typically the cluster name"))
			})
		})

		Context("when error contains 'length'", func() {
			It("should wrap with length constraint guidance", func() {
				originalErr := errors.New("tag value exceeds maximum length")
				wrappedErr := handler.WrapConfigurationError(originalErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("invalid security group configuration for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("AWS tag values must be 256 characters or less"))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please use a shorter value for the karpenter.sh/discovery tag"))
			})
		})

		Context("when error is generic", func() {
			It("should wrap with generic guidance", func() {
				originalErr := errors.New("some generic error")
				wrappedErr := handler.WrapConfigurationError(originalErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("invalid security group configuration for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("some generic error"))
			})
		})
	})

	Describe("WrapTemplateGenerationError", func() {
		Context("when error is nil", func() {
			It("should return nil", func() {
				err := handler.WrapTemplateGenerationError("test operation", nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when error is InvalidParameterValue AWS error", func() {
			It("should wrap with parameter validation guidance", func() {
				awsErr := &mockAPIError{
					code:    "InvalidParameterValue",
					message: "Invalid parameter",
				}
				wrappedErr := handler.WrapTemplateGenerationError("generate template", awsErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("failed to generate template for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("invalid parameter value"))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please check your cluster configuration"))
			})
		})

		Context("when error is InvalidVpcId.NotFound AWS error", func() {
			It("should wrap with VPC guidance", func() {
				awsErr := &mockAPIError{
					code:    "InvalidVpcId.NotFound",
					message: "VPC not found",
				}
				wrappedErr := handler.WrapTemplateGenerationError("create security group", awsErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("failed to create security group for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("VPC not found"))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please ensure the VPC exists"))
			})
		})

		Context("when error is UnauthorizedOperation AWS error", func() {
			It("should wrap with IAM permissions guidance", func() {
				awsErr := &mockAPIError{
					code:    "UnauthorizedOperation",
					message: "Access denied",
				}
				wrappedErr := handler.WrapTemplateGenerationError("create tags", awsErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("failed to create tags for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("insufficient permissions"))
				Expect(wrappedErr.Error()).To(ContainSubstring("ec2:CreateSecurityGroup, ec2:CreateTags, ec2:DescribeSecurityGroups"))
			})
		})

		Context("when error contains 'metadata'", func() {
			It("should wrap with metadata guidance", func() {
				originalErr := errors.New("metadata configuration error")
				wrappedErr := handler.WrapTemplateGenerationError("process metadata", originalErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("failed to process metadata for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("cluster metadata is missing or invalid"))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please ensure your cluster configuration includes valid metadata with the required tags"))
			})
		})

		Context("when error contains 'karpenter.sh/discovery'", func() {
			It("should wrap with Karpenter discovery guidance", func() {
				originalErr := errors.New("karpenter.sh/discovery tag issue")
				wrappedErr := handler.WrapTemplateGenerationError("configure tagging", originalErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("failed to configure tagging for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("To enable automatic security group tagging"))
				Expect(wrappedErr.Error()).To(ContainSubstring("1) Karpenter is enabled (karpenter.version is specified)"))
				Expect(wrappedErr.Error()).To(ContainSubstring("2) karpenter.sh/discovery tag is present in metadata.tags"))
			})
		})

		Context("when error is generic", func() {
			It("should wrap with generic troubleshooting guidance", func() {
				originalErr := errors.New("generic template error")
				wrappedErr := handler.WrapTemplateGenerationError("build template", originalErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("failed to build template for cluster \"test-cluster\""))
				Expect(wrappedErr.Error()).To(ContainSubstring("generic template error"))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please check your cluster configuration"))
				Expect(wrappedErr.Error()).To(ContainSubstring("ec2:CreateSecurityGroup, ec2:CreateTags, ec2:DescribeSecurityGroups"))
			})
		})
	})

	Describe("Error handling for nested configuration structures", func() {
		Context("when validating complex nested configuration scenarios", func() {
			It("should handle nil metadata in nested structure", func() {
				err := handler.ValidateTaggingPrerequisites(true, true, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("karpenter.sh/discovery tag value cannot be empty"))
			})

			It("should handle deeply nested configuration errors", func() {
				nestedErr := errors.New("nested configuration: metadata.tags.karpenter.sh/discovery is nil")
				wrappedErr := handler.WrapConfigurationError(nestedErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("invalid security group configuration"))
				Expect(wrappedErr.Error()).To(ContainSubstring("Please ensure your cluster configuration includes valid metadata and tags"))
			})

			It("should handle configuration merging errors", func() {
				mergingErr := errors.New("configuration merging failed: empty tag value after merge")
				wrappedErr := handler.WrapConfigurationError(mergingErr)
				Expect(wrappedErr).To(HaveOccurred())
				Expect(wrappedErr.Error()).To(ContainSubstring("The karpenter.sh/discovery tag value cannot be empty"))
			})
		})
	})

	Describe("Default value assignment error scenarios", func() {
		Context("when default values cause validation errors", func() {
			It("should handle empty default values", func() {
				err := handler.ValidateTaggingPrerequisites(true, true, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("karpenter.sh/discovery tag value cannot be empty"))
			})

			It("should handle invalid default tag length", func() {
				longDefaultValue := make([]byte, 300) // Exceeds 256 character limit
				for i := range longDefaultValue {
					longDefaultValue[i] = 'x'
				}
				err := handler.ValidateTaggingPrerequisites(true, true, string(longDefaultValue))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("karpenter.sh/discovery tag value exceeds maximum length"))
			})
		})
	})
})

// mockAPIError implements smithy.APIError for testing
type mockAPIError struct {
	code    string
	message string
}

func (e *mockAPIError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *mockAPIError) ErrorCode() string {
	return e.code
}

func (e *mockAPIError) ErrorMessage() string {
	return e.message
}

func (e *mockAPIError) ErrorFault() smithy.ErrorFault {
	return smithy.FaultClient
}
