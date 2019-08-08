package kubernetes_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/yaml"

	. "github.com/weaveworks/eksctl/pkg/kubernetes"
)

var _ = Describe("Kubernetes namespace object helpers", func() {
	var (
		clientSet *fake.Clientset
		err       error
	)

	BeforeEach(func() {
		clientSet = fake.NewSimpleClientset()
	})

	It("can create a namespace object", func() {
		ns := NewNamespace("ns123")

		Expect(ns.APIVersion).To(Equal("v1"))
		Expect(ns.Kind).To(Equal("Namespace"))
		Expect(ns.Name).To(Equal("ns123"))

		Expect(ns.Labels).To(BeEmpty())

		js, err := json.Marshal(ns)
		Expect(err).ToNot(HaveOccurred())

		expected := `{
				"apiVersion": "v1",
				"kind": "Namespace",
				"metadata": {
					"creationTimestamp": null,
					"name": "ns123"
				},
				"spec": {},
				"status": {}
			}`
		Expect(js).To(MatchJSON([]byte(expected)))
	})

	It("can create a clean serialised namespace object", func() {
		ys := NewNamespaceYAML("ns123")
		ns := &corev1.Namespace{}

		err = yaml.Unmarshal(ys, ns)

		Expect(err).ToNot(HaveOccurred())

		Expect(ns.APIVersion).To(Equal("v1"))
		Expect(ns.Kind).To(Equal("Namespace"))
		Expect(ns.Name).To(Equal("ns123"))

		Expect(*ns).To(Equal(*NewNamespace("ns123")))
	})

	It("can create namespace using fake client and check confirm that it exists", func() {
		err = MaybeCreateNamespace(clientSet, "ns-1")
		Expect(err).ToNot(HaveOccurred())

		ok, err := CheckNamespaceExists(clientSet, "ns-1")
		Expect(err).ToNot(HaveOccurred())
		Expect(ok).To(BeTrue())
	})
})
