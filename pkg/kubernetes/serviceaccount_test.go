package kubernetes_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	. "github.com/weaveworks/eksctl/pkg/kubernetes"
)

var _ = Describe("Kubernetes serviceaccount object helpers", func() {
	var (
		clientSet *fake.Clientset
		err       error
	)

	BeforeEach(func() {
		clientSet = fake.NewSimpleClientset()
	})

	It("can create a serviceaccount object", func() {
		sa := NewServiceAccount(metav1.ObjectMeta{Name: "sa-1", Namespace: "ns-1"})

		Expect(sa.APIVersion).To(Equal("v1"))
		Expect(sa.Kind).To(Equal("ServiceAccount"))
		Expect(sa.Name).To(Equal("sa-1"))
		Expect(sa.Namespace).To(Equal("ns-1"))

		Expect(sa.Labels).To(BeEmpty())

		js, err := json.Marshal(sa)
		Expect(err).ToNot(HaveOccurred())

		expected := `{
				"apiVersion": "v1",
				"kind": "ServiceAccount",
				"metadata": {
		  			"creationTimestamp": null,
					"name": "sa-1",
					"namespace": "ns-1"
				}
			}`
		Expect(js).To(MatchJSON([]byte(expected)))
	})

	It("can create serviceaccount using fake client, and update it in incremental manner with overrides", func() {
		sa := metav1.ObjectMeta{Name: "sa-1", Namespace: "ns-1"}

		err = MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())

		ok, err := CheckNamespaceExists(clientSet, sa.Namespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(ok).To(BeTrue())

		ok, err = CheckServiceAccountExists(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())
		Expect(ok).To(BeTrue())

		{
			resp, err := clientSet.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.Name, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.Labels).To(BeEmpty())
			Expect(resp.Annotations).To(BeEmpty())
		}

		sa.Labels = map[string]string{
			"foo": "bar",
		}
		sa.Annotations = map[string]string{
			"test": "1",
		}

		err = MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())

		{
			resp, err := clientSet.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.Name, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.Labels).To(HaveKey("foo"))
			Expect(resp.Annotations).To(HaveKeyWithValue("test", "1"))
		}

		delete(sa.Labels, "foo")
		sa.Annotations["test"] = "2"

		err = MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())

		{
			resp, err := clientSet.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.Name, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.Labels).To(HaveKey("foo"))
			Expect(resp.Annotations).To(HaveKeyWithValue("test", "2"))
		}
	})

	It("can delete existsing service account, and doesn't fail if it doesn't exist", func() {
		sa := metav1.ObjectMeta{Name: "sa-2", Namespace: "ns-2"}

		err = MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())

		ok, err := CheckServiceAccountExists(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())
		Expect(ok).To(BeTrue())

		// should delete it
		err = MaybeDeleteServiceAccount(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())

		ok, err = CheckServiceAccountExists(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())
		Expect(ok).To(BeFalse())

		// shouldn't fail if it doesn't exist
		err = MaybeDeleteServiceAccount(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())

		ok, err = CheckServiceAccountExists(clientSet, sa)
		Expect(err).ToNot(HaveOccurred())
		Expect(ok).To(BeFalse())
	})
})
