package defaultaddons_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("default addons - coredns", func() {
	var (
		rawClient *testutils.FakeRawClient
		ct        *testutils.CollectionTracker
	)

	loadSample := func(kubernetesVersion string, expectCount int) {
		samplePath := "testdata/sample-" + kubernetesVersion + ".json"
		It("can load "+samplePath+" and create objests that don't exist", func() {
			rawClient = testutils.NewFakeRawClient()

			sampleAddons := testutils.LoadSamples(samplePath)

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
			Expect(ct.CreatedItems()).To(HaveLen(expectCount))
		})
	}

	loadSampleAndCheck := func(kubernetesVersion, coreDNSVersion string) {
		Context("load "+kubernetesVersion, func() {

			loadSample(kubernetesVersion, 10)

			It("can load "+kubernetesVersion+" sample", func() {
				checkCoreDNSImage(rawClient, "eu-west-1", "v"+coreDNSVersion)

				createReqs := []string{
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

				Expect(rawClient.Collection.Created()).To(HaveLen(len(createReqs)))
				for _, k := range createReqs {
					Expect(rawClient.Collection.Created()).To(HaveKey(k))
				}

				Expect(rawClient.Collection.Updated()).To(HaveLen(0))
			})
		})
	}

	Describe("[1.10 –> 1.11] can update from kube-dns to coredns", func() {

		loadSample("1.10", 6)

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
			_, err = InstallCoreDNS(rawClient, "eu-west-1", "1.11.x", nil, false)
			Expect(err).ToNot(HaveOccurred())

			createReqs := []string{
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
			Expect(rawClient.Collection.Created()).To(HaveLen(len(createReqs)))
			for _, k := range createReqs {
				Expect(rawClient.Collection.Created()).To(HaveKey(k))
			}

			Expect(rawClient.Collection.UpdatedItems()).To(HaveLen(1))
			Expect(rawClient.Collection.Updated()).To(HaveKey(
				"PUT [/namespaces/kube-system/services/kube-dns] (kube-dns)",
			))

			svcNew, err := rawClient.ClientSet().CoreV1().Services(metav1.NamespaceSystem).Get(KubeDNS, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			Expect(svcNew.Spec.ClusterIP).To(Equal(svcOld.Spec.ClusterIP))

			coreDNS, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(CoreDNS, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(coreDNS.Spec.Replicas).ToNot(BeNil())
			Expect(*coreDNS.Spec.Replicas == 2).To(BeTrue())
			Expect(coreDNS.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(coreDNS.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/eks/coredns:v1.1.3"),
			)
		})

	})

	loadSampleAndCheck("1.11", "1.1.3")

	Context("[1.11 –> 1.12] can update coredns", func() {

		loadSample("1.11", 10)

		It("can load 1.11 sample", func() {
			checkCoreDNSImage(rawClient, "eu-west-1", "v1.1.3")
		})

		It("can update to correct version", func() {
			_, err := UpdateCoreDNS(rawClient, "eu-west-2", "1.12.x", false)
			Expect(err).ToNot(HaveOccurred())
			checkCoreDNSImage(rawClient, "eu-west-2", "v1.2.2")

			createReqs := []string{
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

			Expect(rawClient.Collection.Created()).To(HaveLen(len(createReqs)))
			for _, k := range createReqs {
				Expect(rawClient.Collection.Created()).To(HaveKey(k))
			}

			updateReqs := []string{
				"PUT [/namespaces/kube-system/serviceaccounts/coredns] (coredns)",
				"PUT [/namespaces/kube-system/configmaps/coredns] (coredns)",
				"PUT [/namespaces/kube-system/services/kube-dns] (kube-dns)",
				"PUT [/clusterroles/system:coredns] (system:coredns)",
				"PUT [/clusterrolebindings/system:coredns] (system:coredns)",
				"PUT [/namespaces/kube-system/deployments/coredns] (coredns)",
			}

			Expect(rawClient.Collection.Updated()).To(HaveLen(len(updateReqs)))
			for _, k := range updateReqs {
				Expect(rawClient.Collection.Updated()).To(HaveKey(k))
			}
		})
	})

	loadSampleAndCheck("1.12", "1.2.2")

	Context("[1.12 –> 1.13] can update coredns", func() {

		loadSample("1.12", 10)

		It("can load 1.11 sample", func() {
			checkCoreDNSImage(rawClient, "eu-west-1", "v1.2.2")
		})

		It("can update to correct version", func() {
			_, err := UpdateCoreDNS(rawClient, "eu-west-2", "1.13.x", false)
			Expect(err).ToNot(HaveOccurred())
			checkCoreDNSImage(rawClient, "eu-west-2", "v1.2.6")

			createReqs := []string{
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

			Expect(rawClient.Collection.Created()).To(HaveLen(len(createReqs)))
			for _, k := range createReqs {
				Expect(rawClient.Collection.Created()).To(HaveKey(k))
			}

			updateReqs := []string{
				"PUT [/namespaces/kube-system/serviceaccounts/coredns] (coredns)",
				"PUT [/namespaces/kube-system/configmaps/coredns] (coredns)",
				"PUT [/namespaces/kube-system/services/kube-dns] (kube-dns)",
				"PUT [/clusterroles/system:coredns] (system:coredns)",
				"PUT [/clusterrolebindings/system:coredns] (system:coredns)",
				"PUT [/namespaces/kube-system/deployments/coredns] (coredns)",
			}

			Expect(rawClient.Collection.Updated()).To(HaveLen(len(updateReqs)))
			for _, k := range updateReqs {
				Expect(rawClient.Collection.Updated()).To(HaveKey(k))
			}
		})
	})

	loadSampleAndCheck("1.13", "1.2.6")
})

func checkCoreDNSImage(rawClient *testutils.FakeRawClient, region, imageTag string) {
	coreDNS, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(CoreDNS, metav1.GetOptions{})

	Expect(err).ToNot(HaveOccurred())
	Expect(coreDNS).ToNot(BeNil())
	Expect(coreDNS.Spec.Template.Spec.Containers).To(HaveLen(1))

	Expect(coreDNS.Spec.Template.Spec.Containers[0].Image).To(
		Equal("602401143452.dkr.ecr." + region + ".amazonaws.com/eks/coredns:" + imageTag),
	)
}
