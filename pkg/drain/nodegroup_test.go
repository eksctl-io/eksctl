package drain_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/weaveworks/eksctl/pkg/drain/evictor"

	"github.com/weaveworks/eksctl/pkg/drain/fakes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
)

var _ = Describe("Drain", func() {
	var (
		mockNG        mocks.KubeNodeGroup
		fakeClientSet *fake.Clientset
		fakeEvictor   *fakes.FakeEvictor
		nodeName      = "node-1"
		ctx           context.Context
	)

	BeforeEach(func() {
		mockNG = mocks.KubeNodeGroup{}
		mockNG.Mock.On("NameString").Return("node-1")
		mockNG.Mock.On("ListOptions").Return(metav1.ListOptions{})
		fakeClientSet = fake.NewSimpleClientset()
		fakeEvictor = new(fakes.FakeEvictor)
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		DeferCleanup(cancel)
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

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{
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
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second*10, time.Second, 0, false, false, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(BeNumerically(">=", 2))
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(Equal(1))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(0)).To(Equal(pod))

			node, err := fakeClientSet.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
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

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("times out and errors", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, 0, time.Second, 0, false, false, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err := nodeGroupDrainer.Drain(ctx)
			Expect(err).To(MatchError(`timed out waiting for nodegroup "node-1" to be drained`))
		})
	})

	When("Evictions are not supported", func() {
		BeforeEach(func() {
			fakeEvictor.CanUseEvictionsReturns(fmt.Errorf("error1"))
		})

		It("errors", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, 0, false, false, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := nodeGroupDrainer.Drain(ctx)
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

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{
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
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, time.Second*0, false, true, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			Expect(nodeGroupDrainer.Drain(ctx)).To(Succeed())

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(BeNumerically(">=", 2))
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(Equal(1))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(0)).To(Equal(pod))

			node, err := fakeClientSet.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Spec.Unschedulable).To(BeTrue())
		})
	})

	When("undo is true", func() {
		BeforeEach(func() {
			_, err := fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{
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
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, time.Second*0, true, false, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain(ctx)
			Expect(err).NotTo(HaveOccurred())

			node, err := fakeClientSet.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(node.Spec.Unschedulable).To(BeFalse())

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(BeZero())
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(BeZero())
		})
	})

	When("an eviction fails recoverably", func() {
		var pod corev1.Pod

		BeforeEach(func() {
			pod = corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-1",
				},
			}

			for i := 0; i < 2; i++ {
				fakeEvictor.GetPodsForEvictionReturnsOnCall(i, &evictor.PodDeleteList{
					Items: []evictor.PodDelete{
						{
							Pod: pod,
							Status: evictor.PodDeleteStatus{
								Delete: true,
							},
						},
					},
				}, nil)
			}
			fakeEvictor.GetPodsForEvictionReturnsOnCall(2, nil, nil)

			fakeEvictor.EvictOrDeletePodReturnsOnCall(0, apierrors.NewTooManyRequestsError("error1"))
			fakeEvictor.EvictOrDeletePodReturnsOnCall(1, nil)

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{
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
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, 0, false, false, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			Expect(nodeGroupDrainer.Drain(ctx)).To(Succeed())

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(Equal(3))
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(Equal(2))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(0)).To(Equal(pod))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(1)).To(Equal(pod))
		})
	})

	When("an eviction fails irrecoverably", func() {
		var pod corev1.Pod
		var evictionError error

		BeforeEach(func() {
			pod = corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "ns-1",
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

			evictionError = errors.New("error1")
			fakeEvictor.EvictOrDeletePodReturns(evictionError)

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: corev1.NodeSpec{
					Unschedulable: false,
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, 0, false, false, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			err := nodeGroupDrainer.Drain(ctx)
			Expect(err.Error()).To(ContainSubstring("unrecoverable error evicting pod: ns-1/pod-1"))

			Expect(fakeEvictor.GetPodsForEvictionCallCount()).To(Equal(1))
			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(Equal(1))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(0)).To(Equal(pod))
		})
	})

	When("eviction fails recoverably with multiple pods", func() {
		var pods []corev1.Pod
		var evictionError error

		BeforeEach(func() {
			pods = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-1",
						Namespace: "ns-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-2",
						Namespace: "ns-2",
					},
				},
			}

			fakeEvictor.GetPodsForEvictionReturns(&evictor.PodDeleteList{
				Items: []evictor.PodDelete{
					{
						Pod: pods[0],
						Status: evictor.PodDeleteStatus{
							Delete: true,
						},
					},
					{
						Pod: pods[1],
						Status: evictor.PodDeleteStatus{
							Delete: true,
						},
					},
				},
			}, nil)

			evictionError = apierrors.NewTooManyRequestsError("error1")
			fakeEvictor.EvictOrDeletePodReturns(evictionError)

			_, err := fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: corev1.NodeSpec{
					Unschedulable: false,
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("it attempts to drain all pods", func() {
			nodeGroupDrainer := drain.NewNodeGroupDrainer(fakeClientSet, &mockNG, time.Second, time.Second, time.Second*0, false, false, 1)
			nodeGroupDrainer.SetDrainer(fakeEvictor)

			_ = nodeGroupDrainer.Drain(ctx)

			Expect(fakeEvictor.EvictOrDeletePodCallCount()).To(BeNumerically(">=", 2))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(0)).To(Equal(pods[0]))
			Expect(fakeEvictor.EvictOrDeletePodArgsForCall(1)).To(Equal(pods[1]))
		})
	})
})
