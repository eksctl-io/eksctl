package drain_test

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/weaveworks/eksctl/pkg/drain/evictor"

	"github.com/weaveworks/eksctl/pkg/drain/fakes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Drain", func() {
	var (
		mockNG        mocks.KubeNodeGroup
		fakeClientSet *fake.Clientset
		fakeDrainer   *fakes.FakeDrainer
	)

	BeforeEach(func() {
		mockNG = mocks.KubeNodeGroup{}
		mockNG.Mock.On("NameString").Return("node-1")
		mockNG.Mock.On("ListOptions").Return(metav1.ListOptions{})
		fakeClientSet = fake.NewSimpleClientset()
		fakeDrainer = new(fakes.FakeDrainer)
	})

	When("all nodes drain successfully", func() {
		var pod corev1.Pod

		BeforeEach(func() {
			pod = corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-1",
				},
			}

			fakeDrainer.GetPodsForDeletionReturnsOnCall(0, &evictor.PodDeleteList{
				Items: []evictor.PodDelete{
					{
						Pod: pod,
						Status: evictor.PodDeleteStatus{
							Delete: true,
						},
					},
				},
			}, nil)
			fakeDrainer.GetPodsForDeletionReturnsOnCall(1, &evictor.PodDeleteList{}, nil)

			fakeDrainer.EvictOrDeletePodReturns(nil)

			fakeClientSet.CoreV1().Nodes().Create(&corev1.Node{})
		})

		It("does not error", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second*10, time.Second, false)
			nodeGroupDrainer.SetDrainer(fakeDrainer)

			err := nodeGroupDrainer.Drain()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeDrainer.GetPodsForDeletionCallCount()).To(Equal(2))
			Expect(fakeDrainer.EvictOrDeletePodCallCount()).To(Equal(1))
			Expect(fakeDrainer.EvictOrDeletePodArgsForCall(0)).To(Equal(pod))
		})
	})

	When("Evictions are not supported", func() {
		BeforeEach(func() {
			fakeDrainer.CanUseEvictionsReturns(fmt.Errorf("error1"))
		})

		It("errors", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, false)
			nodeGroupDrainer.SetDrainer(fakeDrainer)

			err := nodeGroupDrainer.Drain()
			Expect(err).To(MatchError("checking if cluster implements policy API: error1"))
		})
	})
})
