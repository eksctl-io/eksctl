package cluster_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("New", func() {
	DescribeTable("constructing a cluster action", func(clusterExists, stackExists bool, expectedCluster cluster.Cluster, expectedError string) {
		cfg := api.NewClusterConfig()
		cfg.Metadata.Name = "my-cluster"
		p := mockprovider.NewMockProvider()
		ctl := &eks.ClusterProvider{
			AWSProvider: p,
			Status:      &eks.ProviderStatus{},
		}

		describeCluster := p.MockEKS().On("DescribeCluster", mock.Anything, &awseks.DescribeClusterInput{
			Name: aws.String(cfg.Metadata.Name),
		}).Once()
		if clusterExists {
			describeCluster.Return(&awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster(cfg.Metadata.Name, ekstypes.ClusterStatusActive),
			}, nil)
		} else {
			describeCluster.Return(nil, &ekstypes.ResourceNotFoundException{
				Message: aws.String("cluster not found"),
			})
		}

		listStacksOutput := &cloudformation.ListStacksOutput{}
		if stackExists {
			stackName := "eksctl-" + cfg.Metadata.Name + "-cluster"
			listStacksOutput.StackSummaries = []cfntypes.StackSummary{{
				StackName: aws.String(stackName),
			}}
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, &cloudformation.DescribeStacksInput{
				StackName: aws.String(stackName),
			}).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []cfntypes.Stack{{
					StackName:   aws.String(stackName),
					StackStatus: cfntypes.StackStatusDeleteFailed,
					Tags: []cfntypes.Tag{{
						Key:   aws.String(api.ClusterNameTag),
						Value: aws.String(cfg.Metadata.Name),
					}},
				}},
			}, nil).Once()
		}
		p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).
			Return(listStacksOutput, nil).Once()

		actualCluster, err := cluster.New(context.Background(), cfg, ctl)
		if expectedError != "" {
			Expect(err).To(MatchError(expectedError))
			Expect(actualCluster).To(BeNil())
		} else {
			Expect(err).NotTo(HaveOccurred())
			Expect(actualCluster).To(BeAssignableToTypeOf(expectedCluster))
		}
		if !clusterExists {
			Expect(ctl.Status.ClusterInfo).To(BeNil())
		}
		p.MockEKS().AssertExpectations(GinkgoT())
		p.MockCloudFormation().AssertExpectations(GinkgoT())
	},
		Entry("a deleted EKS cluster with a remaining stack is owned", false, true, &cluster.OwnedCluster{}, ""),
		Entry("a deleted EKS cluster without a stack does not exist", false, false, nil, `cluster "my-cluster" does not exist`),
		Entry("an existing EKS cluster with a stack is owned", true, true, &cluster.OwnedCluster{}, ""),
		Entry("an existing EKS cluster without a stack is unowned", true, false, &cluster.UnownedCluster{}, ""),
	)
})
