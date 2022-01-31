package kubernetes

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("NodeGroup", func() {
	Describe("GetNodegroupKubernetesVersion", func() {
		const (
			ngName = "test-ng"
		)

		type ngEntry struct {
			mockCalls func(*fake.Clientset)

			expectedVersion string
			expError        error
		}

		DescribeTable("with node returned by kube API", func(t ngEntry) {
			fakeClientSet := fake.NewSimpleClientset()

			if t.mockCalls != nil {
				t.mockCalls(fakeClientSet)
			}

			version, err := GetNodegroupKubernetesVersion(fakeClientSet.CoreV1().Nodes(), ngName)

			if t.expectedVersion == "" {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(t.expError.Error())))
				return
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(Equal(t.expectedVersion))
		},

			Entry("[happy path] returns correct version", ngEntry{
				mockCalls: func(c *fake.Clientset) {
					_, err := c.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-node",
							Labels: map[string]string{
								api.NodeGroupNameLabel: ngName,
							},
						},
						Status: v1.NodeStatus{
							NodeInfo: v1.NodeSystemInfo{
								KubeletVersion: "1.19.6",
							},
						},
					}, metav1.CreateOptions{})
					Expect(err).NotTo(HaveOccurred())
				},
				expectedVersion: "1.19.6",
			}),

			Entry("[happy path] returns correct version with version trimming", ngEntry{
				mockCalls: func(c *fake.Clientset) {
					_, err := c.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-node",
							Labels: map[string]string{
								api.NodeGroupNameLabel: ngName,
							},
						},
						Status: v1.NodeStatus{
							NodeInfo: v1.NodeSystemInfo{
								KubeletVersion: "v1.19.6-eks-49a6c0",
							},
						},
					}, metav1.CreateOptions{})
					Expect(err).NotTo(HaveOccurred())
				},
				expectedVersion: "1.19.6",
			}),

			Entry("fails to list nodes", ngEntry{
				mockCalls: func(c *fake.Clientset) {
					_, err := c.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-node",
						},
					}, metav1.CreateOptions{})
					Expect(err).NotTo(HaveOccurred())
				},
				expError: errors.New("no nodes were found"),
			}),
		)
	})
})
