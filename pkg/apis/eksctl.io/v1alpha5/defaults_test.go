package v1alpha5

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

var _ = Describe("ClusterConfig validation", func() {
	Describe("cloudWatch.clusterLogging", func() {
		var (
			cfg *ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		It("should handle unknown types", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"anything"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})

		It("should have no logging by default", func() {
			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(BeEmpty())
		})

		It("should expand `['*']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"*"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})

		It("should expand `['all']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"all"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})

		It("should expand `['api', '*', 'audit']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"api", "*", "audit"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})

		It("should expand `['authenticator', 'controllermanager', 'all']` to all", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"authenticator", "controllermanager", "all"}

			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.CloudWatch.ClusterLogging.EnableTypes).To(Equal(SupportedCloudWatchClusterLogTypes()))
		})
	})

	Context("SSH settings", func() {

		It("Providing an SSH key enables SSH when SSH.Allow not set", func() {
			testKeyPath := "some/path/to/file.pub"

			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeSize: &DefaultNodeVolumeSize,
					SSH: &NodeGroupSSH{
						PublicKeyPath: &testKeyPath,
					},
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(*testNodeGroup.SSH.Allow).To(BeTrue())
			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo(testKeyPath))
		})

		It("Enabling SSH without a key uses default key", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeSize: &DefaultNodeVolumeSize,
					SSH: &NodeGroupSSH{
						Allow: Enabled(),
					},
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo("~/.ssh/id_rsa.pub"))
		})

		It("Providing an SSH key and explicitly disabling SSH keeps SSH disabled", func() {
			testKeyPath := "some/path/to/file.pub"

			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeSize: &DefaultNodeVolumeSize,
					SSH: &NodeGroupSSH{
						Allow:         Disabled(),
						PublicKeyPath: &testKeyPath,
					},
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(*testNodeGroup.SSH.Allow).To(BeFalse())
			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo(testKeyPath))
		})
	})

	Context("volume settings", func() {
		It("sets up defaults for the main volume", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(*testNodeGroup.VolumeType).To(Equal(DefaultNodeVolumeType))
			Expect(*testNodeGroup.VolumeSize).To(Equal(DefaultNodeVolumeSize))
		})
		It("sets up defaults for any additional volume", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AdditionalVolumes: []*VolumeMapping{
						{
							VolumeName: aws.String("test"),
						},
					},
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(*testNodeGroup.AdditionalVolumes[0].VolumeType).To(Equal(DefaultNodeVolumeType))
			Expect(*testNodeGroup.AdditionalVolumes[0].VolumeSize).To(Equal(DefaultNodeVolumeSize))
		})
		It("sets up defaults for gp3", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeType: aws.String(NodeVolumeTypeGP3),
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(*testNodeGroup.VolumeType).To(Equal(NodeVolumeTypeGP3))
			Expect(*testNodeGroup.VolumeIOPS).To(Equal(DefaultNodeVolumeGP3IOPS))
			Expect(*testNodeGroup.VolumeThroughput).To(Equal(DefaultNodeVolumeThroughput))
		})

		It("[Outposts] sets up defaults for the main volume", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(*testNodeGroup.VolumeType).To(Equal(NodeVolumeTypeGP2))
			Expect(*testNodeGroup.VolumeSize).To(Equal(DefaultNodeVolumeSize))
		})

		It("[Outposts] sets up defaults for additional volumes", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AdditionalVolumes: []*VolumeMapping{
						{
							VolumeName: aws.String("test"),
						},
					},
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(*testNodeGroup.AdditionalVolumes[0].VolumeType).To(Equal(NodeVolumeTypeGP2))
			Expect(*testNodeGroup.AdditionalVolumes[0].VolumeSize).To(Equal(DefaultNodeVolumeSize))
		})
	})

	Context("Bottlerocket Settings", func() {
		It("enables SSH with NodeGroup", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
					SSH: &NodeGroupSSH{
						Allow: Enabled(),
					},
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(*testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeTrue())
		})

		It("leaves EnableAdminContainer unset if SSH is disabled", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
					SSH: &NodeGroupSSH{
						Allow: Disabled(),
					},
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeNil())
		})

		It("has default NodeGroup configuration", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(testNodeGroup.Bottlerocket).NotTo(BeNil())
			Expect(testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeNil())
		})

		It("tolerates non standard casing of AMI Family", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: "BoTTleRocKet",
				},
			}

			err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(testNodeGroup.NodeGroupBase.AMIFamily).To(Equal(NodeImageFamilyBottlerocket))
		})
	})

	Context("Cluster NAT settings", func() {

		It("Cluster NAT defaults to single NAT gateway mode", func() {
			testVpc := &ClusterVPC{}
			testVpc.NAT = DefaultClusterNAT()

			Expect(*testVpc.NAT.Gateway).To(BeIdenticalTo(ClusterSingleNAT))

		})

	})

	Context("Container Runtime settings", func() {
		Context("Kubernetes version 1.23 or lower", func() {
			When("ami family is windows", func() {
				It("defaults to docker as a container runtime", func() {
					testNodeGroup := NodeGroup{
						NodeGroupBase: &NodeGroupBase{
							AMIFamily: NodeImageFamilyWindowsServer2019CoreContainer,
						},
					}
					err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
					Expect(err).NotTo(HaveOccurred())
					Expect(*testNodeGroup.ContainerRuntime).To(Equal(ContainerRuntimeDockerForWindows))
				})
			})
			When("ami family is AmazonLinux2023", func() {
				It("defaults to containerd as a container runtime", func() {
					testNodeGroup := NodeGroup{
						NodeGroupBase: &NodeGroupBase{
							AMIFamily: NodeImageFamilyAmazonLinux2023,
						},
					}
					err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Version: Version1_23}, false)
					Expect(err).NotTo(HaveOccurred())
					Expect(*testNodeGroup.ContainerRuntime).To(Equal(ContainerRuntimeContainerD))
				})
			})
		})

		Context("Kubernetes version 1.24 or greater", func() {
			It("defaults to containerd as a container runtime", func() {
				testNodeGroup := NodeGroup{
					NodeGroupBase: &NodeGroupBase{},
				}
				err := SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Version: DockershimDeprecationVersion}, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(*testNodeGroup.ContainerRuntime).To(Equal(ContainerRuntimeContainerD))
			})
		})
	})

	Describe("Cluster Managed Shared Node Security Group settings", func() {
		var (
			cfg *ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		It("should be enabled by default", func() {
			SetClusterConfigDefaults(cfg)
			Expect(*cfg.VPC.ManageSharedNodeSecurityGroupRules).To(BeTrue())
		})

		It("should fail validation if disabled without a defined shared node security group", func() {
			cfg.VPC.ManageSharedNodeSecurityGroupRules = Disabled()
			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})

		It("should pass validation if disabled with a defined shared node security group", func() {
			cfg.VPC.SharedNodeSecurityGroup = "sg-0123456789"
			cfg.VPC.ManageSharedNodeSecurityGroupRules = Disabled()
			SetClusterConfigDefaults(cfg)
			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Context("Authentication Mode", func() {
		var (
			cfg *ClusterConfig
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		It("should be set to API_AND_CONFIG_MAP by default", func() {
			SetClusterConfigDefaults(cfg)
			Expect(cfg.AccessConfig.AuthenticationMode).To(Equal(ekstypes.AuthenticationModeApiAndConfigMap))
		})

		It("should be set to CONFIG_MAP when control plane is on outposts", func() {
			cfg.Outpost = &Outpost{
				ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
			SetClusterConfigDefaults(cfg)
			Expect(cfg.AccessConfig.AuthenticationMode).To(Equal(ekstypes.AuthenticationModeConfigMap))
		})
	})

	Describe("ClusterConfig", func() {
		var cfg *ClusterConfig

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		Describe("RemoteNetworkConfig", func() {
			It("should set default credentials provider to SSM", func() {
				cfg.RemoteNetworkConfig = &RemoteNetworkConfig{}
				SetClusterConfigDefaults(cfg)
				Expect(*cfg.RemoteNetworkConfig.IAM.Provider).To(Equal(SSMProvider))
			})
		})

		Describe("SetDefaultFargateProfile", func() {
			It("should create a default Fargate profile with two selectors matching default and kube-system w/o any label", func() {
				Expect(cfg.FargateProfiles).To(HaveLen(0))
				cfg.SetDefaultFargateProfile()
				Expect(cfg.FargateProfiles).To(HaveLen(1))
				profile := cfg.FargateProfiles[0]
				Expect(profile.Name).To(Equal("fp-default"))
				Expect(profile.Selectors).To(HaveLen(2))
				Expect(profile.Selectors[0].Namespace).To(Equal("default"))
				Expect(profile.Selectors[0].Labels).To(HaveLen(0))
				Expect(profile.Selectors[1].Namespace).To(Equal("kube-system"))
				Expect(profile.Selectors[1].Labels).To(HaveLen(0))
			})
		})

		DescribeTable("default AMI family", func(kubernetesVersion, expectedAMIFamily string) {
			mng := NewManagedNodeGroup()
			err := SetManagedNodeGroupDefaults(mng, &ClusterMeta{
				Version: kubernetesVersion,
			}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(mng.AMIFamily).To(Equal(expectedAMIFamily))
		},
			Entry("EKS 1.33 uses AL2023", "1.33", NodeImageFamilyAmazonLinux2023),
			Entry("EKS 1.32 uses AL2023", "1.32", NodeImageFamilyAmazonLinux2023),
			Entry("EKS 1.31 uses AL2023", "1.31", NodeImageFamilyAmazonLinux2023),
			Entry("EKS 1.30 uses AL2023", "1.30", NodeImageFamilyAmazonLinux2023),
			Entry("EKS 1.29 uses AL2", "1.29", NodeImageFamilyAmazonLinux2023),
			Entry("EKS 1.28 uses AL2", "1.28", NodeImageFamilyAmazonLinux2023),
		)
	})

	Context("Node Group AMI family changes for K8s version 1.33", func() {
		It("Should return an error when AL2 AMI family was explicitly selected for K8s version 1.33", func() {
			ng := &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyAmazonLinux2,
				},
			}
			err := SetNodeGroupDefaults(ng, &ClusterMeta{Version: "1.33"}, false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("AmazonLinux2 is not supported for Kubernetes version 1.33"))
		})

		It("AL2 AMI family should remain unchanged for K8s version 1.32", func() {
			ng := &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyAmazonLinux2,
				},
			}
			err := SetNodeGroupDefaults(ng, &ClusterMeta{Version: "1.32"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyAmazonLinux2))
		})

		It("Non-AL2 AMI families should remain unchanged for K8s version 1.33", func() {
			ng := &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
				},
			}
			err := SetNodeGroupDefaults(ng, &ClusterMeta{Version: "1.33"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyBottlerocket))
		})

		It("Default AMI family should be AL2023 for K8s version 1.33", func() {
			ng := &NodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}
			err := SetNodeGroupDefaults(ng, &ClusterMeta{Version: "1.33"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyAmazonLinux2023))
		})

		It("Default AMI family should be AL2 for K8s version 1.32", func() {
			ng := &NodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}
			err := SetNodeGroupDefaults(ng, &ClusterMeta{Version: "1.32"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyAmazonLinux2023))
		})
	})

	Context("Managed Node Group AMI family changes for K8s version 1.33", func() {
		It("Should return an error when AL2 AMI family was explicitly selected for K8s version 1.33", func() {
			ng := &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyAmazonLinux2,
				},
			}
			err := SetManagedNodeGroupDefaults(ng, &ClusterMeta{Version: "1.33"}, false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("AmazonLinux2 is not supported for Kubernetes version 1.33"))
		})

		It("AL2 AMI family should remain unchanged for K8s version 1.32", func() {
			ng := &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyAmazonLinux2,
				},
			}
			err := SetManagedNodeGroupDefaults(ng, &ClusterMeta{Version: "1.32"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyAmazonLinux2))
		})

		It("Non-AL2 AMI families should remain unchanged for K8s version 1.33", func() {
			ng := &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
				},
			}
			err := SetManagedNodeGroupDefaults(ng, &ClusterMeta{Version: "1.33"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyBottlerocket))
		})

		It("Default AMI family should be AL2023 for K8s version 1.33", func() {
			ng := &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}
			err := SetManagedNodeGroupDefaults(ng, &ClusterMeta{Version: "1.33"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyAmazonLinux2023))
		})

		It("Default AMI family should be AL2 for K8s version 1.29", func() {
			ng := &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}
			err := SetManagedNodeGroupDefaults(ng, &ClusterMeta{Version: "1.29"}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(ng.AMIFamily).To(Equal(NodeImageFamilyAmazonLinux2023))
		})
	})
})
