package kubernetes_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("Kubernetes client wrappers", func() {
	Describe("can create or replace missing objects", func() {
		It("can update objects that already exist", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.16.json")
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, false, false, ct)

				exists, err := rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeTrue()) // The Kubernetes resource already exists.

				_, err = rc.CreateOrReplace(false)
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "GET", "GET", "PUT"}))

				exists, err = rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeTrue()) // The Kubernetes resource still exists.
			}

			Expect(ct.Updated()).ToNot(BeEmpty())
			Expect(ct.UpdatedItems()).To(HaveLen(10))
			Expect(ct.Created()).To(BeEmpty())
			Expect(ct.CreatedItems()).To(BeEmpty())
		})

		It("can create objects that don't exist yet", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.16.json")
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, true, false, ct)

				exists, err := rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeFalse()) // The Kubernetes resource has not been created yet.

				_, err = rc.CreateOrReplace(false)
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "GET", "POST"}))

				exists, err = rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeTrue()) // The Kubernetes resource has not been created yet.
			}

			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(10))
			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.UpdatedItems()).To(BeEmpty())
		})

		It("can update objects that already exist", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.16.json")
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, false, false, ct)

				exists, err := rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeTrue()) // The Kubernetes resource already exists.

				_, err = rc.CreateOrReplace(false)
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "GET", "GET", "PUT"}))

				exists, err = rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeTrue()) // The Kubernetes resource still exists.
			}

			Expect(ct.Updated()).ToNot(BeEmpty())
			Expect(ct.UpdatedItems()).To(HaveLen(10))
			Expect(ct.Created()).To(BeEmpty())
			Expect(ct.CreatedItems()).To(BeEmpty())
		})

		It("can create objects and update objects in a union", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.16.json")

			rawClient := testutils.NewFakeRawClient()

			rawClient.UseUnionTracker = true

			for _, item := range sampleAddons {
				rc, err := rawClient.NewRawResource(item)
				Expect(err).ToNot(HaveOccurred())
				_, err = rc.CreateOrReplace(false)
				Expect(err).ToNot(HaveOccurred())
			}

			ct := rawClient.Collection

			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(10))
			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.UpdatedItems()).To(BeEmpty())

			dsl, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).List(metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(dsl.Items).To(HaveLen(2))

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get("aws-node", metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni:v1.6.1"),
			)

			saTest1 := &corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ServiceAccount",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test1",
					Namespace: metav1.NamespaceDefault,
				},
			}

			saTest2a := &corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ServiceAccount",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test2",
					Namespace: metav1.NamespaceDefault,
					Labels:    map[string]string{"test": "2a"},
				},
			}
			saTest2b := &corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ServiceAccount",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test2",
					Namespace: metav1.NamespaceDefault,
					Labels:    map[string]string{"test": "2b"},
				},
			}

			for _, item := range []runtime.Object{saTest1, saTest2a, saTest2b} {
				rc, err := rawClient.NewRawResource(item)
				Expect(err).ToNot(HaveOccurred())
				_, err = rc.CreateOrReplace(false)
				Expect(err).ToNot(HaveOccurred())
			}

			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(10 + 2))
			Expect(ct.UpdatedItems()).ToNot(BeEmpty())
			Expect(ct.UpdatedItems()).To(HaveLen(1))

			_, err = rawClient.ClientSet().CoreV1().ServiceAccounts(metav1.NamespaceDefault).Get("test1", metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			_, err = rawClient.ClientSet().CoreV1().ServiceAccounts(metav1.NamespaceDefault).Get("test2", metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			_, err = rawClient.ClientSet().CoreV1().ServiceAccounts(metav1.NamespaceDefault).Create(saTest1)
			Expect(err).To(HaveOccurred())

			err = rawClient.ClientSet().CoreV1().ServiceAccounts(metav1.NamespaceDefault).Delete("test1", &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())

			// saving a clientset instance results in objects being trackable,
			// but only as far as the clientset instance is concerned
			// we need to find a way to fix this, the test is only to document
			// this limitation
			c := rawClient.ClientSet().CoreV1().ServiceAccounts(metav1.NamespaceDefault)
			err = c.Delete("test1", &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())
			err = c.Delete("test1", &metav1.DeleteOptions{})
			Expect(err).To(HaveOccurred())

			// however deletions of raw resources are trackable
			rc, err := rawClient.NewRawResource(saTest1)
			Expect(err).ToNot(HaveOccurred())
			_, err = rc.Helper.Delete(rc.Info.Namespace, rc.Info.Name)
			Expect(err).ToNot(HaveOccurred())
			_, err = rawClient.ClientSet().CoreV1().ServiceAccounts(metav1.NamespaceDefault).Get("test1", metav1.GetOptions{})
			Expect(err).To(HaveOccurred())
		})

		It("can delete existing objects", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.16.json")
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, false, false, ct)

				exists, err := rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeTrue()) // The Kubernetes resource already exists.

				status, err := rc.DeleteSync()
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(fmt.Sprintf("deleted %q", rc)))
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "GET", "DELETE", "GET"}))

				exists, err = rc.Exists()
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeFalse()) // The Kubernetes resource no longer exists.
			}

			Expect(ct.Created()).To(BeEmpty())
			Expect(ct.CreatedItems()).To(BeEmpty())
			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.UpdatedItems()).To(BeEmpty())
			Expect(ct.Deleted()).ToNot(BeEmpty())
			Expect(ct.DeletedItems()).To(HaveLen(10))
		})
	})
})
