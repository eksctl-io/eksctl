package eks_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("EKS API wrapper", func() {
	toLogTypes := func(logTypes []string) []ekstypes.LogType {
		ret := make([]ekstypes.LogType, len(logTypes))
		for i, lt := range logTypes {
			ret[i] = ekstypes.LogType(lt)
		}
		return ret
	}
	Describe("can update cluster configuration for logging", func() {
		var (
			ctl *ClusterProvider

			cfg *api.ClusterConfig
			err error

			sentClusterLogging []ekstypes.LogSetup

			p *mockprovider.MockProvider
		)

		BeforeEach(func() {
			p = mockprovider.NewMockProvider()
			ctl = &ClusterProvider{
				AWSProvider: p,
				Status:      &ProviderStatus{},
			}

			cfg = api.NewClusterConfig()

			updateClusterConfigOutput := &awseks.UpdateClusterConfigOutput{
				Update: &ekstypes.Update{
					Id:   aws.String("u123"),
					Type: ekstypes.UpdateTypeLoggingUpdate,
				},
			}

			describeClusterOutput := &awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster("testcluster", ekstypes.ClusterStatusActive),
			}

			describeClusterOutput.Cluster.Logging = &ekstypes.Logging{
				ClusterLogging: []ekstypes.LogSetup{
					{
						Enabled: api.Enabled(),
						Types:   []ekstypes.LogType{"api", "audit"},
					},
					{
						Enabled: api.Disabled(),
						Types:   []ekstypes.LogType{"controllerManager", "authenticator", "scheduler"},
					},
				},
			}

			p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				return true
			})).Return(describeClusterOutput, nil)

			p.MockEKS().On("UpdateClusterConfig", mock.Anything, mock.MatchedBy(func(input *awseks.UpdateClusterConfigInput) bool {
				Expect(input.Logging).NotTo(BeNil())

				Expect(input.Logging.ClusterLogging[0].Enabled).NotTo(BeNil())
				Expect(input.Logging.ClusterLogging[1].Enabled).NotTo(BeNil())

				Expect(*input.Logging.ClusterLogging[0].Enabled).To(BeTrue())
				Expect(*input.Logging.ClusterLogging[1].Enabled).To(BeFalse())

				sentClusterLogging = input.Logging.ClusterLogging

				return true
			})).Return(updateClusterConfigOutput, nil)

			describeUpdateInput := &awseks.DescribeUpdateInput{}

			describeUpdateOutput := &awseks.DescribeUpdateOutput{
				Update: &ekstypes.Update{
					Id:     aws.String("u123"),
					Type:   ekstypes.UpdateTypeLoggingUpdate,
					Status: ekstypes.UpdateStatusSuccessful,
				},
			}

			p.MockEKS().On("DescribeUpdate", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeUpdateInput) bool {
				*describeUpdateInput = *input
				return true
			}), mock.Anything).Return(describeUpdateOutput, nil)
		})

		It("should get current config", func() {
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
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

			err = ctl.UpdateClusterConfigForLogging(context.Background(), cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(sentClusterLogging[0].Types).To(BeEmpty())

			Expect(sentClusterLogging[1].Types).NotTo(BeEmpty())
			Expect(sentClusterLogging[1].Types).To(Equal(toLogTypes(api.SupportedCloudWatchClusterLogTypes())))
		})

		It("should expand `['*']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"*"}

			api.SetClusterConfigDefaults(cfg)
			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(api.SupportedCloudWatchClusterLogTypes()))

			err = ctl.UpdateClusterConfigForLogging(context.Background(), cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(sentClusterLogging[0].Types).NotTo(BeEmpty())
			Expect(sentClusterLogging[0].Types).To(Equal(toLogTypes(cfg.CloudWatch.ClusterLogging.EnableTypes)))

			Expect(sentClusterLogging[1].Types).To(BeEmpty())
		})

		It("should enable some logging facilities and disable others", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"authenticator", "controllerManager"}

			api.SetClusterConfigDefaults(cfg)
			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal([]string{"authenticator", "controllerManager"}))

			err = ctl.UpdateClusterConfigForLogging(context.Background(), cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(sentClusterLogging[0].Types).NotTo(BeEmpty())
			Expect(sentClusterLogging[0].Types).To(Equal(toLogTypes(cfg.CloudWatch.ClusterLogging.EnableTypes)))

			Expect(sentClusterLogging[1].Types).NotTo(BeEmpty())
			Expect(sentClusterLogging[1].Types).To(Equal([]ekstypes.LogType{"api", "audit", "scheduler"}))
		})
		It("should update logRetentionInDays in case if it's greater than 0", func() {
			p.MockCloudWatchLogs().On("PutRetentionPolicy", mock.Anything, &cloudwatchlogs.PutRetentionPolicyInput{
				LogGroupName:    aws.String(fmt.Sprintf("/aws/eks/%s/cluster", cfg.Metadata.Name)),
				RetentionInDays: aws.Int32(int32(30)),
			}).Return(&cloudwatchlogs.PutRetentionPolicyOutput{}, nil)
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"authenticator"}
			cfg.CloudWatch.ClusterLogging.LogRetentionInDays = 30

			api.SetClusterConfigDefaults(cfg)
			Expect(api.ValidateClusterConfig(cfg)).To(Succeed())
			Expect(ctl.UpdateClusterConfigForLogging(context.Background(), cfg)).To(Succeed())
		})
	})
})
