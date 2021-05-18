package tests

import (
	"context"
	"fmt"

	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func AssertNodeTaints(clientset kubernetes.Interface, nodeGroupName string, expectedTaints []corev1.Taint) {
	nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", api.NodeGroupNameLabel, nodeGroupName),
	})
	Expect(err).ToNot(HaveOccurred())

	for _, node := range nodeList.Items {
		for i, t := range node.Spec.Taints {
			expected := expectedTaints[i]
			Expect(t.Key).To(Equal(expected.Key))
			Expect(t.Value).To(Equal(expected.Value))
			Expect(t.Effect).To(Equal(expected.Effect))
		}
	}
}
