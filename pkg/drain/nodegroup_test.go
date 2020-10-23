package drain_test

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Drain", func() {
	When("all nodes drain successfully", func() {
		It("does not error", func() {
			mockNG := mocks.KubeNodeGroup{}
			mockNG.Mock.On("NameString").Return("node-1")
			mockNG.Mock.On("ListOptions").Return(metav1.ListOptions{})

			fakeClientSet := fake.NewSimpleClientset()

			err := drain.NodeGroup(fakeClientSet, &mockNG, time.Second, time.Second, false)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("one of the nodes cannot be evicted", func() {
		It("returns an error", func() {
			mockNG := mocks.KubeNodeGroup{}
			mockNG.Mock.On("NameString").Return("node-1")
			mockNG.Mock.On("ListOptions").Return(metav1.ListOptions{})

			fakeClientSet := fake.NewSimpleClientset()
			//fakeClientSet.CoreV1().Pods(metav1.NamespaceAll).Create(v1.Pod{})

			err := drain.NodeGroup(fakeClientSet, &mockNG, time.Second, time.Second, false)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
