package defaultaddons_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("default addons", func() {
	Describe("can load a set of resources and create a fake client", func() {
		It("can create the fake client and verify objects get loaded client", func() {
			clientSet, _ := testutils.NewFakeClientSetWithSamples("testdata/sample-1.12.json")

			nsl, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(nsl.Items).To(HaveLen(0))

			dl, err := clientSet.AppsV1().Deployments(metav1.NamespaceAll).List(metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(dl.Items).To(HaveLen(1))
			Expect(dl.Items[0].Spec.Template.Spec.Containers).To(HaveLen(1))

			kubeProxy, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(KubeProxy, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(kubeProxy).ToNot(BeNil())
			Expect(kubeProxy.Spec.Template.Spec.Containers).To(HaveLen(1))
		})
	})
})
