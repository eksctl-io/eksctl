package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"k8s.io/client-go/kubernetes/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	utilsstrings "github.com/weaveworks/eksctl/pkg/utils/strings"
)

var _ = Describe("EKS API wrapper", func() {
	var (
		ctl         *ClusterProviderImpl
		cfg         *api.ClusterConfig
		awsProvider *mockprovider.MockAwsProvider
	)

	BeforeEach(func() {
		fakeClientSet := fake.NewSimpleClientset()
		awsProvider = mockprovider.NewMockAwsProvider()
		kubeProvider := mockprovider.NewMockKubeProvider(fakeClientSet)
		ctl = NewWithMocks(awsProvider, kubeProvider)
		cfg = api.NewClusterConfig()
	})

	Describe("can update cluster tags", func() {
		var (
			err      error
			sentTags map[string]string
		)
		BeforeEach(func() {
			updateClusterTagsOutput := &awseks.TagResourceOutput{}

			awsProvider.MockEKS().On("TagResource", mock.MatchedBy(func(input *awseks.TagResourceInput) bool {
				Expect(input.Tags).ToNot(BeEmpty())

				sentTags = utilsstrings.ToValuesMap(input.Tags)

				return true
			})).Return(updateClusterTagsOutput, nil)
			awsProvider.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				return true
			})).Return(&awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster("testcluster", awseks.ClusterStatusActive),
			}, nil)
		})
		It("shouldn't call the API when there are no tags", func() {
			err = ctl.UpdateClusterTags(cfg)
			Expect(sentTags).To(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})
		It("should call the API with a non-empty tags map", func() {
			cfg.Metadata.Tags = map[string]string{
				"env": "prod",
			}
			err = ctl.UpdateClusterTags(cfg)
			Expect(sentTags).To(BeEquivalentTo(cfg.Metadata.Tags))
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("can update cluster configuration for logging", func() {
		var (
			err                error
			sentClusterLogging []*awseks.LogSetup
		)

		BeforeEach(func() {
			fakeClientSet := fake.NewSimpleClientset()
			awsProvider := mockprovider.NewMockAwsProvider()
			kubeProvider := mockprovider.NewMockKubeProvider(fakeClientSet)
			ctl = NewWithMocks(awsProvider, kubeProvider)
			cfg = api.NewClusterConfig()

			updateClusterConfigOutput := &awseks.UpdateClusterConfigOutput{
				Update: &awseks.Update{
					Id:   aws.String("u123"),
					Type: aws.String(awseks.UpdateTypeLoggingUpdate),
				},
			}

			describeClusterOutput := &awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster("testcluster", awseks.ClusterStatusActive),
			}

			describeClusterOutput.Cluster.Logging = &awseks.Logging{
				ClusterLogging: []*awseks.LogSetup{
					{
						Enabled: api.Enabled(),
						Types:   aws.StringSlice([]string{"api", "audit"}),
					},
					{
						Enabled: api.Disabled(),
						Types:   aws.StringSlice([]string{"controllerManager", "authenticator", "scheduler"}),
					},
				},
			}

			awsProvider.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				return true
			})).Return(describeClusterOutput, nil)

			awsProvider.MockEKS().On("UpdateClusterConfig", mock.MatchedBy(func(input *awseks.UpdateClusterConfigInput) bool {
				Expect(input.Logging).ToNot(BeNil())

				Expect(input.Logging.ClusterLogging[0].Enabled).ToNot(BeNil())
				Expect(input.Logging.ClusterLogging[1].Enabled).ToNot(BeNil())

				Expect(*input.Logging.ClusterLogging[0].Enabled).To(BeTrue())
				Expect(*input.Logging.ClusterLogging[1].Enabled).To(BeFalse())

				sentClusterLogging = input.Logging.ClusterLogging

				return true
			})).Return(updateClusterConfigOutput, nil)

			describeUpdateInput := &awseks.DescribeUpdateInput{}

			describeUpdateOutput := &awseks.DescribeUpdateOutput{
				Update: &awseks.Update{
					Id:     aws.String("u123"),
					Type:   aws.String(awseks.UpdateTypeLoggingUpdate),
					Status: aws.String(awseks.UpdateStatusSuccessful),
				},
			}

			awsProvider.MockEKS().On("DescribeUpdateRequest", mock.MatchedBy(func(input *awseks.DescribeUpdateInput) bool {
				*describeUpdateInput = *input
				return true
			})).Return(awsProvider.Client.MockRequestForGivenOutput(describeUpdateInput, describeUpdateOutput), describeUpdateOutput)
		})

		It("should get current config", func() {
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			enabled.HasAll("api", "audit")
			disabled.HasAll("controllerManager", "authenticator", "scheduler")
		})

		It("should expand `[]` to none", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{}

			api.SetClusterConfigDefaults(cfg)
			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(BeEmpty())

			err = ctl.UpdateClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(sentClusterLogging[0].Types).To(BeEmpty())

			Expect(sentClusterLogging[1].Types).ToNot(BeEmpty())
			Expect(sentClusterLogging[1].Types).To(Equal(aws.StringSlice(api.SupportedCloudWatchClusterLogTypes())))
		})

		It("should expand `['*']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"*"}

			api.SetClusterConfigDefaults(cfg)
			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(api.SupportedCloudWatchClusterLogTypes()))

			err = ctl.UpdateClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(sentClusterLogging[0].Types).ToNot(BeEmpty())
			Expect(sentClusterLogging[0].Types).To(Equal(aws.StringSlice(cfg.CloudWatch.ClusterLogging.EnableTypes)))

			Expect(sentClusterLogging[1].Types).To(BeEmpty())
		})

		It("should enable some logging facilities and disable others", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"authenticator", "controllerManager"}

			api.SetClusterConfigDefaults(cfg)
			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal([]string{"authenticator", "controllerManager"}))

			err = ctl.UpdateClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(sentClusterLogging[0].Types).ToNot(BeEmpty())
			Expect(sentClusterLogging[0].Types).To(Equal(aws.StringSlice(cfg.CloudWatch.ClusterLogging.EnableTypes)))

			Expect(sentClusterLogging[1].Types).ToNot(BeEmpty())
			Expect(sentClusterLogging[1].Types).To(Equal(aws.StringSlice([]string{"api", "audit", "scheduler"})))
		})
	})
})
