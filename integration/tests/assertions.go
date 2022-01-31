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

func AssertNodeTaints(nodeList *corev1.NodeList, expectedTaints []corev1.Taint) {
	//unset the time so the structs can be compared
	for _, node := range nodeList.Items {
		for _, t := range node.Spec.Taints {
			t.TimeAdded = nil
		}
	}

	for _, node := range nodeList.Items {
		for _, taint := range expectedTaints {
			Expect(node.Spec.Taints).To(ContainElement(taint))
		}
	}
}

func ListNodes(clientset kubernetes.Interface, nodeGroupName string) *corev1.NodeList {
	nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", api.NodeGroupNameLabel, nodeGroupName),
	})
	Expect(err).NotTo(HaveOccurred())
	return nodeList
}
