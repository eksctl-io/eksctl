package defaultaddons_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("default addons - aws-node", func() {
	Describe("can update aws-node add-on", func() {
		var (
			rawClient *testutils.FakeRawClient
			ct        *testutils.CollectionTracker
		)

		It("can load sample for 1.12 and create objests that don't exist", func() {
			sampleAddons := testutils.LoadSamples("testdata/sample-1.12.json")

			rawClient = testutils.NewFakeRawClient()

			rawClient.AssumeObjectsMissing = true

			for _, item := range sampleAddons {
				rc, err := rawClient.NewRawResource(item)
				Expect(err).ToNot(HaveOccurred())
				_, err = rc.CreateOrReplace(false)
				Expect(err).ToNot(HaveOccurred())
			}

			ct = rawClient.Collection

			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(10))
		})

		It("has newly created objects", func() {
			rawClient.ClientSetUseUpdatedObjects = false // must be set for initial verification, and for subsequent UpdateAWSNode

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni:v1.4.1"),
			)

		})

		It("can update 1.12 sample to latest", func() {
			rawClient.AssumeObjectsMissing = false

			_, err := UpdateAWSNode(rawClient, "eu-west-2", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(rawClient.Collection.UpdatedItems()).To(HaveLen(4))
			Expect(rawClient.Collection.CreatedItems()).To(HaveLen(10))

			rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/amazon-k8s-cni:v1.5.0"),
			)

			rawClient.ClearUpdated()
		})

		It("can update 1.12 sample for different region", func() {
			rawClient.ClientSetUseUpdatedObjects = false // must be set for subsequent UpdateAWSNode

			_, err := UpdateAWSNode(rawClient, "us-east-1", false)
			Expect(err).ToNot(HaveOccurred())

			rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon-k8s-cni:v1.5.0"),
			)
		})
	})
})
