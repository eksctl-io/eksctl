package defaultaddons_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	da "github.com/weaveworks/eksctl/pkg/addons/default"

	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("default addons - aws-node", func() {
	var (
		input     da.AddonInput
		rawClient *testutils.FakeRawClient
		ct        *testutils.CollectionTracker
	)

	BeforeEach(func() {
		rawClient = testutils.NewFakeRawClient()
		input = da.AddonInput{
			RawClient:           rawClient,
			ControlPlaneVersion: "1.16.0",
			Region:              "eu-west-1",
		}
	})
	Describe("properly checks for multi-architecture support", func() {
		loadSample := func(f string) {
			sampleAddons := testutils.LoadSamples(f)

			rawClient.AssumeObjectsMissing = true

			for _, item := range sampleAddons {
				rc, err := rawClient.NewRawResource(item)
				Expect(err).NotTo(HaveOccurred())
				_, err = rc.CreateOrReplace(false)
				Expect(err).NotTo(HaveOccurred())
			}

			ct = rawClient.Collection

			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.Created()).NotTo(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(10))
		}
		It("reports that 1.15 sample needs an update", func() {
			loadSample("testdata/sample-1.15.json")
			input.ControlPlaneVersion = "1.15.0"
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := da.DoesAWSNodeSupportMultiArch(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(needsUpdate).To(BeFalse())
		})
		It("reports that sample with 1.6.3-eksbuild.1 doesn't need an update", func() {
			loadSample("testdata/sample-1.16-eksbuild.1.json")
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := da.DoesAWSNodeSupportMultiArch(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(needsUpdate).To(BeTrue())
		})
		It("reports that sample with 1.7.6 doesn't need an update", func() {
			loadSample("testdata/sample-1.16-v1.7.json")
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := da.DoesAWSNodeSupportMultiArch(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(needsUpdate).To(BeTrue())
		})
	})

	Describe("can update aws-node add-on to multi-architecture images", func() {
		BeforeEach(func() {
			sampleAddons := testutils.LoadSamples("testdata/sample-1.15.json")

			rawClient.AssumeObjectsMissing = true

			for _, item := range sampleAddons {
				rc, err := rawClient.NewRawResource(item)
				Expect(err).NotTo(HaveOccurred())
				_, err = rc.CreateOrReplace(false)
				Expect(err).NotTo(HaveOccurred())
			}

			ct = rawClient.Collection

			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.Created()).NotTo(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(10))
		})

		It("can update the aws-node successfully", func() {
			By("updating the 1.15 sample to latest multi-architecture image", func() {
				rawClient.AssumeObjectsMissing = false

				preUpdateAwsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				_, err = da.UpdateAWSNode(input, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(rawClient.Collection.UpdatedItems()).To(HaveLen(3))
				Expect(rawClient.Collection.UpdatedItems()).NotTo(ContainElement(PointTo(MatchFields(IgnoreMissing|IgnoreExtras, Fields{
					"TypeMeta": MatchFields(IgnoreMissing|IgnoreExtras, Fields{"Kind": Equal("ServiceAccount")}),
				}))))
				Expect(rawClient.Collection.CreatedItems()).To(HaveLen(10))

				rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.Containers[0].Image).NotTo(
					Equal(preUpdateAwsNode.Spec.Template.Spec.Containers[0].Image),
				)
				Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
					HavePrefix("602401143452.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni-init"),
				)
				rawClient.ClearUpdated()
			})

			By("updating the 1.15 sample for different region to multi-architecture image", func() {
				input.Region = "us-east-1"
				rawClient.ClientSetUseUpdatedObjects = false // must be set for subsequent UpdateAWSNode

				preUpdateAwsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				_, err = da.UpdateAWSNode(input, false)
				Expect(err).NotTo(HaveOccurred())

				rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.Containers[0].Image).NotTo(
					Equal(preUpdateAwsNode.Spec.Template.Spec.Containers[0].Image),
				)
				Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
					HavePrefix("602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon-k8s-cni-init"),
				)
			})

			By("updating the 1.15 sample for china region to multi-architecture image", func() {
				input.Region = "cn-northwest-1"
				rawClient.ClientSetUseUpdatedObjects = false // must be set for subsequent UpdateAWSNode

				preUpdateAwsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				_, err = da.UpdateAWSNode(input, false)
				Expect(err).NotTo(HaveOccurred())

				rawClient.ClientSetUseUpdatedObjects = true // for verification of updated objects

				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.Containers[0].Image).NotTo(
					Equal(preUpdateAwsNode.Spec.Template.Spec.Containers[0].Image),
				)
				Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
					HavePrefix("961992271922.dkr.ecr.cn-northwest-1.amazonaws.com.cn/amazon-k8s-cni-init"),
				)
			})

			By("detecting matching image version when determining plan", func() {
				// updating from latest to latest needs no updating
				needsUpdate, err := da.UpdateAWSNode(input, true)
				Expect(err).NotTo(HaveOccurred())
				Expect(needsUpdate).To(BeFalse())
			})
		})
	})
})
