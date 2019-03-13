package defaultaddons_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("default addons - 1.10/aws-node", func() {
	Describe("can load a resources and create fake client", func() {

		It("can create objests that don't exist, and convert into a clientset", func() {
			sampleAddons := testutils.LoadSamples("testdata/sample-1.10.json")

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

			Expect(rawClient.Collection.CreatedItems()).To(HaveLen(6))

			rawClient.ClientSetUseUpdatedObjects = false // must be set for initial verification, and for subsequent UpdateAWSNode

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/amazon-k8s-cni:v1.3.2"),
			)

			rawClient.AssumeObjectsMissing = false

			err = UpdateAWSNode(rawClient, "eu-west-2")
			Expect(err).ToNot(HaveOccurred())
			Expect(rawClient.Collection.UpdatedItems()).To(HaveLen(4))
			Expect(rawClient.Collection.CreatedItems()).To(HaveLen(6))

			rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

			awsNode, err = rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-2.amazonaws.com/amazon-k8s-cni:v1.3.2"),
			)

			rawClient.ClearUpdated()
			rawClient.ClientSetUseUpdatedObjects = false // must be set for subsequent UpdateAWSNode

			err = UpdateAWSNode(rawClient, "us-east-1")
			Expect(err).ToNot(HaveOccurred())

			rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

			awsNode, err = rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon-k8s-cni:v1.3.2"),
			)
		})

	})
})
