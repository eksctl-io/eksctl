package cluster

import (
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"k8s.io/client-go/kubernetes/fake"
)

type drainerMock struct {
	mock.Mock
}

func (drainer *drainerMock) Drain(nodeGroups []eks.KubeNodeGroup, plan bool, maxGracePeriod, nodeDrainWaitPeriod time.Duration, undo, disableEviction bool) error {
	args := drainer.Called(nodeGroups, plan, maxGracePeriod, nodeDrainWaitPeriod, undo, disableEviction)
	return args.Error(0)
}

var _ = Describe("Delete", func() {
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
		ctl = &eks.ClusterProvider{Provider: p, Status: &eks.ProviderStatus{}}
	})

	Context("draining node groups", func() {
		It("drain the node groups without disabling the eviction", func() {
			c := NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			nodeGroupStacks := []manager.NodeGroupStack{{NodeGroupName: "ng-1"}}
			kubeNodeGroups := cmdutils.ToKubeNodeGroups(cfg)
			var nodeDrainWaitPeriod time.Duration = 0
			plan := false
			undo := false
			disableEviction := false

			mockedDrainer := &drainerMock{}
			mockedDrainer.On("Drain", kubeNodeGroups, plan, ctl.Provider.WaitTimeout(), nodeDrainWaitPeriod, undo, disableEviction).Return(nil)
			vpcCniDeleterCalled := 0
			vpcCniDeleter := func(clusterName string, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
				vpcCniDeleterCalled++
			}

			err := drainAllNodeGroups(c.cfg, c.ctl, fakeClientSet, nodeGroupStacks, disableEviction, mockedDrainer, vpcCniDeleter)
			Expect(err).NotTo(HaveOccurred())
			mockedDrainer.AssertNumberOfCalls(GinkgoT(), "Drain", 1)
			Expect(vpcCniDeleterCalled).To(Equal(1))
		})

		It("drain the node groups with disabling the eviction", func() {
			c := NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			nodeGroupStacks := []manager.NodeGroupStack{{NodeGroupName: "ng-1"}}
			kubeNodeGroups := cmdutils.ToKubeNodeGroups(cfg)
			var nodeDrainWaitPeriod time.Duration = 0
			plan := false
			undo := false
			disableEviction := true

			mockedDrainer := &drainerMock{}
			mockedDrainer.On("Drain", kubeNodeGroups, plan, ctl.Provider.WaitTimeout(), nodeDrainWaitPeriod, undo, disableEviction).Return(nil)
			vpcCniDeleterCalled := 0
			vpcCniDeleter := func(clusterName string, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
				vpcCniDeleterCalled++
			}

			err := drainAllNodeGroups(c.cfg, c.ctl, fakeClientSet, nodeGroupStacks, disableEviction, mockedDrainer, vpcCniDeleter)
			Expect(err).NotTo(HaveOccurred())
			mockedDrainer.AssertNumberOfCalls(GinkgoT(), "Drain", 1)
			Expect(vpcCniDeleterCalled).To(Equal(1))
		})

		It("does nothing when there are no node group stacks", func() {
			c := NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			nodeGroupStacks := []manager.NodeGroupStack{}
			kubeNodeGroups := cmdutils.ToKubeNodeGroups(cfg)
			var nodeDrainWaitPeriod time.Duration = 0
			plan := false
			undo := false
			disableEviction := false

			mockedDrainer := &drainerMock{}
			mockedDrainer.On("Drain", kubeNodeGroups, plan, ctl.Provider.WaitTimeout(), nodeDrainWaitPeriod, undo, disableEviction).Return(nil)
			vpcCniDeleterCalled := 0
			vpcCniDeleter := func(clusterName string, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
				vpcCniDeleterCalled++
			}

			err := drainAllNodeGroups(c.cfg, c.ctl, fakeClientSet, nodeGroupStacks, disableEviction, mockedDrainer, vpcCniDeleter)
			Expect(err).NotTo(HaveOccurred())
			mockedDrainer.AssertNotCalled(GinkgoT(), "Drain")
			Expect(vpcCniDeleterCalled).To(Equal(0))
		})
	})
})
