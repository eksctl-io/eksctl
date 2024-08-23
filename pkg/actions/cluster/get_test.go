package cluster_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/cluster"
	"github.com/weaveworks/eksctl/pkg/actions/cluster/fakes"
	mgrfakes "github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Get", func() {
	var (
		awsProvider             *fakes.FakeProviderConstructor
		stackCollectionProvider *fakes.FakeStackManagerConstructor
		intialProvider          *mockprovider.MockProvider
	)

	BeforeEach(func() {
		intialProvider = mockprovider.NewMockProvider()
		intialProvider.SetRegion("us-west-2")
		awsProvider = new(fakes.FakeProviderConstructor)
		stackCollectionProvider = new(fakes.FakeStackManagerConstructor)
		cluster.SetProviderConstructor(awsProvider.Spy)
		cluster.SetStackManagerConstructor(stackCollectionProvider.Spy)
	})

	When("listing in a specific region", func() {
		var stackManager *mgrfakes.FakeStackManager
		BeforeEach(func() {
			stackManager = new(mgrfakes.FakeStackManager)
			stackCollectionProvider.Returns(stackManager)
		})

		When("it succeeds", func() {
			BeforeEach(func() {

				intialProvider.MockEKS().On("ListClusters", mock.Anything, &awseks.ListClustersInput{
					MaxResults: aws.Int32(100),
					Include:    []string{"all"},
				}, mock.Anything).Return(&awseks.ListClustersOutput{
					Clusters: []string{"cluster1", "cluster2", "cluster3"},
				}, nil)

				stackManager.ListClusterStackNamesReturns(nil, nil)
				stackManager.HasClusterStackFromListReturnsOnCall(0, true, nil)
				stackManager.HasClusterStackFromListReturnsOnCall(1, false, nil)
				stackManager.HasClusterStackFromListReturnsOnCall(2, false, fmt.Errorf("foo"))
			})
			It("returns the clusters in that region", func() {
				clusters, err := cluster.GetClusters(context.Background(), intialProvider, false, 100)
				Expect(err).NotTo(HaveOccurred())
				Expect(clusters).To(ConsistOf(
					cluster.Description{
						Name:   "cluster1",
						Region: "us-west-2",
						Owned:  "True",
					},
					cluster.Description{
						Name:   "cluster2",
						Region: "us-west-2",
						Owned:  "False",
					},
					cluster.Description{
						Name:   "cluster3",
						Region: "us-west-2",
						Owned:  "Unknown",
					},
				))

				Expect(stackCollectionProvider.CallCount()).To(Equal(1))
				provider, _ := stackCollectionProvider.ArgsForCall(0)
				Expect(provider).To(Equal(intialProvider))
				Expect(awsProvider.CallCount()).To(Equal(0))

				Expect(stackManager.HasClusterStackFromListCallCount()).To(Equal(3))
				_, _, clusterName := stackManager.HasClusterStackFromListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster1"))
				_, _, clusterName = stackManager.HasClusterStackFromListArgsForCall(1)
				Expect(clusterName).To(Equal("cluster2"))
				_, _, clusterName = stackManager.HasClusterStackFromListArgsForCall(2)
				Expect(clusterName).To(Equal("cluster3"))
			})
		})

		When("ListClusterStackNames errors", func() {
			BeforeEach(func() {
				stackManager.ListClusterStackNamesReturns(nil, fmt.Errorf("foo"))
			})

			It("errors", func() {
				_, err := cluster.GetClusters(context.Background(), intialProvider, false, 100)
				Expect(err).To(MatchError(`failed to list cluster stacks in region "us-west-2": foo`))
			})
		})

		When("ListClusters errors", func() {
			BeforeEach(func() {
				intialProvider.MockEKS().On("ListClusters", mock.Anything, &awseks.ListClustersInput{
					MaxResults: aws.Int32(100),
					Include:    []string{"all"},
				}, mock.Anything).Return(nil, fmt.Errorf("foo"))
			})

			It("errors", func() {
				_, err := cluster.GetClusters(context.Background(), intialProvider, false, 100)
				Expect(err).To(MatchError(`failed to list clusters in region "us-west-2": foo`))
			})
		})
	})

	When("listing all regions", func() {
		var (
			providerRegion1     *mockprovider.MockProvider
			providerRegion2     *mockprovider.MockProvider
			stackManagerRegion1 *mgrfakes.FakeStackManager
			stackManagerRegion2 *mgrfakes.FakeStackManager
		)

		BeforeEach(func() {
			providerRegion1 = mockprovider.NewMockProvider()
			providerRegion1.SetRegion("us-west-1")
			providerRegion2 = mockprovider.NewMockProvider()
			providerRegion2.SetRegion("us-west-2")

			stackManagerRegion1 = new(mgrfakes.FakeStackManager)
			stackCollectionProvider.ReturnsOnCall(0, stackManagerRegion1)
			stackManagerRegion2 = new(mgrfakes.FakeStackManager)
			stackCollectionProvider.ReturnsOnCall(1, stackManagerRegion2)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				awsProvider.ReturnsOnCall(0, &eks.ClusterProvider{AWSProvider: providerRegion1}, nil)
				awsProvider.ReturnsOnCall(1, &eks.ClusterProvider{AWSProvider: providerRegion2}, nil)
				intialProvider.MockEC2().On("DescribeRegions", mock.Anything, &ec2.DescribeRegionsInput{}).Return(&ec2.DescribeRegionsOutput{
					Regions: []ec2types.Region{
						{
							RegionName: aws.String("us-west-1"),
						},
						{
							RegionName: aws.String("us-west-2"),
						},
					},
				}, nil)

				providerRegion1.MockEKS().On("ListClusters", mock.Anything, &awseks.ListClustersInput{
					MaxResults: aws.Int32(100),
					Include:    []string{"all"},
				}, mock.Anything).Return(&awseks.ListClustersOutput{
					Clusters: []string{"cluster1"},
				}, nil)

				providerRegion2.MockEKS().On("ListClusters", mock.Anything, &awseks.ListClustersInput{
					MaxResults: aws.Int32(100),
					Include:    []string{"all"},
				}, mock.Anything).Return(&awseks.ListClustersOutput{
					Clusters: []string{"cluster2"},
				}, nil)

				stackManagerRegion1.ListClusterStackNamesReturns(nil, nil)
				stackManagerRegion1.HasClusterStackFromListReturnsOnCall(0, true, nil)
				stackManagerRegion2.ListClusterStackNamesReturns(nil, nil)
				stackManagerRegion2.HasClusterStackFromListReturnsOnCall(0, false, nil)
			})

			It("returns the clusters across all authorised regions", func() {
				clusters, err := cluster.GetClusters(context.Background(), intialProvider, true, 100)
				Expect(err).NotTo(HaveOccurred())
				Expect(clusters).To(ConsistOf(
					cluster.Description{
						Name:   "cluster1",
						Region: "us-west-1",
						Owned:  "True",
					},
					cluster.Description{
						Name:   "cluster2",
						Region: "us-west-2",
						Owned:  "False",
					},
				))

				Expect(stackCollectionProvider.CallCount()).To(Equal(2))
				provider, _ := stackCollectionProvider.ArgsForCall(0)
				Expect(provider).To(Equal(providerRegion1))
				provider, _ = stackCollectionProvider.ArgsForCall(1)
				Expect(provider).To(Equal(providerRegion2))

				Expect(awsProvider.CallCount()).To(Equal(2))
				_, cfg, _ := awsProvider.ArgsForCall(0)
				Expect(cfg.Region).To(Equal("us-west-1"))
				_, cfg, _ = awsProvider.ArgsForCall(1)
				Expect(cfg.Region).To(Equal("us-west-2"))

				Expect(stackManagerRegion1.HasClusterStackFromListCallCount()).To(Equal(1))
				_, _, clusterName := stackManagerRegion1.HasClusterStackFromListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster1"))

				Expect(stackManagerRegion2.HasClusterStackFromListCallCount()).To(Equal(1))
				_, _, clusterName = stackManagerRegion2.HasClusterStackFromListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster2"))
			})
		})

		When("DescribeRegion errors", func() {
			BeforeEach(func() {
				intialProvider.MockEC2().On("DescribeRegions", mock.Anything, &ec2.DescribeRegionsInput{}).Return(nil, fmt.Errorf("foo"))
			})

			It("errors", func() {
				_, err := cluster.GetClusters(context.Background(), intialProvider, true, 100)
				Expect(err).To(MatchError(`failed to describe regions: foo`))
			})
		})

		When("error occurs in a region", func() {
			BeforeEach(func() {
				awsProvider.ReturnsOnCall(0, &eks.ClusterProvider{AWSProvider: providerRegion1}, nil)
				awsProvider.ReturnsOnCall(1, nil, fmt.Errorf("foo"))
				intialProvider.MockEC2().On("DescribeRegions", mock.Anything, &ec2.DescribeRegionsInput{}).Return(&ec2.DescribeRegionsOutput{
					Regions: []ec2types.Region{
						{
							RegionName: aws.String("us-west-1"),
						},
						{
							RegionName: aws.String("us-west-2"),
						},
					},
				}, nil)

				providerRegion1.MockEKS().On("ListClusters", mock.Anything, &awseks.ListClustersInput{
					MaxResults: aws.Int32(100),
					Include:    []string{"all"},
				}, mock.Anything).Return(&awseks.ListClustersOutput{
					Clusters: []string{"cluster1"},
				}, nil)

				stackManagerRegion1.ListClusterStackNamesReturns(nil, nil)
				stackManagerRegion1.HasClusterStackFromListReturnsOnCall(0, true, nil)
			})

			It("returns the clusters in the regions it was successful in", func() {
				clusters, err := cluster.GetClusters(context.Background(), intialProvider, true, 100)
				Expect(err).NotTo(HaveOccurred())
				Expect(clusters).To(ConsistOf(
					cluster.Description{
						Name:   "cluster1",
						Region: "us-west-1",
						Owned:  "True",
					},
				))

				Expect(stackCollectionProvider.CallCount()).To(Equal(1))
				provider, _ := stackCollectionProvider.ArgsForCall(0)
				Expect(provider).To(Equal(providerRegion1))

				Expect(awsProvider.CallCount()).To(Equal(2))
				_, cfg, _ := awsProvider.ArgsForCall(0)
				Expect(cfg.Region).To(Equal("us-west-1"))
				_, cfg, _ = awsProvider.ArgsForCall(1)
				Expect(cfg.Region).To(Equal("us-west-2"))

				Expect(stackManagerRegion1.HasClusterStackFromListCallCount()).To(Equal(1))
				_, _, clusterName := stackManagerRegion1.HasClusterStackFromListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster1"))
			})
		})
	})
})
