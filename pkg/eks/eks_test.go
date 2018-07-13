package eks

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/testutils/mocks"
)

type awsMocks struct {
	CFN *mocks.CloudFormationAPI
	EKS *mocks.EKSAPI
	EC2 *mocks.EC2API
	STS *mocks.STSAPI
}

func newMocks() *awsMocks {
	return &awsMocks{
		CFN: &mocks.CloudFormationAPI{},
		EKS: &mocks.EKSAPI{},
		EC2: &mocks.EC2API{},
		STS: &mocks.STSAPI{},
	}
}

func newClusterProvider(clusterName string, mocks *awsMocks) *ClusterProvider {
	config := &ClusterConfig{
		ClusterName: clusterName,
	}

	cp := &ClusterProvider{
		cfg: config,
		svc: &providerServices{
			cfn: mocks.CFN,
			eks: mocks.EKS,
			ec2: mocks.EC2,
			sts: mocks.STS,
		},
	}

	return cp
}

func TestListAllWithClusterName(t *testing.T) {
	assert := assert.New(t)
	clusterName := "test-cluster"
	clusterStatus := eks.ClusterStatusActive
	logger.Level = 3

	// Setup required mocks
	mocks := newMocks()

	mockDescribeClusterFn := func(input *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
		output := &eks.DescribeClusterOutput{}
		output.Cluster = &eks.Cluster{
			Name:   &clusterName,
			Status: &clusterStatus,
		}
		return output
	}
	mocks.EKS.On("DescribeCluster", mock.MatchedBy(func(input *eks.DescribeClusterInput) bool {
		return input.Name != nil
	})).Return(mockDescribeClusterFn, nil)

	// Get clusterprovider
	cp := newClusterProvider(clusterName, mocks)

	assert.Nil(cp.ListClusters())
	//TODO: Do we need to test the output to stdout

	mocks.EKS.AssertNumberOfCalls(t, "DescribeCluster", 1)
	mocks.CFN.AssertNumberOfCalls(t, "ListStacksPages", 0)
}

func TestListAllWithClusterNameWithVerboseLogging(t *testing.T) {
	assert := assert.New(t)
	clusterName := "test-cluster"
	clusterStatus := eks.ClusterStatusActive
	logger.Level = 4

	// Setup required mocks
	mocks := newMocks()

	mockDescribeClusterFn := func(input *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
		output := &eks.DescribeClusterOutput{}
		output.Cluster = &eks.Cluster{
			Name:   &clusterName,
			Status: &clusterStatus,
		}
		return output
	}
	mocks.EKS.On("DescribeCluster", mock.MatchedBy(func(input *eks.DescribeClusterInput) bool {
		return input.Name != nil
	})).Return(mockDescribeClusterFn, nil)
	mocks.CFN.On("ListStacksPages", mock.Anything, mock.Anything).Return(nil)

	// Get clusterprovider
	cp := newClusterProvider(clusterName, mocks)

	assert.Nil(cp.ListClusters())

	mocks.EKS.AssertNumberOfCalls(t, "DescribeCluster", 1)
	mocks.CFN.AssertNumberOfCalls(t, "ListStacksPages", 1)
}

func TestListAllWithNoClusterName(t *testing.T) {
	assert := assert.New(t)
	clusterName := ""
	clusterStatus := eks.ClusterStatusActive
	logger.Level = 3

	// Setup required mocks
	mocks := newMocks()

	mockListClustersFn := func(inout *eks.ListClustersInput) *eks.ListClustersOutput {
		clusterName1 := "cluster1"
		clusterName2 := "cluster2"
		clusterNames := []*string{&clusterName1, &clusterName2}
		output := &eks.ListClustersOutput{
			Clusters: clusterNames,
		}
		return output
	}
	mocks.EKS.On("ListClusters", mock.MatchedBy(func(input *eks.ListClustersInput) bool {
		return true
	})).Return(mockListClustersFn, nil)

	mockDescribeClusterFn := func(input *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
		output := &eks.DescribeClusterOutput{}
		output.Cluster = &eks.Cluster{
			Name:   input.Name,
			Status: &clusterStatus,
		}
		return output
	}
	mocks.EKS.On("DescribeCluster", mock.MatchedBy(func(input *eks.DescribeClusterInput) bool {
		return input.Name != nil
	})).Return(mockDescribeClusterFn, nil)

	// Get clusterprovider
	cp := newClusterProvider(clusterName, mocks)

	assert.Nil(cp.ListClusters())

	mocks.EKS.AssertNumberOfCalls(t, "ListClusters", 1)
	mocks.EKS.AssertNumberOfCalls(t, "DescribeCluster", 2)
}
