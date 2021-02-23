package drain_test

import (
	"context"
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
		fakeEvictor   *fakes.FakeEvictor
		nodeName      = "node-1"
	)

	BeforeEach(func() {
		mockNG = mocks.KubeNodeGroup{}
		mockNG.Mock.On("NameString").Return("node-1")
		mockNG.Mock.On("ListOptions").Return(metav1.ListOptions{})
		fakeClientSet = fake.NewSimpleClientset()
		fakeEvictor = new(fakes.FakeEvictor)
	})

	When("all nodes drain successfully", func() {
		var pod corev1.Pod

		BeforeEach(func() {
			pod = corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-1",
				},
			}

			fakeEvictor.GetPodsForEvictionReturnsOnCall(0, &evictor.PodDeleteList{
				Items: []evictor.PodDelete{
					{
						Pod: pod,
						Status: evictor.PodDeleteStatus{
							Delete: true,
						},
					},
				},
			}, nil)
			fakeEvictor.GetPodsForEvictionReturnsOnCall(1, &evictor.PodDeleteList{}, nil)

			fakeEvictor.EvictOrDeletePodReturns(nil)

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: corev1.NodeSpec{
					Unschedulable: false,
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not error", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second*10, time.Second, false, false)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(Equal(2))
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(Equal(1))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(0)).To(Equal(pod))

			node, err := fakeClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Spec.Unschedulable).To(BeTrue())
		})
	})

	When("the nodes never drain successfully", func() {
		var pod corev1.Pod

		BeforeEach(func() {
			pod = corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-1",
				},
			}

			fakeEvictor.GetPodsForEvictionReturns(&evictor.PodDeleteList{
				Items: []evictor.PodDelete{
					{
						Pod: pod,
						Status: evictor.PodDeleteStatus{
							Delete: true,
						},
					},
				},
			}, nil)

			fakeEvictor.EvictOrDeletePodReturns(nil)

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("times out and errors", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second*2, time.Second, false, false)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain()
			Expect(err).To(MatchError("timed out (after 2s) waiting for nodegroup \"node-1\" to be drained"))
		})
	})

	When("Evictions are not supported", func() {
		BeforeEach(func() {
			fakeEvictor.CanUseEvictionsReturns(fmt.Errorf("error1"))
		})

		It("errors", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, false, false)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain()
			Expect(err).To(MatchError("checking if cluster implements policy API: error1"))
		})
	})

	When("Evictions are disabled", func() {
		var pod corev1.Pod

		BeforeEach(func() {
			pod = corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-1",
				},
			}

			fakeEvictor.GetPodsForEvictionReturnsOnCall(0, &evictor.PodDeleteList{
				Items: []evictor.PodDelete{
					{
						Pod: pod,
						Status: evictor.PodDeleteStatus{
							Delete: true,
						},
					},
				},
			}, nil)
			fakeEvictor.GetPodsForEvictionReturnsOnCall(1, &evictor.PodDeleteList{}, nil)

			fakeEvictor.EvictOrDeletePodReturns(nil)

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: corev1.NodeSpec{
					Unschedulable: false,
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not error", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second*10, time.Second, false, true)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(Equal(2))
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(Equal(1))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(0)).To(Equal(pod))

			node, err := fakeClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Spec.Unschedulable).To(BeTrue())
		})
	})

	When("undo is true", func() {
		BeforeEach(func() {
			_, err := fakeClientSet.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: corev1.NodeSpec{
					Unschedulable: true,
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("uncordons all the nodes", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second*10, time.Second, true, false)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain()
			Expect(err).NotTo(HaveOccurred())

			node, err := fakeClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Spec.Unschedulable).To(BeFalse())

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(BeZero())
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(BeZero())
		})
	})
})
