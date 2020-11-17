package cluster_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cluster"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Owned Clusters", func() {
	Describe("upgrade", func() {
		var (
			ownedCluster *cluster.OwnedCluster
			p            *mockprovider.MockProvider
			err          error
		)
		BeforeEach(func() {
			p = mockprovider.NewMockProvider()
			cfg := api.NewClusterConfig()
			cfg.Metadata.Name = "owned"
			stackManager := manager.NewStackCollection(p, &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Name: "owned",
				},
			})
			ctl := &eks.ClusterProvider{
				Provider: p,
				Status: &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: testutils.NewFakeCluster("owned", "idc"),
					},
				},
			}

			p.MockEKS().On("DescribeCluster", mock.Anything).Return(&awseks.DescribeClusterOutput{
				Cluster: &awseks.Cluster{
					ResourcesVpcConfig: &awseks.VpcConfigResponse{
						EndpointPrivateAccess: aws.Bool(false),
						EndpointPublicAccess:  aws.Bool(true),
					},
					Status: aws.String(awseks.ClusterStatusActive),
					CertificateAuthority: &awseks.Certificate{
						Data: aws.String("Zm9vCg=="),
					},
					Endpoint: aws.String("foo"),
					Arn:      aws.String("123"),
					Version:  aws.String(api.Version1_17),
				},
			}, nil)

			p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				pager := args.Get(1).(func(*cloudformation.ListStacksOutput, bool) bool)
				pager(&cloudformation.ListStacksOutput{
					StackSummaries: []*cloudformation.StackSummary{
						{
							StackName: aws.String("eksctl-owned-cluster"),
							StackId:   aws.String("1"),
						},
					},
				}, true)
			}).Return(nil)

			p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					{
						StackName:   aws.String("eksctl-owned-cluster"),
						StackId:     aws.String("1"),
						StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
						Tags: []*cloudformation.Tag{
							{
								Key:   aws.String(api.ClusterNameTag),
								Value: aws.String("owned"),
							},
						},
						Outputs: []*cloudformation.Output{
							{
								OutputKey:   aws.String("VPC"),
								OutputValue: aws.String("my-vpc"),
							},
							{
								OutputKey:   aws.String("SecurityGroup"),
								OutputValue: aws.String("my-sg"),
							},
						},
					},
				},
			}, nil)

			ownedCluster, err = cluster.NewOwnedCluster(cfg, ctl, stackManager)
			Expect(err).NotTo(HaveOccurred())
		})

		It("upgrades the cluster", func() {
			err := ownedCluster.Upgrade(false)

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
