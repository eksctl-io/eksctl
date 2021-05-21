package v1alpha5_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

var _ = Describe("ClusterConfig validation", func() {
	Describe("nodeGroups[*].name", func() {
		var (
			cfg *api.ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "ng0"
			ng1 := cfg.NewNodeGroup()
			ng1.Name = "ng1"
		})

		It("should handle unique nodegroups", func() {
			err = api.ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())

			for i, ng := range cfg.NodeGroups {
				err = api.ValidateNodeGroup(i, ng)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should handle non-unique nodegroups", func() {
			cfg.NodeGroups[0].Name = "ng"
			cfg.NodeGroups[1].Name = "ng"

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})

		It("should handle unamed nodegroups", func() {
			cfg.NodeGroups[0].Name = ""

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("nodeGroups[*].volumeX", func() {
		var (
			cfg *api.ClusterConfig
			ng0 *api.NodeGroup
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			ng0 = cfg.NewNodeGroup()
			ng0.Name = "ng0"
		})

		When("volumeIOPS is set", func() {
			BeforeEach(func() {
				ng0.VolumeIOPS = aws.Int(3000)
			})

			When("VolumeType is gp3", func() {
				BeforeEach(func() {
					*ng0.VolumeType = api.NodeVolumeTypeGP3
				})

				It("does not fail", func() {
					Expect(api.ValidateNodeGroup(0, ng0)).To(Succeed())
				})

				When(fmt.Sprintf("the value of volumeIOPS is < %d", api.MinGP3Iops), func() {
					It("returns an error", func() {
						ng0.VolumeIOPS = aws.Int(api.MinGP3Iops - 1)
						Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("value for nodeGroups[0].volumeIOPS must be within range 3000-16000"))
					})
				})

				When(fmt.Sprintf("the value of volumeIOPS is > %d", api.MaxGP3Iops), func() {
					It("returns an error", func() {
						ng0.VolumeIOPS = aws.Int(api.MaxGP3Iops + 1)
						Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("value for nodeGroups[0].volumeIOPS must be within range 3000-16000"))
					})
				})
			})

			When("VolumeType is io1", func() {
				BeforeEach(func() {
					*ng0.VolumeType = api.NodeVolumeTypeIO1
				})

				It("does not fail", func() {
					Expect(api.ValidateNodeGroup(0, ng0)).To(Succeed())
				})

				When(fmt.Sprintf("the value of volumeIOPS is < %d", api.MinIO1Iops), func() {
					It("returns an error", func() {
						ng0.VolumeIOPS = aws.Int(api.MinIO1Iops - 1)
						Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("value for nodeGroups[0].volumeIOPS must be within range 100-64000"))
					})
				})

				When(fmt.Sprintf("the value of volumeIOPS is > %d", api.MaxIO1Iops), func() {
					It("returns an error", func() {
						ng0.VolumeIOPS = aws.Int(api.MaxIO1Iops + 1)
						Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("value for nodeGroups[0].volumeIOPS must be within range 100-64000"))
					})
				})
			})

			When("VolumeType is one for which IOPS is not supported", func() {
				It("returns an error", func() {
					*ng0.VolumeType = api.NodeVolumeTypeGP2
					Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("nodeGroups[0].volumeIOPS is only supported for io1 and gp3 volume types"))
				})
			})
		})

		When("volumeThroughput is set", func() {
			BeforeEach(func() {
				ng0.VolumeThroughput = aws.Int(125)
			})

			When("VolumeType is gp3", func() {
				BeforeEach(func() {
					*ng0.VolumeType = api.NodeVolumeTypeGP3
				})

				It("does not fail", func() {
					Expect(api.ValidateNodeGroup(0, ng0)).To(Succeed())
				})

				When(fmt.Sprintf("the value of volumeThroughput is < %d", api.MinThroughput), func() {
					It("returns an error", func() {
						ng0.VolumeThroughput = aws.Int(api.MinThroughput - 1)
						Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("value for nodeGroups[0].volumeThroughput must be within range 125-1000"))
					})
				})

				When(fmt.Sprintf("the value of volumeIOPS is > %d", api.MaxThroughput), func() {
					It("returns an error", func() {
						ng0.VolumeThroughput = aws.Int(api.MaxThroughput + 1)
						Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("value for nodeGroups[0].volumeThroughput must be within range 125-1000"))
					})
				})
			})

			When("VolumeType is one for which Throughput is not supported", func() {
				It("returns an error", func() {
					*ng0.VolumeType = api.NodeVolumeTypeGP2
					Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError("nodeGroups[0].volumeThroughput is only supported for gp3 volume type"))
				})
			})
		})
	})

	Describe("nodeGroups[*].iam", func() {
		var (
			cfg *api.ClusterConfig
			err error
			ng1 *api.NodeGroup
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()

			ng0 := cfg.NewNodeGroup()
			ng0.Name = "ng0"

			ng0.IAM.AttachPolicyARNs = []string{
				"arn:aws:iam::aws:policy/Foo",
				"arn:aws:iam::aws:policy/Bar",
			}
			ng0.IAM.WithAddonPolicies.ExternalDNS = api.Enabled()
			ng0.IAM.WithAddonPolicies.AWSLoadBalancerController = api.Enabled()
			ng0.IAM.WithAddonPolicies.ImageBuilder = api.Enabled()

			ng1 = cfg.NewNodeGroup()
			ng1.Name = "ng1"
		})

		JustBeforeEach(func() {
			err = api.ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())

			for i, ng := range cfg.NodeGroups {
				err = api.ValidateNodeGroup(i, ng)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should allow setting only instanceProfileARN", func() {
			ng1.IAM.InstanceProfileARN = "p1"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should allow setting only instanceRoleARN", func() {
			ng1.IAM.InstanceRoleARN = "r1"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should allow setting instanceProfileARN and instanceRoleARN", func() {
			ng1.IAM.InstanceProfileARN = "p1"
			ng1.IAM.InstanceRoleARN = "r1"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not allow setting instanceProfileARN and instanceRoleName", func() {
			ng1.IAM.InstanceProfileARN = "p1"
			ng1.IAM.InstanceRoleName = "aRole"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceProfileARN and nodeGroups[1].iam.instanceRoleName cannot be set at the same time"))
		})

		It("should not allow setting instanceProfileARN and instanceRolePermissionsBoundary", func() {
			ng1.IAM.InstanceProfileARN = "p1"
			ng1.IAM.InstanceRolePermissionsBoundary = "aPolicy"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceProfileARN and nodeGroups[1].iam.instanceRolePermissionsBoundary cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and instanceRolePermissionsBoundary", func() {
			ng1.IAM.InstanceRoleARN = "r1"
			ng1.IAM.InstanceRolePermissionsBoundary = "p1"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.instanceRolePermissionsBoundary cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and instanceRoleName", func() {
			ng1.IAM.InstanceRoleARN = "r1"
			ng1.IAM.InstanceRoleName = "aRole"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.instanceRoleName cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and attachPolicyARNs", func() {
			ng1.IAM.InstanceRoleARN = "r1"
			ng1.IAM.AttachPolicyARNs = []string{
				"arn:aws:iam::aws:policy/Foo",
				"arn:aws:iam::aws:policy/Bar",
			}

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.attachPolicyARNs cannot be set at the same time"))
		})

		It("should not allow setting instanceRoleARN and withAddonPolicies", func() {
			ng1.IAM.InstanceRoleARN = "r1"

			ng1.IAM.WithAddonPolicies.ExternalDNS = api.Enabled()

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.withAddonPolicies.externalDNS cannot be set at the same time"))
		})

	})

	Describe("iam.{withOIDC,serviceAccounts}", func() {
		var (
			cfg *api.ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
		})

		It("should pass when iam.withOIDC is unset", func() {
			cfg.IAM.WithOIDC = nil

			err = api.ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when iam.withOIDC is disabled", func() {
			cfg.IAM.WithOIDC = api.Disabled()

			err = api.ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when iam.withOIDC is enabled", func() {
			cfg.IAM.WithOIDC = api.Enabled()

			err = api.ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when iam.withOIDC is disabled and some iam.serviceAccounts are given", func() {
			cfg.IAM.WithOIDC = api.Disabled()

			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{{}, {}}
			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[1].Name = "sa-2"

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("iam.withOIDC must be enabled explicitly"))
		})

		It("should pass when iam.withOIDC is enabled and some iam.serviceAccounts are given", func() {
			cfg.IAM.WithOIDC = api.Enabled()

			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{{}, {}}

			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[0].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[1].Name = "sa-2"
			cfg.IAM.ServiceAccounts[1].AttachPolicy = map[string]interface{}{"Statement": "foo"}

			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when unnamed iam.serviceAccounts[1] is given", func() {
			cfg.IAM.WithOIDC = api.Enabled()

			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{{}, {}}
			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[0].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[1].Name = ""
			cfg.IAM.ServiceAccounts[1].AttachPolicy = map[string]interface{}{"Statement": "foo"}

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("iam.serviceAccounts[1].name must be set"))
		})

		It("should fail when iam.serviceAccounts[1] has no policy", func() {
			cfg.IAM.WithOIDC = api.Enabled()

			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{{}, {}}
			cfg.IAM.ServiceAccounts[0].Name = "sa-1"
			cfg.IAM.ServiceAccounts[0].AttachPolicyARNs = []string{""}

			cfg.IAM.ServiceAccounts[1].Name = "sa-2"

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(SatisfyAll(
				ContainSubstring("iam.serviceAccounts[1]"),
				ContainSubstring("attachPolicy"),
				ContainSubstring("must be set"),
			))
		})

		It("should fail when non-uniquely named iam.serviceAccounts are given", func() {
			cfg.IAM.WithOIDC = api.Enabled()

			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{{}, {}, {}, {}, {}}
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

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("<namespace>/<name> of iam.serviceAccounts[2] \"ns-2/sa-2\" is not unique"))

			cfg.IAM.ServiceAccounts[2].Namespace = "ns-3"

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("<namespace>/<name> of iam.serviceAccounts[4] \"/sa-1\" is not unique"))
		})
	})

	Describe("cloudWatch.clusterLogging", func() {
		var (
			cfg *api.ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
		})

		It("should handle known types", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"api"}

			err = api.ValidateClusterConfig(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle unknown types", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"anything"}

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("cluster endpoint access config", func() {
		var (
			cfg *api.ClusterConfig
			vpc *api.ClusterVPC
			err error
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			vpc = api.NewClusterVPC()
			cfg.VPC = vpc
		})

		It("should not error on private=true, public=true", func() {
			cfg.VPC.ClusterEndpoints =
				&api.ClusterEndpoints{PrivateAccess: api.Enabled(), PublicAccess: api.Enabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not error on private=false, public=true", func() {
			cfg.VPC.ClusterEndpoints =
				&api.ClusterEndpoints{PrivateAccess: api.Disabled(), PublicAccess: api.Enabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not error on private=true, public=false", func() {
			cfg.VPC.ClusterEndpoints =
				&api.ClusterEndpoints{PrivateAccess: api.Enabled(), PublicAccess: api.Disabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should error on private=false, public=false", func() {
			cfg.VPC.ClusterEndpoints = &api.ClusterEndpoints{PrivateAccess: api.Disabled(), PublicAccess: api.Disabled()}
			err = cfg.ValidateClusterEndpointConfig()
			Expect(err).To(BeIdenticalTo(api.ErrClusterEndpointNoAccess))
		})
	})

	Describe("cpuCredits", func() {
		var ng *api.NodeGroup
		BeforeEach(func() {
			unlimited := "unlimited"
			ng = &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{},
				InstancesDistribution: &api.NodeGroupInstancesDistribution{
					InstanceTypes: []string{"t3.medium", "t3.large"},
				},
				CPUCredits: &unlimited,
			}
		})

		It("works independent of instanceType", func() {
			Context("unset", func() {
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})
			Context("set", func() {
				ng.InstanceType = "mixed"
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		It("errors if no instance distribution", func() {
			ng.InstancesDistribution = nil
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).To(HaveOccurred())
		})

		It("errors if no instance types", func() {
			ng.InstancesDistribution.InstanceTypes = []string{}
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ssh flags", func() {
		var (
			testKeyPath = "some/path/to/file.pub"
			testKeyName = "id_rsa.pub"
			testKey     = "THIS IS A KEY"
			ng          *api.NodeGroup
		)

		BeforeEach(func() {
			ng = newNodeGroup()
		})

		It("fails when a key path and a key name are specified", func() {
			SSHConfig := &api.NodeGroupSSH{
				Allow:         api.Enabled(),
				PublicKeyPath: &testKeyPath,
				PublicKeyName: &testKeyName,
			}

			ng.SSH = SSHConfig
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).To(MatchError("only one of publicKeyName, publicKeyPath or publicKey can be specified for SSH per node-group"))
		})

		It("fails when a key path and a key are specified", func() {
			SSHConfig := &api.NodeGroupSSH{
				Allow:         api.Enabled(),
				PublicKeyPath: &testKeyPath,
				PublicKey:     &testKey,
			}

			ng.SSH = SSHConfig
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).To(MatchError("only one of publicKeyName, publicKeyPath or publicKey can be specified for SSH per node-group"))
		})

		It("fails when a key name and a key are specified", func() {
			SSHConfig := &api.NodeGroupSSH{
				Allow:         api.Enabled(),
				PublicKeyName: &testKeyName,
				PublicKey:     &testKey,
			}

			ng.SSH = SSHConfig
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).To(MatchError("only one of publicKeyName, publicKeyPath or publicKey can be specified for SSH per node-group"))
		})

		Context("Instances distribution", func() {
			var ng *api.NodeGroup
			BeforeEach(func() {
				ng = &api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						InstanceTypes:                       []string{"t3.medium", "t3.large"},
						OnDemandBaseCapacity:                newInt(1),
						SpotInstancePools:                   newInt(2),
						OnDemandPercentageAboveBaseCapacity: newInt(50),
					},
				}
			})

			It("It doesn't panic when instance distribution is not enabled", func() {
				ng.InstancesDistribution = nil
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It doesn't fail when instance distribution is enabled and instanceType is \"mixed\"", func() {
				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{"t3.medium"}

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})

			It("It fails when instance distribution is enabled and instanceType set", func() {
				ng.InstanceType = "t3.small"

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("instanceType should be \"mixed\" or unset when using the instances distribution feature"))
			})

			It("It fails when the instance distribution doesn't have any instance type", func() {
				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{}

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("at least two instance types have to be specified for mixed nodegroups"))
			})

			It("It fails when the onDemandBaseCapacity is not above 0", func() {
				ng.InstancesDistribution.OnDemandBaseCapacity = newInt(-1)

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("onDemandBaseCapacity should be 0 or more"))
			})

			It("It fails when the spotInstancePools is not between 1 and 20", func() {
				ng.InstancesDistribution.SpotInstancePools = newInt(0)

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("spotInstancePools should be between 1 and 20"))

				ng.InstancesDistribution.SpotInstancePools = newInt(21)
				err = api.ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("spotInstancePools should be between 1 and 20"))
			})

			It("It fails when the onDemandPercentageAboveBaseCapacity is not between 0 and 100", func() {
				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(-1)

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("percentageAboveBase should be between 0 and 100"))

				ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity = newInt(101)
				err = api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("percentageAboveBase should be between 0 and 100"))
			})

			It("It fails when the spotAllocationStrategy is not a supported strategy", func() {
				ng.InstancesDistribution.SpotAllocationStrategy = strings.Pointer("unsupported-strategy")

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("spotAllocationStrategy should be one of: lowest-price, capacity-optimized, capacity-optimized-prioritized"))
			})

			It("It fails when the spotAllocationStrategy is capacity-optimized and spotInstancePools is specified", func() {
				ng.InstancesDistribution.SpotAllocationStrategy = strings.Pointer("capacity-optimized")
				ng.InstancesDistribution.SpotInstancePools = newInt(2)

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("spotInstancePools cannot be specified when also specifying spotAllocationStrategy: capacity-optimized"))
			})

			It("It fails when the spotAllocationStrategy is capacity-optimized-prioritized and spotInstancePools is specified", func() {
				ng.InstancesDistribution.SpotAllocationStrategy = strings.Pointer("capacity-optimized-prioritized")
				ng.InstancesDistribution.SpotInstancePools = newInt(2)

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("spotInstancePools cannot be specified when also specifying spotAllocationStrategy: capacity-optimized-prioritized"))
			})

			It("It does not fail when the spotAllocationStrategy is lowest-price and spotInstancePools is specified", func() {
				ng.InstancesDistribution.SpotAllocationStrategy = strings.Pointer("lowest-price")
				ng.InstancesDistribution.SpotInstancePools = newInt(2)

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("kubelet extra config", func() {
		Context("Instances distribution", func() {

			var ng *api.NodeGroup
			BeforeEach(func() {
				ng = newNodeGroup()
			})

			It("Forbids overriding basic fields", func() {
				testKeys := []string{"kind", "apiVersion", "address", "clusterDomain", "authentication",
					"authorization", "serverTLSBootstrap"}

				for _, key := range testKeys {
					ng.KubeletExtraConfig = &api.InlineDocument{
						key: "should-not-be-allowed",
					}
					err := api.ValidateNodeGroup(0, ng)
					Expect(err).To(MatchError(fmt.Sprintf("cannot override \"%s\" in kubelet config, as it's critical to eksctl functionality", key)))
				}
			})

			It("Allows other kubelet options", func() {
				ng.KubeletExtraConfig = &api.InlineDocument{
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
				err := api.ValidateNodeGroup(0, ng)
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

			var ng *api.NodeGroup
			BeforeEach(func() {
				ng = newNodeGroup()
			})

			It("Forbids setting volumeKmsKeyID without volumeEncrypted", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = nil
				ng.VolumeKmsKeyID = &kmsKeyID
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
			})

			It("Forbids setting volumeKmsKeyID with volumeEncrypted: false", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = &disabled
				ng.VolumeKmsKeyID = &kmsKeyID
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred())
			})

			It("Allows setting volumeKmsKeyID with volumeEncrypted: true", func() {
				ng.Name = nodegroup
				ng.VolumeSize = &volSize
				ng.VolumeEncrypted = &enabled
				ng.VolumeKmsKeyID = &kmsKeyID
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).ToNot(HaveOccurred())
			})

		})
	})

	Describe("FargateProfile", func() {
		Describe("Validate", func() {
			It("returns an error when the profile's name is empty", func() {
				profile := api.FargateProfile{
					Selectors: []api.FargateProfileSelector{
						{Namespace: "default"},
					},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: empty name"))
			})

			It("returns an error when the profile has a nil selectors array", func() {
				profile := api.FargateProfile{
					Name: "default",
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"default\": no profile selector"))
			})

			It("returns an error when the profile has an empty selectors array", func() {
				profile := api.FargateProfile{
					Name:      "default",
					Selectors: []api.FargateProfileSelector{},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"default\": no profile selector"))
			})

			It("returns an error when the profile's selectors do not have any namespace defined", func() {
				profile := api.FargateProfile{
					Name: "default",
					Selectors: []api.FargateProfileSelector{
						{},
					},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"default\": invalid profile selector at index #0: empty namespace"))
			})

			It("returns an error when the profile's name starts with eks-", func() {
				profile := api.FargateProfile{
					Name: "eks-foo",
					Selectors: []api.FargateProfileSelector{
						{Namespace: "default"},
					},
				}
				err := profile.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile \"eks-foo\": name should NOT start with \"eks-\""))
			})

			It("passes when a name and at least one selector with a namespace is defined", func() {
				profile := api.FargateProfile{
					Name: "default",
					Selectors: []api.FargateProfileSelector{
						{Namespace: "default"},
					},
				}
				err := profile.Validate()
				Expect(err).ToNot(HaveOccurred())
			})

			It("passes when a name and multiple selectors with a namespace is defined", func() {
				profile := api.FargateProfile{
					Name: "default",
					Selectors: []api.FargateProfileSelector{
						{
							Namespace: "default",
						},
						{
							Namespace: "dev",
						},
					},
				}
				err := profile.Validate()
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Bottlerocket node groups", func() {
		It("returns an error with unsupported fields", func() {
			cmd := "/usr/bin/some-command"
			doc := api.InlineDocument{
				"cgroupDriver": "systemd",
			}

			ngs := map[string]*api.NodeGroup{
				"PreBootstrapCommands": {
					NodeGroupBase: &api.NodeGroupBase{
						PreBootstrapCommands: []string{"/usr/bin/env true"},
					}},
				"OverrideBootstrapCommand": {
					NodeGroupBase: &api.NodeGroupBase{
						OverrideBootstrapCommand: &cmd,
					}},
				"KubeletExtraConfig": {KubeletExtraConfig: &doc},
				"overlapping Bottlerocket settings": {
					NodeGroupBase: &api.NodeGroupBase{
						Bottlerocket: &api.NodeGroupBottlerocket{
							Settings: &api.InlineDocument{
								"kubernetes": map[string]interface{}{
									"node-labels": map[string]string{
										"mylabel.example.com": "value",
									},
								},
							},
						},
					},
				},
			}

			for name, ng := range ngs {
				if ng.NodeGroupBase == nil {
					ng.NodeGroupBase = &api.NodeGroupBase{}
				}
				ng.AMIFamily = api.NodeImageFamilyBottlerocket
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred(), "foo", name)
			}
		})

		It("has no error with supported fields", func() {
			x := 32
			ngs := []*api.NodeGroup{
				{NodeGroupBase: &api.NodeGroupBase{Labels: map[string]string{"label": "label-value"}}},
				{NodeGroupBase: &api.NodeGroupBase{MaxPodsPerNode: x}},
				{
					NodeGroupBase: &api.NodeGroupBase{
						ScalingConfig: &api.ScalingConfig{
							MinSize: &x,
						},
					},
				},
			}

			for i, ng := range ngs {
				ng.AMIFamily = api.NodeImageFamilyBottlerocket
				Expect(api.ValidateNodeGroup(i, ng)).To(Succeed())
			}
		})
	})

	type kmsFieldCase struct {
		secretsEncryption *api.SecretsEncryption
		errSubstr         string
	}

	DescribeTable("KMS field validation", func(k kmsFieldCase) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Version = "1.15"

		clusterConfig.SecretsEncryption = k.secretsEncryption
		err := api.ValidateClusterConfig(clusterConfig)
		if k.errSubstr != "" {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(k.errSubstr))
		} else {
			Expect(err).ToNot(HaveOccurred())
		}
	},
		Entry("Nil secretsEncryption", kmsFieldCase{
			secretsEncryption: nil,
		}),
		Entry("Empty secretsEncryption.keyARN", kmsFieldCase{
			secretsEncryption: &api.SecretsEncryption{},
			errSubstr:         "secretsEncryption.keyARN is required",
		}),
	)

	Describe("Supported AMI Families", func() {
		var ng *api.NodeGroup
		BeforeEach(func() {
			ng = api.NewNodeGroup()
		})

		It("succeeds when AMI family is in supported list", func() {
			ng.AMIFamily = "AmazonLinux2"
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).NotTo(HaveOccurred())
		})

		It("fails when the AMIFamily is not supported", func() {
			ng.AMIFamily = "SomeTrash"
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).To(MatchError("AMI Family SomeTrash is not supported - use one of: AmazonLinux2, Ubuntu2004, Ubuntu1804, Bottlerocket, WindowsServer2019CoreContainer, WindowsServer2019FullContainer, WindowsServer2004CoreContainer"))
		})
	})

	Describe("Windows node groups", func() {
		It("returns an error with unsupported fields", func() {
			cmd := "start /wait msiexec.exe"
			doc := api.InlineDocument{
				"cgroupDriver": "systemd",
			}

			ngs := map[string]*api.NodeGroup{
				"OverrideBootstrapCommand": {NodeGroupBase: &api.NodeGroupBase{OverrideBootstrapCommand: &cmd}},
				"KubeletExtraConfig":       {KubeletExtraConfig: &doc, NodeGroupBase: &api.NodeGroupBase{}},
			}

			for name, ng := range ngs {
				ng.AMIFamily = api.NodeImageFamilyWindowsServer2019CoreContainer
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(HaveOccurred(), "Expected an error when provided %s", name)
			}
		})

		It("has no error with supported fields", func() {
			x := 32
			ngs := []*api.NodeGroup{
				{NodeGroupBase: &api.NodeGroupBase{Labels: map[string]string{"label": "label-value"}}},
				{NodeGroupBase: &api.NodeGroupBase{MaxPodsPerNode: x}},
				{NodeGroupBase: &api.NodeGroupBase{ScalingConfig: &api.ScalingConfig{MinSize: &x}}},
				{NodeGroupBase: &api.NodeGroupBase{PreBootstrapCommands: []string{"start /wait msiexec.exe"}}},
			}

			for i, ng := range ngs {
				ng.AMIFamily = api.NodeImageFamilyWindowsServer2019CoreContainer
				Expect(api.ValidateNodeGroup(i, ng)).To(Succeed())
			}
		})
	})

	type labelsTaintsEntry struct {
		labels map[string]string
		taints []api.NodeGroupTaint
		valid  bool
	}

	DescribeTable("Nodegroup label and taints validation", func(e labelsTaintsEntry) {
		ng := newNodeGroup()
		ng.Labels = e.labels
		ng.Taints = e.taints
		err := api.ValidateNodeGroup(0, ng)
		if e.valid {
			Expect(err).ToNot(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		Entry("disallowed label", labelsTaintsEntry{
			labels: map[string]string{
				"node-role.kubernetes.io/os": "linux",
			},
		}),

		Entry("disallowed label 2", labelsTaintsEntry{
			labels: map[string]string{
				"alpha.service-controller.kubernetes.io/test": "value",
			},
		}),

		Entry("empty labels and taints", labelsTaintsEntry{
			labels: map[string]string{},
			taints: []api.NodeGroupTaint{},
			valid:  true,
		}),

		Entry("allowed labels", labelsTaintsEntry{
			labels: map[string]string{
				"kubernetes.io/hostname":           "supercomputer",
				"beta.kubernetes.io/os":            "linux",
				"kubelet.kubernetes.io/palindrome": "telebuk",
			},
			valid: true,
		}),

		Entry("valid taints", labelsTaintsEntry{
			taints: []api.NodeGroupTaint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: "NoSchedule",
				},
				{
					Key:    "key2",
					Effect: "NoSchedule",
				},
				{
					Key:    "key3",
					Effect: "PreferNoSchedule",
				},
			},
			valid: true,
		}),

		Entry("missing taint effect", labelsTaintsEntry{
			taints: []api.NodeGroupTaint{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
		}),

		Entry("unsupported taint effect", labelsTaintsEntry{
			taints: []api.NodeGroupTaint{
				{
					Key:    "key2",
					Value:  "value1",
					Effect: "NoEffect",
				},
			},
		}),

		Entry("invalid value", labelsTaintsEntry{
			taints: []api.NodeGroupTaint{
				{
					Key:    "key3",
					Value:  "v@lue",
					Effect: "NoSchedule",
				},
			},
		}),
	)

	type updateConfigEntry struct {
		unavailable           *int64
		unavailablePercentage *int64
		valid                 bool
	}

	DescribeTable("UpdateConfig", func(e updateConfigEntry) {
		ng := newNodeGroup()
		ng.UpdateConfig = &api.NodeGroupUpdateConfig{
			MaxUnavailable:             e.unavailable,
			MaxUnavailableInPercentage: e.unavailablePercentage,
		}
		err := api.ValidateNodeGroup(0, ng)
		if e.valid {
			Expect(err).ToNot(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		Entry("max unavailable set", updateConfigEntry{
			unavailable: aws.Int64(1),
			valid:       true,
		}),
		Entry("max unavailable specified in percentage", updateConfigEntry{
			unavailablePercentage: aws.Int64(1),
			valid:                 true,
		}),
		Entry("both set", updateConfigEntry{
			unavailable:           aws.Int64(1),
			unavailablePercentage: aws.Int64(1),
			valid:                 false,
		}),
	)
})

func newInt(value int) *int {
	v := value
	return &v
}

func newNodeGroup() *api.NodeGroup {
	return &api.NodeGroup{
		NodeGroupBase: &api.NodeGroupBase{},
	}
}
