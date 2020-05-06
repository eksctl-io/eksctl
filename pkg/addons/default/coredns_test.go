package defaultaddons_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/addons"
	. "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("default addons - coredns", func() {
	var (
		rawClient *testutils.FakeRawClient
		ct        *testutils.CollectionTracker
	)

	loadSample := func(kubernetesVersion string, expectCount int) {
		samplePath := "testdata/sample-" + kubernetesVersion + ".json"
		It("can load "+samplePath+" and create objects that don't exist", func() {
			rawClient = testutils.NewFakeRawClient()

			sampleAddons := testutils.LoadSamples(samplePath)

			rawClient.UseUnionTracker = true

			for _, item := range sampleAddons {
				rc, err := rawClient.NewRawResource(item)
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
				checkCoreDNSImage(rawClient, "eu-west-1", "v"+coreDNSVersion, true)

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

	loadSampleAndCheck("1.13", "1.2.6")

	Context("[1.13 –> 1.14] can update coredns", func() {

		loadSample("1.13", 10)

		It("can load 1.13 sample", func() {
			checkCoreDNSImage(rawClient, "eu-west-1", "v1.2.6", true)
		})

		It("can update to correct version", func() {
			_, err := UpdateCoreDNS(rawClient, "eu-west-2", "1.14.x", false)
			Expect(err).ToNot(HaveOccurred())
			checkCoreDNSImage(rawClient, "eu-west-2", "v1.6.6", false)

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
})

func checkCoreDNSImage(rawClient *testutils.FakeRawClient, region, imageTag string, setImage bool) {
	coreDNS, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(CoreDNS, metav1.GetOptions{})

	Expect(err).ToNot(HaveOccurred())
	Expect(coreDNS).ToNot(BeNil())
	Expect(coreDNS.Spec.Template.Spec.Containers).To(HaveLen(1))

	if setImage {
		err := addons.UseRegionalImage(&coreDNS.Spec.Template, region)
		Expect(err).ToNot(HaveOccurred())
	}

	Expect(coreDNS.Spec.Template.Spec.Containers[0].Image).To(
		Equal("602401143452.dkr.ecr." + region + ".amazonaws.com/eks/coredns:" + imageTag),
	)
}
