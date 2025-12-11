package addons_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("NvidiaDevicePlugin", func() {
	Describe("SetTolerations", func() {
		var (
			plugin *addons.NvidiaDevicePlugin
			spec   *corev1.PodTemplateSpec
			config *api.ClusterConfig
		)

		BeforeEach(func() {
			spec = &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Tolerations: []corev1.Toleration{},
				},
			}
			config = &api.ClusterConfig{}
		})

		Context("with NodeGroups", func() {
			It("should add tolerations for AmazonLinux2 nodegroups with NVIDIA instances", func() {
				config.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "nvidia-ng",
							InstanceType: "g4dn.xlarge",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "nvidia.com/gpu", Value: "true", Effect: "NoSchedule"},
							{Key: "workload", Value: "gpu", Effect: "NoExecute"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(HaveLen(2))
				Expect(spec.Spec.Tolerations).To(ContainElement(corev1.Toleration{
					Key:   "nvidia.com/gpu",
					Value: "true",
				}))
				Expect(spec.Spec.Tolerations).To(ContainElement(corev1.Toleration{
					Key:   "workload",
					Value: "gpu",
				}))
			})

			It("should add tolerations for AmazonLinux2023 nodegroups with NVIDIA instances", func() {
				config.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "nvidia-ng",
							InstanceType: "g5.xlarge",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2023,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "gpu-workload", Value: "ml", Effect: "NoSchedule"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(HaveLen(1))
				Expect(spec.Spec.Tolerations[0].Key).To(Equal("gpu-workload"))
				Expect(spec.Spec.Tolerations[0].Value).To(Equal("ml"))
			})

			It("should not add tolerations for non-NVIDIA instances", func() {
				config.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "cpu-ng",
							InstanceType: "m5.large",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "cpu-only", Value: "true", Effect: "NoSchedule"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(BeEmpty())
			})

			It("should not add tolerations for unsupported AMI families", func() {
				config.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "nvidia-ng",
							InstanceType: "g4dn.xlarge",
							AMIFamily:    api.NodeImageFamilyUbuntu2004,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "nvidia.com/gpu", Value: "true", Effect: "NoSchedule"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(BeEmpty())
			})
		})

		Context("with ManagedNodeGroups", func() {
			It("should add tolerations for AmazonLinux2 managed nodegroups with NVIDIA instances", func() {
				config.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "managed-nvidia-ng",
							InstanceType: "g4dn.2xlarge",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "managed-gpu", Value: "nvidia", Effect: "NoSchedule"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(HaveLen(1))
				Expect(spec.Spec.Tolerations[0].Key).To(Equal("managed-gpu"))
				Expect(spec.Spec.Tolerations[0].Value).To(Equal("nvidia"))
			})

			It("should add tolerations for AmazonLinux2023 managed nodegroups with NVIDIA instances", func() {
				config.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "managed-nvidia-ng",
							InstanceType: "g5.4xlarge",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2023,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "ml-workload", Value: "training", Effect: "NoExecute"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(HaveLen(1))
				Expect(spec.Spec.Tolerations[0].Key).To(Equal("ml-workload"))
				Expect(spec.Spec.Tolerations[0].Value).To(Equal("training"))
			})
		})

		Context("with existing tolerations", func() {
			It("should not duplicate existing tolerations", func() {
				spec.Spec.Tolerations = []corev1.Toleration{
					{Key: "existing-taint", Value: "existing-value"},
				}

				config.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "nvidia-ng",
							InstanceType: "g4dn.xlarge",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "existing-taint", Value: "different-value", Effect: "NoSchedule"},
							{Key: "new-taint", Value: "new-value", Effect: "NoSchedule"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(HaveLen(2))
				Expect(spec.Spec.Tolerations).To(ContainElement(corev1.Toleration{
					Key:   "existing-taint",
					Value: "existing-value",
				}))
				Expect(spec.Spec.Tolerations).To(ContainElement(corev1.Toleration{
					Key:   "new-taint",
					Value: "new-value",
				}))
			})
		})

		Context("with mixed nodegroup types", func() {
			It("should combine taints from both regular and managed nodegroups", func() {
				config.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "nvidia-ng",
							InstanceType: "g4dn.xlarge",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "regular-gpu", Value: "nvidia", Effect: "NoSchedule"},
						},
					},
				}
				config.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "managed-nvidia-ng",
							InstanceType: "g5.xlarge",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2023,
						},
						Taints: []api.NodeGroupTaint{
							{Key: "managed-gpu", Value: "nvidia", Effect: "NoSchedule"},
						},
					},
				}

				plugin = addons.NewNvidiaDevicePlugin(nil, "us-west-2", false, config).(*addons.NvidiaDevicePlugin)
				err := plugin.SetTolerations(spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Spec.Tolerations).To(HaveLen(2))
				Expect(spec.Spec.Tolerations).To(ContainElement(corev1.Toleration{
					Key:   "regular-gpu",
					Value: "nvidia",
				}))
				Expect(spec.Spec.Tolerations).To(ContainElement(corev1.Toleration{
					Key:   "managed-gpu",
					Value: "nvidia",
				}))
			})
		})
	})
})
