package cluster

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"time"

	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

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
			disableEviction := false

			nodeGroupDrainerCalled := 0
			nodeGroupDrainer := func(nodeGroups []eks.KubeNodeGroup, plan bool, maxGracePeriod time.Duration, disableEviction bool) error {
				nodeGroupDrainerCalled++
				return nil
			}
			vpcCniDeleterCalled := 0
			vpcCniDeleter := func(clusterName string, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
				vpcCniDeleterCalled++
			}

			err := drainAllNodeGroups(c.cfg, c.ctl, fakeClientSet, nodeGroupStacks, &disableEviction, nodeGroupDrainer, vpcCniDeleter)
			Expect(err).NotTo(HaveOccurred())
			Expect(nodeGroupDrainerCalled).To(Equal(1))
			Expect(vpcCniDeleterCalled).To(Equal(1))
		})

		It("drain the node groups with disabling the eviction", func() {
			c := NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			nodeGroupStacks := []manager.NodeGroupStack{{NodeGroupName: "ng-1"}}
			disableEviction := true

			nodeGroupDrainerCalled := 0
			nodeGroupDrainer := func(nodeGroups []eks.KubeNodeGroup, plan bool, maxGracePeriod time.Duration, disableEviction bool) error {
				nodeGroupDrainerCalled++
				return nil
			}
			vpcCniDeleterCalled := 0
			vpcCniDeleter := func(clusterName string, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
				vpcCniDeleterCalled++
			}

			err := drainAllNodeGroups(c.cfg, c.ctl, fakeClientSet, nodeGroupStacks, &disableEviction, nodeGroupDrainer, vpcCniDeleter)
			Expect(err).NotTo(HaveOccurred())
			Expect(nodeGroupDrainerCalled).To(Equal(1))
			Expect(vpcCniDeleterCalled).To(Equal(1))
		})

		It("does nothing when there are no node group stacks", func() {
			c := NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			nodeGroupStacks := []manager.NodeGroupStack{}
			disableEviction := false

			nodeGroupDrainerCalled := 0
			nodeGroupDrainer := func(nodeGroups []eks.KubeNodeGroup, plan bool, maxGracePeriod time.Duration, disableEviction bool) error {
				nodeGroupDrainerCalled++
				return nil
			}
			vpcCniDeleterCalled := 0
			vpcCniDeleter := func(clusterName string, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
				vpcCniDeleterCalled++
			}

			err := drainAllNodeGroups(c.cfg, c.ctl, fakeClientSet, nodeGroupStacks, &disableEviction, nodeGroupDrainer, vpcCniDeleter)
			Expect(err).NotTo(HaveOccurred())
			Expect(nodeGroupDrainerCalled).To(Equal(0))
			Expect(vpcCniDeleterCalled).To(Equal(0))
		})
	})
})
