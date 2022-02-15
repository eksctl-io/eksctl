package cluster_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

				intialProvider.MockEKS().On("ListClusters", &awseks.ListClustersInput{
					MaxResults: aws.Int64(100),
					Include:    aws.StringSlice([]string{"all"}),
				}).Return(&awseks.ListClustersOutput{
					Clusters: aws.StringSlice([]string{"cluster1", "cluster2", "cluster3"}),
				}, nil)

				stackManager.ListClusterStackNamesReturns(nil, nil)
				stackManager.HasClusterStackUsingCachedListReturnsOnCall(0, true, nil)
				stackManager.HasClusterStackUsingCachedListReturnsOnCall(1, false, nil)
				stackManager.HasClusterStackUsingCachedListReturnsOnCall(2, false, fmt.Errorf("foo"))
			})
			It("returns the clusters in that region", func() {
				clusters, err := cluster.GetClusters(intialProvider, false, 100)
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

				Expect(stackManager.HasClusterStackUsingCachedListCallCount()).To(Equal(3))
				_, clusterName := stackManager.HasClusterStackUsingCachedListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster1"))
				_, clusterName = stackManager.HasClusterStackUsingCachedListArgsForCall(1)
				Expect(clusterName).To(Equal("cluster2"))
				_, clusterName = stackManager.HasClusterStackUsingCachedListArgsForCall(2)
				Expect(clusterName).To(Equal("cluster3"))
			})
		})

		When("ListClusterStackNames errors", func() {
			BeforeEach(func() {
				stackManager.ListClusterStackNamesReturns(nil, fmt.Errorf("foo"))
			})

			It("errors", func() {
				_, err := cluster.GetClusters(intialProvider, false, 100)
				Expect(err).To(MatchError(`failed to list cluster stacks in region "us-west-2": foo`))
			})
		})

		When("ListClusters errors", func() {
			BeforeEach(func() {
				intialProvider.MockEKS().On("ListClusters", &awseks.ListClustersInput{
					MaxResults: aws.Int64(100),
					Include:    aws.StringSlice([]string{"all"}),
				}).Return(nil, fmt.Errorf("foo"))
			})

			It("errors", func() {
				_, err := cluster.GetClusters(intialProvider, false, 100)
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
				awsProvider.ReturnsOnCall(0, &eks.ClusterProvider{Provider: providerRegion1}, nil)
				awsProvider.ReturnsOnCall(1, &eks.ClusterProvider{Provider: providerRegion2}, nil)
				intialProvider.MockEC2().On("DescribeRegions", &awsec2.DescribeRegionsInput{}).Return(&awsec2.DescribeRegionsOutput{
					Regions: []*awsec2.Region{
						{
							RegionName: aws.String("us-west-1"),
						},
						{
							RegionName: aws.String("us-west-2"),
						},
					},
				}, nil)

				providerRegion1.MockEKS().On("ListClusters", &awseks.ListClustersInput{
					MaxResults: aws.Int64(100),
					Include:    aws.StringSlice([]string{"all"}),
				}).Return(&awseks.ListClustersOutput{
					Clusters: aws.StringSlice([]string{"cluster1"}),
				}, nil)

				providerRegion2.MockEKS().On("ListClusters", &awseks.ListClustersInput{
					MaxResults: aws.Int64(100),
					Include:    aws.StringSlice([]string{"all"}),
				}).Return(&awseks.ListClustersOutput{
					Clusters: aws.StringSlice([]string{"cluster2"}),
				}, nil)

				stackManagerRegion1.ListClusterStackNamesReturns(nil, nil)
				stackManagerRegion1.HasClusterStackUsingCachedListReturnsOnCall(0, true, nil)
				stackManagerRegion2.ListClusterStackNamesReturns(nil, nil)
				stackManagerRegion2.HasClusterStackUsingCachedListReturnsOnCall(0, false, nil)
			})

			It("returns the clusters across all authorised regions", func() {
				clusters, err := cluster.GetClusters(intialProvider, true, 100)
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
				cfg, _ := awsProvider.ArgsForCall(0)
				Expect(cfg.Region).To(Equal("us-west-1"))
				cfg, _ = awsProvider.ArgsForCall(1)
				Expect(cfg.Region).To(Equal("us-west-2"))

				Expect(stackManagerRegion1.HasClusterStackUsingCachedListCallCount()).To(Equal(1))
				_, clusterName := stackManagerRegion1.HasClusterStackUsingCachedListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster1"))

				Expect(stackManagerRegion2.HasClusterStackUsingCachedListCallCount()).To(Equal(1))
				_, clusterName = stackManagerRegion2.HasClusterStackUsingCachedListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster2"))
			})
		})

		When("DescribeRegion errors", func() {
			BeforeEach(func() {
				intialProvider.MockEC2().On("DescribeRegions", &awsec2.DescribeRegionsInput{}).Return(nil, fmt.Errorf("foo"))
			})

			It("errors", func() {
				_, err := cluster.GetClusters(intialProvider, true, 100)
				Expect(err).To(MatchError(`failed to describe regions: foo`))
			})
		})

		When("error occurs in a region", func() {
			BeforeEach(func() {
				awsProvider.ReturnsOnCall(0, &eks.ClusterProvider{Provider: providerRegion1}, nil)
				awsProvider.ReturnsOnCall(1, nil, fmt.Errorf("foo"))
				intialProvider.MockEC2().On("DescribeRegions", &awsec2.DescribeRegionsInput{}).Return(&awsec2.DescribeRegionsOutput{
					Regions: []*awsec2.Region{
						{
							RegionName: aws.String("us-west-1"),
						},
						{
							RegionName: aws.String("us-west-2"),
						},
					},
				}, nil)

				providerRegion1.MockEKS().On("ListClusters", &awseks.ListClustersInput{
					MaxResults: aws.Int64(100),
					Include:    aws.StringSlice([]string{"all"}),
				}).Return(&awseks.ListClustersOutput{
					Clusters: aws.StringSlice([]string{"cluster1"}),
				}, nil)

				stackManagerRegion1.ListClusterStackNamesReturns(nil, nil)
				stackManagerRegion1.HasClusterStackUsingCachedListReturnsOnCall(0, true, nil)
			})

			It("returns the clusters in the regions it was successful in", func() {
				clusters, err := cluster.GetClusters(intialProvider, true, 100)
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
				cfg, _ := awsProvider.ArgsForCall(0)
				Expect(cfg.Region).To(Equal("us-west-1"))
				cfg, _ = awsProvider.ArgsForCall(1)
				Expect(cfg.Region).To(Equal("us-west-2"))

				Expect(stackManagerRegion1.HasClusterStackUsingCachedListCallCount()).To(Equal(1))
				_, clusterName := stackManagerRegion1.HasClusterStackUsingCachedListArgsForCall(0)
				Expect(clusterName).To(Equal("cluster1"))
			})
		})
	})
})
