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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)

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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)

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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)

			Expect(*testNodeGroup.SSH.Allow).To(BeFalse())
			Expect(*testNodeGroup.SSH.PublicKeyPath).To(BeIdenticalTo(testKeyPath))
		})
	})

	Context("volume settings", func() {
		It("sets up defaults for the main volume", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(*testNodeGroup.AdditionalVolumes[0].VolumeType).To(Equal(DefaultNodeVolumeType))
			Expect(*testNodeGroup.AdditionalVolumes[0].VolumeSize).To(Equal(DefaultNodeVolumeSize))
		})
		It("sets up defaults for gp3", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					VolumeType: aws.String(NodeVolumeTypeGP3),
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(*testNodeGroup.VolumeType).To(Equal(NodeVolumeTypeGP3))
			Expect(*testNodeGroup.VolumeIOPS).To(Equal(DefaultNodeVolumeGP3IOPS))
			Expect(*testNodeGroup.VolumeThroughput).To(Equal(DefaultNodeVolumeThroughput))
		})

		It("[Outposts] sets up defaults for the main volume", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, true)
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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, true)
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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)

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

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
			Expect(testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeNil())
		})

		It("has default NodeGroup configuration", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: NodeImageFamilyBottlerocket,
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)

			Expect(testNodeGroup.Bottlerocket).NotTo(BeNil())
			Expect(testNodeGroup.Bottlerocket.EnableAdminContainer).To(BeNil())
		})

		It("tolerates non standard casing of AMI Family", func() {
			testNodeGroup := NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: "BoTTleRocKet",
				},
			}

			SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)

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
			It("defaults to dockerd as a container runtime", func() {
				testNodeGroup := NodeGroup{
					NodeGroupBase: &NodeGroupBase{},
				}
				SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Version: Version1_23}, false)
				Expect(*testNodeGroup.ContainerRuntime).To(Equal(ContainerRuntimeDockerD))
			})
			When("ami family is windows", func() {
				It("defaults to docker as a container runtime", func() {
					testNodeGroup := NodeGroup{
						NodeGroupBase: &NodeGroupBase{
							AMIFamily: NodeImageFamilyWindowsServer2019CoreContainer,
						},
					}
					SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{}, false)
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
					SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Version: Version1_23}, false)
					Expect(*testNodeGroup.ContainerRuntime).To(Equal(ContainerRuntimeContainerD))
				})
			})
		})

		Context("Kubernetes version 1.24 or greater", func() {
			It("defaults to containerd as a container runtime", func() {
				testNodeGroup := NodeGroup{
					NodeGroupBase: &NodeGroupBase{},
				}
				SetNodeGroupDefaults(&testNodeGroup, &ClusterMeta{Version: Version1_24}, false)
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
	})
})
