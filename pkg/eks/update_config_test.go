package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("EKS API wrapper", func() {
	Describe("can update cluster configuration for logging", func() {
		var (
			ctl *ClusterProvider

			cfg *api.ClusterConfig
			err error

			sentClusterLogging []*awseks.LogSetup
		)

		BeforeEach(func() {
			p := mockprovider.NewMockProvider()
			ctl = &ClusterProvider{
				Provider: p,
			}

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

			p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				return true
			})).Return(describeClusterOutput, nil)

			p.MockEKS().On("UpdateClusterConfig", mock.MatchedBy(func(input *awseks.UpdateClusterConfigInput) bool {
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

			p.MockEKS().On("DescribeUpdateRequest", mock.MatchedBy(func(input *awseks.DescribeUpdateInput) bool {
				*describeUpdateInput = *input
				return true
			})).Return(p.Client.MockRequestForGivenOutput(describeUpdateInput, describeUpdateOutput), describeUpdateOutput)
		})

		It("should get current config", func() {
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(cfg.Metadata)
			Expect(err).NotTo(HaveOccurred())

			enabled.HasAll("api", "audit")
			disabled.HasAll("controllerManager", "authenticator", "scheduler")
		})

		It("should expand `[]` to none", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{}

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
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

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
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

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
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
