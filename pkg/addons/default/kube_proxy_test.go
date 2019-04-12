package defaultaddons_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("default addons - kube-proxy", func() {
	Describe("can update kube-proxy", func() {
		var (
			clientSet *fake.Clientset
		)

		check := func(imageTag string) {
			kubeProxy, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(KubeProxy, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(kubeProxy).ToNot(BeNil())
			Expect(kubeProxy.Spec.Template.Spec.Containers).To(HaveLen(1))

			Expect(kubeProxy.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/eks/kube-proxy:" + imageTag),
			)
		}

		BeforeEach(func() {
			clientSet, _ = testutils.NewFakeClientSetWithSamples("testdata/sample-1.10.json")
		})

		It("can load 1.10 sample", func() {
			check("v1.10.3")
		})

		It("can update based on control plane version", func() {
			_, err := UpdateKubeProxyImageTag(clientSet, "1.11.0", false)
			Expect(err).ToNot(HaveOccurred())
			check("v1.11.0")
		})

		It("can dry-run update based on control plane version", func() {
			_, err := UpdateKubeProxyImageTag(clientSet, "1.12.1", true)
			Expect(err).ToNot(HaveOccurred())
			check("v1.10.3")
		})
	})
})
