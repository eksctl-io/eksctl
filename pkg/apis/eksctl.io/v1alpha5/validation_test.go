package v1alpha5

import (
	"github.com/bxcodec/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClusterConfig validation", func() {
	Describe("nodeGroups[*].name", func() {
		var (
			cfg *ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "ng0"
			ng1 := cfg.NewNodeGroup()
			ng1.Name = "ng1"
		})

		It("should handle unique nodegroups", func() {
			err = ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())

			for i, ng := range cfg.NodeGroups {
				err = ValidateNodeGroup(i, ng)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should handle non-unique nodegroups", func() {
			cfg.NodeGroups[0].Name = "ng"
			cfg.NodeGroups[1].Name = "ng"

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})

		It("should handle unamed nodegroups", func() {
			cfg.NodeGroups[0].Name = ""

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("nodeGroups[*].iam", func() {
		var (
			cfg *ClusterConfig
			err error
			ng1 *NodeGroup
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()

			ng0 := cfg.NewNodeGroup()
			ng0.Name = "ng0"

			ng0.IAM.AttachPolicyARNs = []string{
				"foo",
				"bar",
			}
			ng0.IAM.WithAddonPolicies.ExternalDNS = Enabled()
			ng0.IAM.WithAddonPolicies.ALBIngress = Enabled()
			ng0.IAM.WithAddonPolicies.ImageBuilder = Enabled()

			ng1 = cfg.NewNodeGroup()
			ng1.Name = "ng1"
		})

		JustBeforeEach(func() {
			err = ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())

			for i, ng := range cfg.NodeGroups {
				err = ValidateNodeGroup(i, ng)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should allow setting only instanceProfileARN", func() {
			ng1.IAM.InstanceProfileARN = "p1"

			err = ValidateNodeGroup(1, ng1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should allow setting only instanceRoleARN", func() {
			ng1.IAM.InstanceRoleARN = "r1"

			err = ValidateNodeGroup(1, ng1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should allow setting instanceProfileARN and instanceRoleARN", func() {
			ng1.IAM.InstanceProfileARN = "p1"
			ng1.IAM.InstanceRoleARN = "r1"

			err = ValidateNodeGroup(1, ng1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not allow setting instanceProfileARN and instanceRoleName", func() {
			ng1.IAM.InstanceProfileARN = "p1"
			ng1.IAM.InstanceRoleName = "aRole"

			err = ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceProfileARN and nodeGroups[1].iam.instanceRoleName cannot be set at the same time"))
		})

		It("should not allow setting instanceProfileARN and instanceRolePermissionsBoundary", func() {
			ng1.IAM.InstanceProfileARN = "p1"
			ng1.IAM.InstanceRolePermissionsBoundary = "aPolicy"

			err = ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceProfileARN and nodeGroups[1].iam.instanceRolePermissionsBoundary cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and instanceRolePermissionsBoundary", func() {
			ng1.IAM.InstanceRoleARN = "r1"
			ng1.IAM.InstanceRolePermissionsBoundary = "p1"

			err = ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.instanceRolePermissionsBoundary cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and instanceRoleName", func() {
			ng1.IAM.InstanceRoleARN = "r1"
			ng1.IAM.InstanceRoleName = "aRole"

			err = ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.instanceRoleName cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and attachPolicyARNs", func() {
			ng1.IAM.InstanceRoleARN = "r1"
			ng1.IAM.AttachPolicyARNs = []string{
				"foo",
				"bar",
			}

			err = ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.attachPolicyARNs cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and withAddonPolicies", func() {
			ng1.IAM.InstanceRoleARN = "r1"

			ng1.IAM.WithAddonPolicies.ExternalDNS = Enabled()

			err = ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.withAddonPolicies.externalDNS cannot be set at the same time"))
		})

	})

	Describe("iam.{withOIDC,serviceAccounts}", func() {
		var (
			cfg *ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		It("should pass when iam.withOIDC is unset", func() {
			cfg.IAM.WithOIDC = nil

			err = ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when iam.withOIDC is disabled", func() {
			cfg.IAM.WithOIDC = Disabled()

			err = ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when iam.withOIDC is enabled", func() {
			cfg.IAM.WithOIDC = Enabled()

			err = ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when iam.withOIDC is disabled and some iam.serviceAccounts are given", func() {
			cfg.IAM.WithOIDC = Disabled()

			cfg.IAM.ServiceAccounts = []*ClusterIAMServiceAccount{{}, {}}
			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[1].Name = "sa-2"

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("iam.withOIDC must be enabled explicitly"))
		})

		It("should pass when iam.withOIDC is enabled and some iam.serviceAccounts are given", func() {
			cfg.IAM.WithOIDC = Enabled()

			cfg.IAM.ServiceAccounts = []*ClusterIAMServiceAccount{{}, {}}

			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[0].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[1].Name = "sa-2"
			cfg.IAM.ServiceAccounts[1].AttachPolicy = map[string]interface{}{"Statement": "foo"}

			err = ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when unnamed iam.serviceAccounts[1] is given", func() {
			cfg.IAM.WithOIDC = Enabled()

			cfg.IAM.ServiceAccounts = []*ClusterIAMServiceAccount{{}, {}}
			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[0].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[1].Name = ""
			cfg.IAM.ServiceAccounts[1].AttachPolicy = map[string]interface{}{"Statement": "foo"}

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("iam.serviceAccounts[1].name must be set"))
		})

		It("should fail when iam.serviceAccounts[1] has no policy", func() {
			cfg.IAM.WithOIDC = Enabled()

			cfg.IAM.ServiceAccounts = []*ClusterIAMServiceAccount{{}, {}}
			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[0].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[1].Name = "sa-2"

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(HavePrefix("iam.serviceAccounts[1].attachPolicyARNs or iam.serviceAccounts[1].attachPolicy must be set"))
		})

		It("should fail when non-uniquely named iam.serviceAccounts are given", func() {
			cfg.IAM.WithOIDC = Enabled()

			cfg.IAM.ServiceAccounts = []*ClusterIAMServiceAccount{{}, {}, {}, {}, {}}
			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[0].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[1].Name = "sa-2"
			cfg.IAM.ServiceAccounts[1].Namespace = "ns-2"
			cfg.IAM.ServiceAccounts[1].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[2].Name = "sa-2"
			cfg.IAM.ServiceAccounts[2].Namespace = "ns-2"
			cfg.IAM.ServiceAccounts[2].AttachPolicy = map[string]interface{}{"Statement": "foo"}

			cfg.IAM.ServiceAccounts[3].Name = "sa-3"
			cfg.IAM.ServiceAccounts[3].AttachPolicy = map[string]interface{}{"Statement": "foo"}

			cfg.IAM.ServiceAccounts[4].Name = "sa-1"
			cfg.IAM.ServiceAccounts[4].AttachPolicy = map[string]interface{}{"Statement": "foo"}

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("<namespace>/<name> of iam.serviceAccounts[2] \"ns-2/sa-2\" is not unique"))

			cfg.IAM.ServiceAccounts[2].Namespace = "ns-3"

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("<namespace>/<name> of iam.serviceAccounts[4] \"/sa-1\" is not unique"))
		})
	})

	Describe("cloudWatch.clusterLogging", func() {
		var (
			cfg *ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		It("should handle known types", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"api"}

			err = ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle unknown types", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"anything"}

			err = ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("cluster endpoint access config", func() {
		var (
			cfg *ClusterConfig
			vpc *ClusterVPC
			err error
		)

		BeforeEach(func() {
			cfg = NewClusterConfig()
			vpc = NewClusterVPC()
			cfg.VPC = vpc
		})

		It("should not error on private=true, public=true", func() {
			cfg.VPC.ClusterEndpoints =
				&ClusterEndpoints{PrivateAccess: Enabled(), PublicAccess: Enabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not error on private=false, public=true", func() {
			cfg.VPC.ClusterEndpoints =
				&ClusterEndpoints{PrivateAccess: Disabled(), PublicAccess: Enabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should error on private=true, public=false", func() {
			cfg.VPC.ClusterEndpoints = &ClusterEndpoints{PrivateAccess: Enabled(), PublicAccess: Disabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).To(BeIdenticalTo(ErrClusterEndpointPrivateOnly))
		})

		It("should error on private=false, public=false", func() {
			cfg.VPC.ClusterEndpoints = &ClusterEndpoints{PrivateAccess: Disabled(), PublicAccess: Disabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).To(BeIdenticalTo(ErrClusterEndpointNoAccess))
		})
	})

	Describe("ssh flags", func() {
		var (
			testKeyPath = "some/path/to/file.pub"
			testKeyName = "id_rsa.pub"
			testKey     = "THIS IS A KEY"
		)

		It("fails when a key path and a key name are specified", func() {
			SSHConfig := &NodeGroupSSH{
				Allow:         Enabled(),
				PublicKeyPath: &testKeyPath,
				PublicKeyName: &testKeyName,
			}

			checkItDetectsError(SSHConfig)
		})

		It("fails when a key path and a key are specified", func() {
			SSHConfig := &NodeGroupSSH{
				Allow:         Enabled(),
				PublicKeyPath: &testKeyPath,
				PublicKey:     &testKey,
			}

			checkItDetectsError(SSHConfig)
		})

		It("fails when a key name and a key are specified", func() {
			SSHConfig := &NodeGroupSSH{
				Allow:         Enabled(),
				PublicKeyName: &testKeyName,
				PublicKey:     &testKey,
			}

			checkItDetectsError(SSHConfig)
		})

		Context("Instances distribution", func() {

			var ng *NodeGroup
			BeforeEach(func() {
				ng = &NodeGroup{
					InstancesDistribution: &NodeGroupInstancesDistribution{
						InstanceTypes: []string{"t3.medium", "t3.large"},
					},
				}
			})

			It("It doesn't panic when instance distribution is not enabled", func() {
				ng.InstancesDistribution = nil
				err := validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when instance distribution is enabled and instanceType is not empty or \"mixed\"", func() {
				err := validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())

				ng.InstanceType = "t3.small"

				err = validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())
			})

			It("It fails when the instance distribution doesn't have at least 2 different instance types", func() {
				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{"t3.medium", "t3.medium"}

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{"t3.medium", "t3.small"}

				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when the onDemandBaseCapacity is not above 0", func() {
				ng.InstancesDistribution.OnDemandBaseCapacity = newInt(-1)

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.OnDemandBaseCapacity = newInt(1)

				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when the spotInstancePools is not between 1 and 20", func() {
				ng.InstancesDistribution.SpotInstancePools = newInt(0)

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.SpotInstancePools = newInt(21)
				err = validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.SpotInstancePools = newInt(2)
				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when the onDemandPercentageAboveBaseCapacity is not between 0 and 100", func() {
				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(-1)

				err := validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(101)
				err = validateInstancesDistribution(ng)
				Expect(err).To(HaveOccurred())

				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(50)
				err = validateInstancesDistribution(ng)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("kubelet extra config", func() {
		Context("Instances distribution", func() {

			var ng *NodeGroup
			BeforeEach(func() {
				ng = &NodeGroup{}
			})

			It("Forbids overriding basic fields", func() {
				testKeys := []string{"kind", "apiVersion", "address", "clusterDomain", "authentication",
					"authorization", "serverTLSBootstrap"}

				for _, key := range testKeys {
					ng.KubeletExtraConfig = &InlineDocument{
						key: "should-not-be-allowed",
					}
					err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig)
					Expect(err).To(HaveOccurred())
				}
			})

			It("Allows other kubelet options", func() {
				ng.KubeletExtraConfig = &InlineDocument{
					"kubeReserved": map[string]string{
						"cpu":               "300m",
						"memory":            "300Mi",
						"ephemeral-storage": "1Gi",
					},
					"kubeReservedCgroup": "/kube-reserved",
					"cgroupDriver":       "systemd",
					"featureGates": map[string]bool{
						"VolumeScheduling":         false,
						"VolumeSnapshotDataSource": true,
					},
				}
				err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig)
				Expect(err).ToNot(HaveOccurred())
			})

		})
	})

	Describe("ebs encryption", func() {
		var (
			nodegroup = "ng1"
			volSize   = 50
			kmsKeyID  = "36c0b54e-64ed-4f2d-a1c7-96558764311e"
			disabled  = false
			enabled   = true
		)

		Context("Encrypted workers", func() {

			var ng *NodeGroup
			BeforeEach(func() {
				ng = &NodeGroup{}
			})

			It("Forbids setting volumeKmsKeyID without volumeEncrypted", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = nil
				ng.VolumeKmsKeyID = &kmsKeyID
				err := ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
			})

			It("Forbids setting volumeKmsKeyID with volumeEncrypted: false", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = &disabled
				ng.VolumeKmsKeyID = &kmsKeyID
				err := ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
			})

			It("Allows setting volumeKmsKeyID with volumeEncrypted: true", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = &enabled
				ng.VolumeKmsKeyID = &kmsKeyID
				err := ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})

		})
	})

	Describe("FargateProfile", func() {
		Describe("Validate", func() {
			It("returns an error when the profile's name is empty", func() {
				profile := FargateProfile{
					Selectors: []FargateProfileSelector{
						FargateProfileSelector{Namespace: "default"},
					},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: empty name"))
			})

			It("returns an error when the profile has a nil selectors array", func() {
				profile := FargateProfile{
					Name: "default",
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"default\": no profile selector"))
			})

			It("returns an error when the profile has an empty selectors array", func() {
				profile := FargateProfile{
					Name:      "default",
					Selectors: []FargateProfileSelector{},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"default\": no profile selector"))
			})

			It("returns an error when the profile's selectors do not have any namespace defined", func() {
				profile := FargateProfile{
					Name: "default",
					Selectors: []FargateProfileSelector{
						FargateProfileSelector{},
					},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"default\": invalid profile selector at index #0: empty namespace"))
			})

			It("returns an error when the profile's name starts with eks-", func() {
				profile := FargateProfile{
					Name: "eks-foo",
					Selectors: []FargateProfileSelector{
						FargateProfileSelector{Namespace: "default"},
					},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"eks-foo\": name should NOT start with \"eks-\""))
			})

			It("passes when a name and at least one selector with a namespace is defined", func() {
				profile := FargateProfile{
					Name: "default",
					Selectors: []FargateProfileSelector{
						FargateProfileSelector{Namespace: "default"},
					},
				}
				err := profile.Validate()
				Expect(err).ToNot(HaveOccurred())
			})

			It("passes on randomly generated fields", func() {
				profile := FargateProfile{}
				err := faker.FakeData(&profile)
				Expect(err).ToNot(HaveOccurred())
				err = profile.Validate()
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

func checkItDetectsError(SSHConfig *NodeGroupSSH) {
	err := validateNodeGroupSSH(SSHConfig)
	Expect(err).To(HaveOccurred())
}

func newInt(value int) *int {
	v := value
	return &v
}
