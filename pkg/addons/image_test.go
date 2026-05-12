package addons_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"

	"github.com/weaveworks/eksctl/pkg/addons"
)

var _ = Describe("UseRegionalImage", func() {
	var spec *corev1.PodTemplateSpec

	Context("with a single container and single init container", func() {
		BeforeEach(func() {
			spec = &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "main",
							Image: "%s.dkr.ecr.%s.%s/amazon-k8s-cni:v1.21.1",
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:  "init",
							Image: "%s.dkr.ecr.%s.%s/amazon-k8s-cni-init:v1.21.1",
						},
					},
				},
			}
		})

		It("regionalizes both containers to the specified region", func() {
			err := addons.UseRegionalImage(spec, "eu-west-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.Spec.Containers[0].Image).To(Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni:v1.21.1"))
			Expect(spec.Spec.InitContainers[0].Image).To(Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni-init:v1.21.1"))
		})
	})

	Context("with multiple containers including aws-eks-nodeagent", func() {
		BeforeEach(func() {
			spec = &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "aws-node",
							Image: "%s.dkr.ecr.%s.%s/amazon-k8s-cni:v1.21.1",
						},
						{
							Name:  "aws-eks-nodeagent",
							Image: "%s.dkr.ecr.%s.%s/amazon/aws-network-policy-agent:v1.3.1",
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:  "aws-vpc-cni-init",
							Image: "%s.dkr.ecr.%s.%s/amazon-k8s-cni-init:v1.21.1",
						},
					},
				},
			}
		})

		It("regionalizes all containers to the specified region", func() {
			err := addons.UseRegionalImage(spec, "eu-central-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.Spec.Containers[0].Image).To(Equal("602401143452.dkr.ecr.eu-central-1.amazonaws.com/amazon-k8s-cni:v1.21.1"))
			Expect(spec.Spec.Containers[1].Image).To(Equal("602401143452.dkr.ecr.eu-central-1.amazonaws.com/amazon/aws-network-policy-agent:v1.3.1"))
			Expect(spec.Spec.InitContainers[0].Image).To(Equal("602401143452.dkr.ecr.eu-central-1.amazonaws.com/amazon-k8s-cni-init:v1.21.1"))
		})

		It("regionalizes to a Chinese region with correct DNS suffix", func() {
			err := addons.UseRegionalImage(spec, "cn-northwest-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.Spec.Containers[0].Image).To(Equal("961992271922.dkr.ecr.cn-northwest-1.amazonaws.com.cn/amazon-k8s-cni:v1.21.1"))
			Expect(spec.Spec.Containers[1].Image).To(Equal("961992271922.dkr.ecr.cn-northwest-1.amazonaws.com.cn/amazon/aws-network-policy-agent:v1.3.1"))
			Expect(spec.Spec.InitContainers[0].Image).To(Equal("961992271922.dkr.ecr.cn-northwest-1.amazonaws.com.cn/amazon-k8s-cni-init:v1.21.1"))
		})
	})

	Context("with a mix of format-string and already-resolved images", func() {
		BeforeEach(func() {
			spec = &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "aws-node",
							Image: "%s.dkr.ecr.%s.%s/amazon-k8s-cni:v1.21.1",
						},
						{
							Name:  "sidecar",
							Image: "public.ecr.aws/some-sidecar:latest",
						},
					},
				},
			}
		})

		It("only regionalizes images that are in format-string format", func() {
			err := addons.UseRegionalImage(spec, "ap-southeast-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.Spec.Containers[0].Image).To(Equal("602401143452.dkr.ecr.ap-southeast-1.amazonaws.com/amazon-k8s-cni:v1.21.1"))
			Expect(spec.Spec.Containers[1].Image).To(Equal("public.ecr.aws/some-sidecar:latest"))
		})
	})

	Context("with no init containers", func() {
		BeforeEach(func() {
			spec = &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "aws-node",
							Image: "%s.dkr.ecr.%s.%s/amazon-k8s-cni:v1.21.1",
						},
					},
				},
			}
		})

		It("does not error when there are no init containers", func() {
			err := addons.UseRegionalImage(spec, "us-east-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.Spec.Containers[0].Image).To(
				Equal("602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon-k8s-cni:v1.21.1"),
			)
		})
	})
})
