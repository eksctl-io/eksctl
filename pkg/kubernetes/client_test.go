package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("default addons", func() {
	Describe("can create or replace missing objects", func() {
		It("can update objects that already exist", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.10.json")
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, false, ct)
				_, err := rc.CreateOrReplace()
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "GET", "PUT"}))
			}

			Expect(ct.Updated()).ToNot(BeEmpty())
			Expect(ct.UpdatedItems()).To(HaveLen(6))
			Expect(ct.Created()).To(BeEmpty())
			Expect(ct.CreatedItems()).To(BeEmpty())
		})

		It("can create objects that don't exist yet", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.10.json")
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, true, ct)
				_, err := rc.CreateOrReplace()
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "POST"}))
			}

			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(6))
			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.UpdatedItems()).To(BeEmpty())
		})

		It("can create objests that don't exist, and convert into a clientset", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.10.json")

			rawClient := testutils.NewFakeRawClient()

			rawClient.AssumeObjectsMissing = true

			for _, item := range sampleAddons {
				rc, err := rawClient.NewRawResource(runtime.RawExtension{Object: item})
				Expect(err).ToNot(HaveOccurred())
				_, err = rc.CreateOrReplace()
				Expect(err).ToNot(HaveOccurred())
			}

			ct := rawClient.Collection

			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(6))

			dsl, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).List(metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(dsl.Items).To(HaveLen(2))

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get("aws-node", metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/amazon-k8s-cni:v1.3.2"),
			)
		})
	})
})
