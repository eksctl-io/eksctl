package defaultaddons_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("default addons - aws-node", func() {
	Describe("properly checks for multi-architecture support", func() {
		var (
			rawClient *testutils.FakeRawClient
			ct        *testutils.CollectionTracker
		)
		loadSample := func(f string) {
			sampleAddons := testutils.LoadSamples(f)

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
		}
		It("reports that 1.15 sample needs an update", func() {
			loadSample("testdata/sample-1.15.json")
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := DoesAWSNodeSupportMultiArch(rawClient, "eu-west-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(needsUpdate).To(BeFalse())
		})
		It("reports that sample with 1.6.3-eksbuild.1 doesn't need an update", func() {
			loadSample("testdata/sample-1.16-eksbuild.1.json")
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := DoesAWSNodeSupportMultiArch(rawClient, "eu-west-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(needsUpdate).To(BeTrue())
		})
		It("reports that sample with 1.7.6 doesn't need an update", func() {
			loadSample("testdata/sample-1.16-v1.7.json")
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := DoesAWSNodeSupportMultiArch(rawClient, "eu-west-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(needsUpdate).To(BeTrue())
		})
	})

	Describe("can update aws-node add-on to multi-architecture images", func() {
		var (
			rawClient *testutils.FakeRawClient
			ct        *testutils.CollectionTracker
		)

		It("can load sample for 1.15 and create objects that don't exist", func() {
			sampleAddons := testutils.LoadSamples("testdata/sample-1.15.json")

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

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni:v1.5.7"),
			)

		})

		It("can update 1.15 sample to latest multi-architecture image", func() {
			rawClient.AssumeObjectsMissing = false

			preUpdateAwsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			_, err = UpdateAWSNode(rawClient, "eu-west-1", false)
			Expect(err).ToNot(HaveOccurred())
			Expect(rawClient.Collection.UpdatedItems()).To(HaveLen(3))
			Expect(rawClient.Collection.UpdatedItems()).ToNot(ContainElement(PointTo(MatchFields(IgnoreMissing|IgnoreExtras, Fields{
				"TypeMeta": MatchFields(IgnoreMissing|IgnoreExtras, Fields{"Kind": Equal("ServiceAccount")}),
			}))))
			Expect(rawClient.Collection.CreatedItems()).To(HaveLen(10))

			rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).ToNot(
				Equal(preUpdateAwsNode.Spec.Template.Spec.Containers[0].Image),
			)
			Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
				HavePrefix("602401143452.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni-init"),
			)
			rawClient.ClearUpdated()
		})

		It("can update 1.15 sample for different region to multi-architecture image", func() {
			rawClient.ClientSetUseUpdatedObjects = false // must be set for subsequent UpdateAWSNode

			preUpdateAwsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			_, err = UpdateAWSNode(rawClient, "us-east-1", false)
			Expect(err).ToNot(HaveOccurred())

			rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).ToNot(
				Equal(preUpdateAwsNode.Spec.Template.Spec.Containers[0].Image),
			)
			Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
				HavePrefix("602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon-k8s-cni-init"),
			)
		})

		It("can update 1.15 sample for china region to multi-architecture image", func() {
			rawClient.ClientSetUseUpdatedObjects = false // must be set for subsequent UpdateAWSNode

			preUpdateAwsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			_, err = UpdateAWSNode(rawClient, "cn-northwest-1", false)
			Expect(err).ToNot(HaveOccurred())

			rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

			awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), AWSNode, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.Containers[0].Image).ToNot(
				Equal(preUpdateAwsNode.Spec.Template.Spec.Containers[0].Image),
			)
			Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
			Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
				HavePrefix("961992271922.dkr.ecr.cn-northwest-1.amazonaws.com.cn/amazon-k8s-cni-init"),
			)
		})

		It("detects matching image version when determining plan", func() {
			// updating from latest to latest needs no updating
			needsUpdate, err := UpdateAWSNode(rawClient, "eu-west-2", true)
			Expect(err).ToNot(HaveOccurred())
			Expect(needsUpdate).To(BeFalse())
		})
	})
})
