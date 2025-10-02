package builder

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Cluster Security Group Configuration Handling", func() {
	var (
		clusterConfig *api.ClusterConfig
		resourceSet   *ClusterResourceSet
		mockProvider  *mockprovider.MockProvider
	)

	BeforeEach(func() {
		mockProvider = mockprovider.NewMockProvider()
		clusterConfig = api.NewClusterConfig()
		clusterConfig.Metadata.Name = "test-cluster"
		clusterConfig.Metadata.Region = "us-west-2"

		resourceSet = NewClusterResourceSet(
			mockProvider.MockEC2(),
			mockProvider.MockSTS(),
			"us-west-2",
			clusterConfig,
			nil,
			false,
		)
	})

	Describe("Karpenter configuration detection", func() {
		Context("when Karpenter is enabled", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = &api.Karpenter{
					Version: "v0.20.0",
				}
			})

			It("should detect Karpenter is enabled", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeTrue())
			})
		})

		Context("when Karpenter is not configured", func() {
			It("should detect Karpenter is not enabled", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeFalse())
			})
		})

		Context("when Karpenter is configured but version is empty", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = &api.Karpenter{
					Version: "",
				}
			})

			It("should detect Karpenter is not enabled", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeFalse())
			})
		})

		Context("when Karpenter is nil", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = nil
			})

			It("should detect Karpenter is not enabled", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeFalse())
			})
		})
	})

	Describe("Karpenter discovery metadata tag detection", func() {
		Context("when karpenter.sh/discovery tag exists in metadata.tags", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "test-cluster",
					"environment":            "test",
				}
			})

			It("should detect the discovery tag exists", func() {
				Expect(resourceSet.hasKarpenterDiscoveryMetadataTag()).To(BeTrue())
			})
		})

		Context("when karpenter.sh/discovery tag does not exist", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"environment": "test",
				}
			})

			It("should detect the discovery tag does not exist", func() {
				Expect(resourceSet.hasKarpenterDiscoveryMetadataTag()).To(BeFalse())
			})
		})

		Context("when metadata.tags is nil", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = nil
			})

			It("should detect the discovery tag does not exist", func() {
				Expect(resourceSet.hasKarpenterDiscoveryMetadataTag()).To(BeFalse())
			})
		})

		Context("when metadata is nil", func() {
			BeforeEach(func() {
				clusterConfig.Metadata = nil
			})

			It("should detect the discovery tag does not exist", func() {
				Expect(resourceSet.hasKarpenterDiscoveryMetadataTag()).To(BeFalse())
			})
		})
	})

	Describe("Combined configuration logic", func() {
		Context("when both Karpenter is enabled AND discovery tag exists", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = &api.Karpenter{
					Version: "v0.20.0",
				}
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "test-cluster",
				}
			})

			It("should determine tagging should be enabled", func() {
				Expect(resourceSet.shouldAddKarpenterDiscoveryTags()).To(BeTrue())
			})
		})

		Context("when only Karpenter is enabled but no discovery tag", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = &api.Karpenter{
					Version: "v0.20.0",
				}
				clusterConfig.Metadata.Tags = map[string]string{
					"environment": "test",
				}
			})

			It("should determine tagging should NOT be enabled", func() {
				Expect(resourceSet.shouldAddKarpenterDiscoveryTags()).To(BeFalse())
			})
		})

		Context("when only discovery tag exists but Karpenter is not enabled", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "test-cluster",
				}
			})

			It("should determine tagging should NOT be enabled", func() {
				Expect(resourceSet.shouldAddKarpenterDiscoveryTags()).To(BeFalse())
			})
		})

		Context("when neither Karpenter is enabled nor discovery tag exists", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"environment": "test",
				}
			})

			It("should determine tagging should NOT be enabled", func() {
				Expect(resourceSet.shouldAddKarpenterDiscoveryTags()).To(BeFalse())
			})
		})
	})

	Describe("Tag generation logic", func() {
		Context("when discovery tag exists with valid value", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "my-cluster-name",
					"environment":            "production",
				}
			})

			It("should generate the correct tag", func() {
				tags := resourceSet.generateKarpenterDiscoveryTags()
				Expect(tags).To(HaveLen(1))
				Expect(tags[0].Key.String()).To(Equal("karpenter.sh/discovery"))
				Expect(tags[0].Value.String()).To(Equal("my-cluster-name"))
			})
		})

		Context("when discovery tag does not exist", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"environment": "production",
				}
			})

			It("should return nil", func() {
				tags := resourceSet.generateKarpenterDiscoveryTags()
				Expect(tags).To(BeNil())
			})
		})

		Context("when metadata.tags is nil", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = nil
			})

			It("should return nil", func() {
				tags := resourceSet.generateKarpenterDiscoveryTags()
				Expect(tags).To(BeNil())
			})
		})

		Context("when metadata is nil", func() {
			BeforeEach(func() {
				clusterConfig.Metadata = nil
			})

			It("should return nil", func() {
				tags := resourceSet.generateKarpenterDiscoveryTags()
				Expect(tags).To(BeNil())
			})
		})
	})

	Describe("Tag generation with validation", func() {
		Context("when configuration is valid", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "valid-cluster-name",
				}
			})

			It("should generate tags successfully", func() {
				tags, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveLen(1))
				Expect(tags[0].Key.String()).To(Equal("karpenter.sh/discovery"))
				Expect(tags[0].Value.String()).To(Equal("valid-cluster-name"))
			})
		})

		Context("when metadata is nil", func() {
			BeforeEach(func() {
				clusterConfig.Metadata = nil
			})

			It("should return an error", func() {
				_, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cluster metadata is nil"))
			})
		})

		Context("when metadata.tags is nil", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = nil
			})

			It("should return an error", func() {
				_, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cluster metadata tags are nil"))
			})
		})

		Context("when discovery tag does not exist", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"environment": "test",
				}
			})

			It("should return an error", func() {
				_, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("karpenter.sh/discovery tag not found in metadata.tags"))
			})
		})

		Context("when discovery tag value is empty", func() {
			BeforeEach(func() {
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "",
				}
			})

			It("should return an error", func() {
				_, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("karpenter.sh/discovery tag value cannot be empty"))
			})
		})

		Context("when discovery tag value exceeds maximum length", func() {
			BeforeEach(func() {
				// Create a string longer than 256 characters
				longValue := make([]byte, 257)
				for i := range longValue {
					longValue[i] = 'a'
				}
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": string(longValue),
				}
			})

			It("should return an error", func() {
				_, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("karpenter.sh/discovery tag value exceeds maximum length of 256 characters"))
			})
		})

		Context("when discovery tag value is at maximum length", func() {
			BeforeEach(func() {
				// Create a string exactly 256 characters
				maxLengthValue := make([]byte, 256)
				for i := range maxLengthValue {
					maxLengthValue[i] = 'a'
				}
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": string(maxLengthValue),
				}
			})

			It("should generate tags successfully", func() {
				tags, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveLen(1))
				Expect(tags[0].Key.String()).To(Equal("karpenter.sh/discovery"))
				Expect(len(tags[0].Value.String())).To(Equal(256))
			})
		})
	})

	Describe("Configuration merging scenarios", func() {
		Context("when cluster has multiple metadata tags including discovery tag", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = &api.Karpenter{
					Version:                   "v0.20.0",
					WithSpotInterruptionQueue: api.Enabled(),
				}
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "production-cluster",
					"environment":            "production",
					"team":                   "platform",
					"cost-center":            "engineering",
				}
			})

			It("should only generate the discovery tag for security group", func() {
				tags := resourceSet.generateKarpenterDiscoveryTags()
				Expect(tags).To(HaveLen(1))
				Expect(tags[0].Key.String()).To(Equal("karpenter.sh/discovery"))
				Expect(tags[0].Value.String()).To(Equal("production-cluster"))
			})

			It("should detect both conditions are met", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeTrue())
				Expect(resourceSet.hasKarpenterDiscoveryMetadataTag()).To(BeTrue())
				Expect(resourceSet.shouldAddKarpenterDiscoveryTags()).To(BeTrue())
			})
		})

		Context("when cluster has Karpenter with additional configuration", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = &api.Karpenter{
					Version:                   "v0.20.0",
					CreateServiceAccount:      api.Enabled(),
					WithSpotInterruptionQueue: api.Enabled(),
					DefaultInstanceProfile:    aws.String("custom-instance-profile"),
				}
				clusterConfig.Metadata.Tags = map[string]string{
					"karpenter.sh/discovery": "advanced-cluster",
				}
			})

			It("should still detect Karpenter is enabled regardless of additional config", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeTrue())
			})

			It("should generate tags correctly with complex Karpenter config", func() {
				tags, err := resourceSet.generateKarpenterDiscoveryTagsWithValidation()
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveLen(1))
				Expect(tags[0].Key.String()).To(Equal("karpenter.sh/discovery"))
				Expect(tags[0].Value.String()).To(Equal("advanced-cluster"))
			})
		})
	})

	Describe("Default value assignment", func() {
		Context("when no configuration is provided", func() {
			It("should have default values that result in no tagging", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeFalse())
				Expect(resourceSet.hasKarpenterDiscoveryMetadataTag()).To(BeFalse())
				Expect(resourceSet.shouldAddKarpenterDiscoveryTags()).To(BeFalse())
			})
		})

		Context("when partial configuration is provided", func() {
			BeforeEach(func() {
				clusterConfig.Karpenter = &api.Karpenter{
					// Version is empty, so Karpenter should not be considered enabled
					CreateServiceAccount: api.Enabled(),
				}
			})

			It("should not enable tagging with incomplete Karpenter config", func() {
				Expect(resourceSet.isKarpenterEnabled()).To(BeFalse())
				Expect(resourceSet.shouldAddKarpenterDiscoveryTags()).To(BeFalse())
			})
		})
	})
})
