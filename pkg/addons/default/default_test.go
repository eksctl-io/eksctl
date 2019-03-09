package defaultaddons_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/addons/default"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("default addons", func() {
	Describe("can load a resources and create fake client", func() {

		var (
			sampleAddons []runtime.Object
			clientSet    *fake.Clientset
		)

		// TODO: make ClientResource work
		// TODO: test UpdateKubeProxy
		// TODO: test UpdateAWSNode

		It("can load sample addons", func() {
			sampleAddonsData, err := ioutil.ReadFile("testdata/sample-1.10.json")
			Expect(err).To(Not(HaveOccurred()))
			sampleAddonsList, err := NewList(sampleAddonsData)
			Expect(err).To(Not(HaveOccurred()))
			Expect(sampleAddonsList).To(Not(BeNil()))
			Expect(sampleAddonsList.Items).To(HaveLen(7))

			for _, item := range sampleAddonsList.Items {
				kind := item.Object.GetObjectKind().GroupVersionKind().Kind
				if kind == "CustomResourceDefinition" {
					continue // fake client doesn't support CRDs, save it from a panic
				}
				sampleAddons = append(sampleAddons, item.Object)
			}
		})

		It("can create the fake client and verify objects get loaded client", func() {
			Expect(sampleAddons).To(HaveLen(6))
			clientSet = fake.NewSimpleClientset(sampleAddons...)

			nsl, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).To(Not(HaveOccurred()))
			Expect(nsl.Items).To(HaveLen(0))

			dl, err := clientSet.AppsV1().Deployments(metav1.NamespaceAll).List(metav1.ListOptions{})
			Expect(err).To(Not(HaveOccurred()))
			Expect(dl.Items).To(HaveLen(1))
			Expect(dl.Items[0].Spec.Template.Spec.Containers).To(HaveLen(3))

			kubeProxy, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(KubeProxy, metav1.GetOptions{})
			Expect(err).To(Not(HaveOccurred()))
			Expect(kubeProxy).To(Not(BeNil()))
			Expect(kubeProxy.Spec.Template.Spec.Containers).To(HaveLen(1))
		})
	})

	Describe("can load object", func() {

		Context("can load and flatten deeply nested lists", func() {
			It("loads all items into flattened list without errors", func() {
				jb, err := ioutil.ReadFile("testdata/misc-sample-nested-list-1.json")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(jb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(6))
			})
		})

		Context("can load and flatten deeply nested lists", func() {
			It("flatten all items into an empty list without errors", func() {
				jb, err := ioutil.ReadFile("testdata/misc-sample-empty-list-1.json")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(jb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(0))
			})
		})

		Context("can combine empty nested lists from a multidoc", func() {
			It("can load without errors", func() {
				yb, err := ioutil.ReadFile("testdata/misc-sample-multidoc-empty-lists-1.yaml")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(yb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(0))
			})
		})

		Context("can combine two empty lists from a multidoc", func() {

			It("can load without errors", func() {
				yb, err := ioutil.ReadFile("testdata/misc-sample-multidoc-empty-lists-2.yaml")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(yb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(0))
			})
		})

		Context("can combine empty and non-empty lists from a multidoc", func() {
			It("can load without errors", func() {
				yb, err := ioutil.ReadFile("testdata/misc-sample-multidoc-nested-lists-1.yaml")
				Expect(err).To(Not(HaveOccurred()))

				list, err := NewList(yb)
				Expect(err).ToNot(HaveOccurred())
				Expect(list).ToNot(BeNil())
				Expect(list.Items).To(HaveLen(4))
			})
		})
	})
})
