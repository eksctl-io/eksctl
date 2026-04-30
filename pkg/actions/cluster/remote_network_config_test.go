package cluster_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/cluster"
	"github.com/weaveworks/eksctl/pkg/actions/cluster/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("UpdateRemoteNetworkConfig", func() {
	var (
		clusterName      string
		p                *mockprovider.MockProvider
		cfg              *api.ClusterConfig
		fakeStackManager *fakes.FakeStackManager
		ctl              *eks.ClusterProvider
		ownedCluster     *cluster.OwnedCluster
	)

	BeforeEach(func() {
		clusterName = "test-cluster"
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		fakeStackManager = new(fakes.FakeStackManager)
		ctl = &eks.ClusterProvider{AWSProvider: p, Status: &eks.ProviderStatus{
			ClusterInfo: &eks.ClusterInfo{
				Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
			},
		}}
		ownedCluster = cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager, &mocks.AutoModeDeleter{})
	})

	mockSuccessfulUpdate := func() {
		p.MockEKS().On("UpdateClusterConfig", mock.Anything, mock.Anything).
			Return(&awseks.UpdateClusterConfigOutput{
				Update: &ekstypes.Update{
					Id:     aws.String("update-123"),
					Status: ekstypes.UpdateStatusSuccessful,
				},
			}, nil)
		p.MockEKS().On("DescribeUpdate", mock.Anything, mock.Anything, mock.Anything).
			Return(&awseks.DescribeUpdateOutput{
				Update: &ekstypes.Update{
					Id:     aws.String("update-123"),
					Status: ekstypes.UpdateStatusSuccessful,
				},
			}, nil)
	}

	It("skips when remoteNetworkConfig is nil", func() {
		cfg.RemoteNetworkConfig = nil
		updated, err := cluster.UpdateRemoteNetworkConfig(ownedCluster, context.Background(), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeFalse())
		p.MockEKS().AssertNotCalled(GinkgoT(), "UpdateClusterConfig", mock.Anything, mock.Anything)
	})

	It("calls UpdateClusterConfig with correct input for node and pod networks", func() {
		cfg.RemoteNetworkConfig = &api.RemoteNetworkConfig{
			RemoteNodeNetworks: []*api.RemoteNetwork{
				{CIDRs: []string{"10.80.0.0/16", "10.81.0.0/16"}},
			},
			RemotePodNetworks: []*api.RemoteNetwork{
				{CIDRs: []string{"10.90.0.0/16"}},
			},
		}
		mockSuccessfulUpdate()

		updated, err := cluster.UpdateRemoteNetworkConfig(ownedCluster, context.Background(), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())

		calls := p.MockEKS().Calls
		var updateCall mock.Call
		for _, call := range calls {
			if call.Method == "UpdateClusterConfig" {
				updateCall = call
				break
			}
		}
		input := updateCall.Arguments[1].(*awseks.UpdateClusterConfigInput)
		Expect(*input.Name).To(Equal(clusterName))
		Expect(input.RemoteNetworkConfig.RemoteNodeNetworks).To(HaveLen(1))
		Expect(input.RemoteNetworkConfig.RemoteNodeNetworks[0].Cidrs).To(Equal([]string{"10.80.0.0/16", "10.81.0.0/16"}))
		Expect(input.RemoteNetworkConfig.RemotePodNetworks).To(HaveLen(1))
		Expect(input.RemoteNetworkConfig.RemotePodNetworks[0].Cidrs).To(Equal([]string{"10.90.0.0/16"}))
	})

	It("passes nil for omitted networks (no defaulting)", func() {
		cfg.RemoteNetworkConfig = &api.RemoteNetworkConfig{
			RemoteNodeNetworks: []*api.RemoteNetwork{
				{CIDRs: []string{"10.80.0.0/16"}},
			},
		}
		mockSuccessfulUpdate()

		updated, err := cluster.UpdateRemoteNetworkConfig(ownedCluster, context.Background(), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())

		calls := p.MockEKS().Calls
		var updateCall mock.Call
		for _, call := range calls {
			if call.Method == "UpdateClusterConfig" {
				updateCall = call
				break
			}
		}
		input := updateCall.Arguments[1].(*awseks.UpdateClusterConfigInput)
		Expect(input.RemoteNetworkConfig.RemoteNodeNetworks).To(HaveLen(1))
		Expect(input.RemoteNetworkConfig.RemotePodNetworks).To(BeNil())
	})

	It("sends empty lists for remove-all case", func() {
		cfg.RemoteNetworkConfig = &api.RemoteNetworkConfig{
			RemoteNodeNetworks: []*api.RemoteNetwork{},
			RemotePodNetworks:  []*api.RemoteNetwork{},
		}
		mockSuccessfulUpdate()

		updated, err := cluster.UpdateRemoteNetworkConfig(ownedCluster, context.Background(), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())

		calls := p.MockEKS().Calls
		var updateCall mock.Call
		for _, call := range calls {
			if call.Method == "UpdateClusterConfig" {
				updateCall = call
				break
			}
		}
		input := updateCall.Arguments[1].(*awseks.UpdateClusterConfigInput)
		Expect(input.RemoteNetworkConfig.RemoteNodeNetworks).To(BeEmpty())
		Expect(input.RemoteNetworkConfig.RemoteNodeNetworks).NotTo(BeNil())
		Expect(input.RemoteNetworkConfig.RemotePodNetworks).To(BeEmpty())
		Expect(input.RemoteNetworkConfig.RemotePodNetworks).NotTo(BeNil())
	})

	It("treats 'No changes detected' as success", func() {
		cfg.RemoteNetworkConfig = &api.RemoteNetworkConfig{
			RemoteNodeNetworks: []*api.RemoteNetwork{
				{CIDRs: []string{"10.80.0.0/16"}},
			},
		}
		p.MockEKS().On("UpdateClusterConfig", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("operation error EKS: UpdateClusterConfig, https response error StatusCode: 400, InvalidParameterException: No changes detected for remoteNetworkConfig"))

		updated, err := cluster.UpdateRemoteNetworkConfig(ownedCluster, context.Background(), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeFalse())
	})

	It("propagates other API errors", func() {
		cfg.RemoteNetworkConfig = &api.RemoteNetworkConfig{
			RemoteNodeNetworks: []*api.RemoteNetwork{
				{CIDRs: []string{"10.80.0.0/16"}},
			},
		}
		p.MockEKS().On("UpdateClusterConfig", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("operation error EKS: UpdateClusterConfig, https response error StatusCode: 400, InvalidParameterException: Only one remoteNodeNetwork is allowed"))

		_, err := cluster.UpdateRemoteNetworkConfig(ownedCluster, context.Background(), false)
		Expect(err).To(MatchError(ContainSubstring("Only one remoteNodeNetwork is allowed")))
	})

	It("does not call API in dry-run mode", func() {
		cfg.RemoteNetworkConfig = &api.RemoteNetworkConfig{
			RemoteNodeNetworks: []*api.RemoteNetwork{
				{CIDRs: []string{"10.80.0.0/16"}},
			},
		}

		updated, err := cluster.UpdateRemoteNetworkConfig(ownedCluster, context.Background(), true)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())
		p.MockEKS().AssertNotCalled(GinkgoT(), "UpdateClusterConfig", mock.Anything, mock.Anything)
	})
})
