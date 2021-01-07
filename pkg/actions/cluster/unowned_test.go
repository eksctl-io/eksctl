package cluster_test

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/weaveworks/eksctl/pkg/utils/strings"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"k8s.io/client-go/kubernetes/fake"

	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	. "github.com/onsi/ginkgo"
	"github.com/weaveworks/eksctl/pkg/actions/cluster"
)

var _ = Describe("Delete", func() {
	It("deletes the cluster", func() {
		clusterName := "my-cluster"
		p := mockprovider.NewMockProvider()
		cfg := api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		clientSet := fake.NewSimpleClientset()

		//mocks are in order of being called
		p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
			Expect(*input.Name).To(Equal(clusterName))
			return true
		})).Return(&awseks.DescribeClusterOutput{
			Cluster: testutils.NewFakeCluster("testcluster", awseks.ClusterStatusActive),
		}, nil)

		p.MockEKS().On("ListFargateProfiles", &awseks.ListFargateProfilesInput{
			ClusterName: strings.Pointer(clusterName),
		}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

		p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil)

		p.MockEC2().On("DescribeKeyPairs", mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

		p.MockEC2().On("DescribeSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

		p.MockEKS().On("ListNodegroupsPages", mock.MatchedBy(func(input *awseks.ListNodegroupsInput) bool {
			Expect(*input.ClusterName).To(Equal(clusterName))
			return true
		}), mock.Anything).Run(func(args mock.Arguments) {
			consume := args[1].(func(ng *awseks.ListNodegroupsOutput, _ bool) bool)
			out := &awseks.ListNodegroupsOutput{
				Nodegroups: aws.StringSlice([]string{"ng-1"}),
			}
			cont := consume(out, true)
			if !cont {
				panic("unexpected return value from the paging function: shouldContinue was false which isn't expected in this test scenario")
			}
		}).Return(nil)

		p.MockEKS().On("ListNodegroups", mock.Anything).Return(&awseks.ListNodegroupsOutput{}, nil)

		p.MockEKS().On("DeleteNodegroup", mock.MatchedBy(func(input *awseks.DeleteNodegroupInput) bool {
			Expect(*input.ClusterName).To(Equal(clusterName))
			Expect(*input.NodegroupName).To(Equal("ng-1"))
			return true
		})).Return(&awseks.DeleteNodegroupOutput{}, nil)

		p.MockEKS().On("DeleteCluster", mock.Anything).Return(&awseks.DeleteClusterOutput{}, nil)

		c := cluster.NewUnownedCluster(cfg, &eks.ClusterProvider{Provider: p, Status: &eks.ProviderStatus{}}, clientSet)
		err := c.Delete(time.Microsecond, false)
		Expect(err).NotTo(HaveOccurred())
	})
})
