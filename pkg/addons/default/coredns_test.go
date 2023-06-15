package defaultaddons_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	da "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("default addons - coredns", func() {
	var (
		rawClient           *testutils.FakeRawClient
		input               da.AddonInput
		region              string
		controlPlaneVersion string
		kubernetesVersion   string
	)

	BeforeEach(func() {
		rawClient = testutils.NewFakeRawClient()
		rawClient.UseUnionTracker = true
		region = "eu-west-2"
		controlPlaneVersion = "1.23.x"
		kubernetesVersion = "1.22"

		input = da.AddonInput{
			RawClient:           rawClient,
			ControlPlaneVersion: controlPlaneVersion,
			Region:              region,
		}
	})

	Context("UpdateCoreDNS", func() {
		var (
			expectedImageTag string
		)

		BeforeEach(func() {
			createCoreDNSFromTestSample(rawClient, kubernetesVersion)
			expectedImageTag = "v1.8.7-eksbuild.2"
		})

		It("updates coredns to the correct version", func() {
			_, err := da.UpdateCoreDNS(context.Background(), input, false)
			Expect(err).NotTo(HaveOccurred())

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

			Expect(coreDNSImage(rawClient)).To(
				Equal("602401143452.dkr.ecr." + region + ".amazonaws.com/eks/coredns:" + expectedImageTag),
			)
		})
	})
})

func createCoreDNSFromTestSample(rawClient *testutils.FakeRawClient, kubernetesVersion string) {
	samplePath := "testdata/sample-" + kubernetesVersion + ".json"
	sampleAddons := testutils.LoadSamples(samplePath)

	for _, item := range sampleAddons {
		rc, err := rawClient.NewRawResource(item)
		Expect(err).NotTo(HaveOccurred())
		_, err = rc.CreateOrReplace(false)
		Expect(err).NotTo(HaveOccurred())
	}

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
}

func coreDNSImage(rawClient *testutils.FakeRawClient) string {
	coreDNS, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(context.Background(), da.CoreDNS, metav1.GetOptions{})

	Expect(err).NotTo(HaveOccurred())
	Expect(coreDNS).NotTo(BeNil())
	Expect(coreDNS.Spec.Template.Spec.Containers).To(HaveLen(1))

	return coreDNS.Spec.Template.Spec.Containers[0].Image
}
