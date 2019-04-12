package defaultaddons_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("default addons - coredns", func() {
	Describe("can update from kubedns to coredns add-on", func() {
		var (
			rawClient *testutils.FakeRawClient
			ct        *testutils.CollectionTracker
		)

		It("can load sample for 1.10 and create objests that don't exist", func() {
			sampleAddons := testutils.LoadSamples("testdata/sample-1.10.json")

			rawClient = testutils.NewFakeRawClient()

			rawClient.UseUnionTracker = true

			for _, item := range sampleAddons {
				rc, err := rawClient.NewRawResource(runtime.RawExtension{Object: item})
				Expect(err).ToNot(HaveOccurred())
				_, err = rc.CreateOrReplace(false)
				Expect(err).ToNot(HaveOccurred())
			}

			ct = rawClient.Collection

			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(6))

			Expect(ct.CreatedItems()).To(HaveLen(6))
		})

		It("has newly created objects", func() {
			kubeDNS, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(kubeDNS.Spec.Template.Spec.Containers).To(HaveLen(3))
			Expect(kubeDNS.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/eks/kube-dns/kube-dns:1.14.10"),
			)
			Expect(kubeDNS.Spec.Template.Spec.Containers[1].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/eks/kube-dns/dnsmasq-nanny:1.14.10"),
			)
			Expect(kubeDNS.Spec.Template.Spec.Containers[2].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/eks/kube-dns/sidecar:1.14.10"),
			)
		})

		It("can update 1.10 sample to latest", func() {
			svcOld, err := rawClient.ClientSet().CoreV1().Services(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			// test client doesn't support watching, and we would have to create some pods, so we set nil timeout
			_, err = InstallCoreDNS(rawClient, "eu-west-1", nil, false)
			Expect(err).ToNot(HaveOccurred())

			updateReqs := []string{
				"POST [/namespaces/kube-system/deployments] (kube-dns)",
				"POST [/clusterrolebindings] (aws-node)",
				"POST [/namespaces/kube-system/serviceaccounts] (coredns)",
				"POST [/namespaces/kube-system/configmaps] (coredns)",
				"POST [/namespaces/kube-system/services] (kube-dns)",
				"POST [/namespaces/kube-system/daemonsets] (aws-node)",
				"POST [/clusterroles] (system:coredns)",
				"POST [/clusterrolebindings] (system:coredns)",
				"POST [/namespaces/kube-system/deployments] (coredns)",
				"POST [/namespaces/kube-system/daemonsets] (kube-proxy)",
				"POST [/clusterroles] (aws-node)",
			}
			Expect(rawClient.Collection.Created()).To(HaveLen(len(updateReqs)))
			for _, k := range updateReqs {
				Expect(rawClient.Collection.Created()).To(HaveKey(k))
			}

			Expect(rawClient.Collection.UpdatedItems()).To(HaveLen(1))
			Expect(rawClient.Collection.Updated()).To(HaveKey(
				"PUT [/namespaces/kube-system/services/kube-dns] (kube-dns)",
			))

			svcNew, err := rawClient.ClientSet().CoreV1().Services(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			Expect(svcNew.Spec.ClusterIP).To(Equal(svcOld.Spec.ClusterIP))

			coreDNS, err := rawClient.ClientSet().ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Get(CoreDNS, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(coreDNS.Spec.Replicas).ToNot(BeNil())
			Expect(*coreDNS.Spec.Replicas == 2).To(BeTrue())
			Expect(coreDNS.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(coreDNS.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/eks/coredns:v1.1.3"),
			)
		})

	})

	Describe("can update coredns", func() {
		var (
			clientSet *fake.Clientset
		)

		check := func(imageTag string) {
			coreDNS, err := clientSet.AppsV1().Deployments(metav1.NamespaceSystem).Get(CoreDNS, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(coreDNS).ToNot(BeNil())
			Expect(coreDNS.Spec.Template.Spec.Containers).To(HaveLen(1))

			Expect(coreDNS.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/eks/coredns:" + imageTag),
			)
		}

		BeforeEach(func() {
			clientSet, _ = testutils.NewFakeClientSetWithSamples("testdata/sample-1.11.json")
		})

		It("can load 1.11 sample", func() {
			check("v1.1.3")
		})

		It("can update based to latest version", func() {
			_, err := UpdateCoreDNSImageTag(clientSet, false)
			Expect(err).ToNot(HaveOccurred())
			check("v1.2.2")
		})

		It("can dry-run update to latest version", func() {
			_, err := UpdateCoreDNSImageTag(clientSet, true)
			Expect(err).ToNot(HaveOccurred())
			check("v1.1.3")
		})
	})
})
