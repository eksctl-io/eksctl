package defaultaddons_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	da "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("AWS Node", func() {
	var (
		input     da.AddonInput
		rawClient *testutils.FakeRawClient
	)

	BeforeEach(func() {
		rawClient = testutils.NewFakeRawClient()
		input = da.AddonInput{
			RawClient:           rawClient,
			ControlPlaneVersion: "1.16.0",
			Region:              "eu-west-1",
		}
	})

	Describe("DoesAWSNodeSupportMultiArch", func() {
		It("reports that 1.15 sample needs an update", func() {
			loadSamples(rawClient, "testdata/sample-1.15.json")
			input.ControlPlaneVersion = "1.15.0"
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := da.DoesAWSNodeSupportMultiArch(context.Background(), rawClient.ClientSet())
			Expect(err).NotTo(HaveOccurred())
			Expect(needsUpdate).To(BeFalse())
		})

		It("reports that sample with 1.6.3-eksbuild.1 doesn't need an update", func() {
			loadSamples(rawClient, "testdata/sample-1.16-eksbuild.1.json")
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := da.DoesAWSNodeSupportMultiArch(context.Background(), rawClient.ClientSet())
			Expect(err).NotTo(HaveOccurred())
			Expect(needsUpdate).To(BeTrue())
		})

		It("reports that sample with 1.7.6 doesn't need an update", func() {
			loadSamples(rawClient, "testdata/sample-1.16-v1.7.json")
			rawClient.AssumeObjectsMissing = false

			needsUpdate, err := da.DoesAWSNodeSupportMultiArch(context.Background(), rawClient.ClientSet())
			Expect(err).NotTo(HaveOccurred())
			Expect(needsUpdate).To(BeTrue())
		})
	})

	Describe("UpdateAWSNode", func() {
		var preUpdateAwsNode *v1.DaemonSet
		const expectedVersion = "v1.18.2"
		BeforeEach(func() {
			loadSamples(rawClient, "testdata/sample-1.15.json")

			var err error
			preUpdateAwsNode, err = rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.AWSNode, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		When("it is out of date", func() {
			It("updates it", func() {
				input.Region = "us-east-1"

				_, err := da.UpdateAWSNode(context.Background(), input, false)
				Expect(err).NotTo(HaveOccurred())

				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
					Equal(fmt.Sprintf("602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon-k8s-cni:%s", expectedVersion)),
				)
				Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
					Equal(fmt.Sprintf("602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon-k8s-cni-init:%s", expectedVersion)),
				)
			})
		})

		When("using a chinese region", func() {
			It("updates it and uses the amazonaws.com.cn address", func() {
				input.Region = "cn-northwest-1"

				_, err := da.UpdateAWSNode(context.Background(), input, false)
				Expect(err).NotTo(HaveOccurred())

				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.AWSNode, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(awsNode.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(awsNode.Spec.Template.Spec.Containers[0].Image).To(
					Equal(fmt.Sprintf("961992271922.dkr.ecr.cn-northwest-1.amazonaws.com.cn/amazon-k8s-cni:%s", expectedVersion)),
				)
				Expect(awsNode.Spec.Template.Spec.InitContainers).To(HaveLen(1))
				Expect(awsNode.Spec.Template.Spec.InitContainers[0].Image).To(
					Equal(fmt.Sprintf("961992271922.dkr.ecr.cn-northwest-1.amazonaws.com.cn/amazon-k8s-cni-init:%s", expectedVersion)),
				)
			})
		})

		When("dry run is true", func() {
			When("it needs an update", func() {
				It("returns true", func() {
					needsUpdate, err := da.UpdateAWSNode(context.Background(), input, true)
					Expect(err).NotTo(HaveOccurred())
					Expect(needsUpdate).To(BeTrue())

					awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.AWSNode, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					//should be unchanged
					Expect(awsNode.Spec).To(Equal(preUpdateAwsNode.Spec))
				})
			})

			When("it doesn't need an update", func() {
				BeforeEach(func() {
					rawClient = testutils.NewFakeRawClient()
					input.RawClient = rawClient
					loadSamples(rawClient, "assets/aws-node.yaml")

					var err error
					preUpdateAwsNode, err = rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.AWSNode, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns false", func() {
					needsUpdate, err := da.UpdateAWSNode(context.Background(), input, true)
					Expect(err).NotTo(HaveOccurred())
					Expect(needsUpdate).To(BeFalse())

					awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.AWSNode, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					//should be unchanged
					Expect(awsNode.Spec).To(Equal(preUpdateAwsNode.Spec))
				})
			})
		})
	})
})

func loadSamples(rawClient *testutils.FakeRawClient, samplesPath string) {
	sampleAddons := testutils.LoadSamples(samplesPath)
	rawClient.AssumeObjectsMissing = true

	for _, item := range sampleAddons {
		rc, err := rawClient.NewRawResource(item)
		Expect(err).NotTo(HaveOccurred())
		_, err = rc.CreateOrReplace(false)
		Expect(err).NotTo(HaveOccurred())
	}
}
