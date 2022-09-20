package cluster_test

import (
	"context"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type drainerMock struct {
	mock.Mock
}

func (d *drainerMock) Drain(_ context.Context, input *nodegroup.DrainInput) error {
	args := d.Called(input)
	return args.Error(0)
}

var _ = Describe("DrainAllNodeGroups", func() {
	var (
		clusterName      string
		p                *mockprovider.MockProvider
		cfg              *api.ClusterConfig
		fakeStackManager *fakes.FakeStackManager
		fakeClientSet    *fake.Clientset
		ctl              *eks.ClusterProvider
	)

	BeforeEach(func() {
		clusterName = "my-cluster"
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		fakeStackManager = new(fakes.FakeStackManager)
		fakeClientSet = fake.NewSimpleClientset()
		ctl = &eks.ClusterProvider{AWSProvider: p, Status: &eks.ProviderStatus{}}
	})

	Context("draining node groups", func() {
		When("disable eviction flag is set to false", func() {
			It("drain the node groups", func() {
				c := cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				nodeGroupStacks := []manager.NodeGroupStack{{NodeGroupName: "ng-1"}}
				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.AWSProvider.WaitTimeout(),
					Parallel:       1,
				}

				mockedDrainer := &drainerMock{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(nil)
				vpcCniDeleterCalled := 0
				vpcCniDeleter := func(_ *api.ClusterConfig, _ *eks.ClusterProvider, _ kubernetes.Interface) {
					vpcCniDeleterCalled++
				}

				err := cluster.DrainAllNodeGroups(context.Background(), cfg, ctl, fakeClientSet, nodeGroupStacks, false, 1, mockedDrainer, vpcCniDeleter, 0)
				Expect(err).NotTo(HaveOccurred())
				mockedDrainer.AssertNumberOfCalls(GinkgoT(), "Drain", 1)
				Expect(vpcCniDeleterCalled).To(Equal(1))
			})
		})

		When("disable eviction flag is set to true", func() {
			It("drain the node groups", func() {
				c := cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				nodeGroupStacks := []manager.NodeGroupStack{{NodeGroupName: "ng-1"}}
				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:      cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod:  ctl.AWSProvider.WaitTimeout(),
					DisableEviction: true,
					Parallel:        1,
				}

				mockedDrainer := &drainerMock{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(nil)
				vpcCniDeleterCalled := 0
				vpcCniDeleter := func(_ *api.ClusterConfig, _ *eks.ClusterProvider, _ kubernetes.Interface) {
					vpcCniDeleterCalled++
				}

				err := cluster.DrainAllNodeGroups(context.Background(), cfg, ctl, fakeClientSet, nodeGroupStacks, true, 1, mockedDrainer, vpcCniDeleter, 0)
				Expect(err).NotTo(HaveOccurred())
				mockedDrainer.AssertNumberOfCalls(GinkgoT(), "Drain", 1)
				Expect(vpcCniDeleterCalled).To(Equal(1))
			})
		})

		When("no node group stacks exist", func() {
			It("does no draining at all", func() {
				c := cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				var nodeGroupStacks []manager.NodeGroupStack
				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.AWSProvider.WaitTimeout(),
					Parallel:       1,
				}

				mockedDrainer := &drainerMock{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(nil)
				vpcCniDeleterCalled := 0
				vpcCniDeleter := func(_ *api.ClusterConfig, _ *eks.ClusterProvider, _ kubernetes.Interface) {
					vpcCniDeleterCalled++
				}

				err := cluster.DrainAllNodeGroups(context.Background(), cfg, ctl, fakeClientSet, nodeGroupStacks, false, 1, mockedDrainer, vpcCniDeleter, 0)
				Expect(err).NotTo(HaveOccurred())
				mockedDrainer.AssertNotCalled(GinkgoT(), "Drain")
				Expect(vpcCniDeleterCalled).To(Equal(0))
			})
		})
	})
})
