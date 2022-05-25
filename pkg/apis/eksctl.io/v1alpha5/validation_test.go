package v1alpha5_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
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
			Expect(err).NotTo(HaveOccurred())

			for i, ng := range cfg.NodeGroups {
				err = api.ValidateNodeGroup(i, ng)
				Expect(err).NotTo(HaveOccurred())
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

	Describe("nodeGroups[*].name validation", func() {
		var (
			cfg *api.ClusterConfig
			err error
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "ng_invalid-name-10"
			ng1 := cfg.NewNodeGroup()
			ng1.Name = "ng100_invalid_name"
			ng2 := cfg.NewNodeGroup()
			ng2.Name = "ng100@invalid-name"
		})

		It("should reject invalid nodegroup names", func() {
			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			for i, ng := range cfg.NodeGroups {
				err = api.ValidateNodeGroup(i, ng)
				Expect(err).To(HaveOccurred())
			}
		})
	})

	Describe("nodeGroups[*].containerRuntime validation", func() {

		It("allows accepted values", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.ContainerRuntime = aws.String(api.ContainerRuntimeDockerForWindows)
			err := api.ValidateNodeGroup(0, ng0)
			Expect(err).NotTo(HaveOccurred())

			ng0.ContainerRuntime = aws.String(api.ContainerRuntimeDockerD)
			err = api.ValidateNodeGroup(0, ng0)
			Expect(err).NotTo(HaveOccurred())

			ng0.ContainerRuntime = aws.String(api.ContainerRuntimeContainerD)
			ng0.AMIFamily = api.NodeImageFamilyAmazonLinux2
			err = api.ValidateNodeGroup(0, ng0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject invalid container runtime", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.ContainerRuntime = aws.String("invalid")
			err := api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			err = api.ValidateNodeGroup(0, ng0)
			Expect(err).To(HaveOccurred())
		})

		It("containerd is only allowed for AL2 or Windows", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.ContainerRuntime = aws.String(api.ContainerRuntimeContainerD)
			ng0.AMIFamily = api.NodeImageFamilyBottlerocket
			err := api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			err = api.ValidateNodeGroup(0, ng0)
			Expect(err).To(HaveOccurred())
			ng0.AMIFamily = api.NodeImageFamilyWindowsServer2019CoreContainer
			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			err = api.ValidateNodeGroup(0, ng0)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("nodeGroups[*].ami validation", func() {
		It("should require overrideBootstrapCommand if ami is set", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.AMI = "ami-1234"
			Expect(api.ValidateNodeGroup(0, ng0)).To(MatchError(ContainSubstring("overrideBootstrapCommand is required when using a custom AMI ")))
		})
		It("should not require overrideBootstrapCommand if ami is set and type is Bottlerocket", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.AMI = "ami-1234"
			ng0.AMIFamily = api.NodeImageFamilyBottlerocket
			Expect(api.ValidateNodeGroup(0, ng0)).To(Succeed())
		})
		It("should accept ami with a overrideBootstrapCommand set", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.AMI = "ami-1234"
			ng0.OverrideBootstrapCommand = aws.String("echo 'yo'")
			Expect(api.ValidateNodeGroup(0, ng0)).To(Succeed())
		})
	})

	Describe("nodeGroups[*].maxInstanceLifetime validation", func() {
		It("should reject if value is below a day", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.MaxInstanceLifetime = aws.Int(5)
			err := api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			err = api.ValidateNodeGroup(0, ng0)
			Expect(err).To(MatchError(ContainSubstring("maximum instance lifetime must have a minimum value of 86,400 seconds (one day), but was: 5")))
		})
		It("setting it if greater than or equal to one day", func() {
			cfg := api.NewClusterConfig()
			ng0 := cfg.NewNodeGroup()
			ng0.Name = "node-group"
			ng0.MaxInstanceLifetime = aws.Int(api.OneDay)
			err := api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			err = api.ValidateNodeGroup(0, ng0)
			Expect(err).NotTo(HaveOccurred())
			ng0.MaxInstanceLifetime = aws.Int(api.OneDay + 1000)
			Expect(err).NotTo(HaveOccurred())
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

			ng0.IAM.AttachPolicy = cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)
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
			Expect(err).NotTo(HaveOccurred())

			for i, ng := range cfg.NodeGroups {
				err = api.ValidateNodeGroup(i, ng)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should allow setting only instanceProfileARN", func() {
			ng1.IAM.InstanceProfileARN = "p1"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should allow setting only instanceRoleARN", func() {
			ng1.IAM.InstanceRoleARN = "r1"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not allow setting AWSLoadBalancerController and albIngress", func() {
			ng1.IAM.WithAddonPolicies.AWSLoadBalancerController = aws.Bool(true)
			ng1.IAM.WithAddonPolicies.DeprecatedALBIngress = aws.Bool(true)

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(MatchError(`"awsLoadBalancerController" and "albIngress" cannot both be configured, please use "awsLoadBalancerController" as "albIngress" is deprecated`))
		})

		It("should allow setting instanceProfileARN and instanceRoleARN", func() {
			ng1.IAM.InstanceProfileARN = "p1"
			ng1.IAM.InstanceRoleARN = "r1"

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).NotTo(HaveOccurred())
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

		It("should not allow setting instanceRoleARN and attachPolicy", func() {
			ng1.IAM.InstanceRoleARN = "r1"
			ng1.IAM.AttachPolicy = cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			err = api.ValidateNodeGroup(1, ng1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("nodeGroups[1].iam.instanceRoleARN and nodeGroups[1].iam.attachPolicy cannot be set at the same time"))
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
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass when iam.withOIDC is disabled", func() {
			cfg.IAM.WithOIDC = api.Disabled()

			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass when iam.withOIDC is enabled", func() {
			cfg.IAM.WithOIDC = api.Enabled()

			err = api.ValidateClusterConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
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
			Expect(api.ValidateClusterConfig(cfg)).To(Succeed())
		})

		It("should handle unknown types", func() {
			cfg.CloudWatch.ClusterLogging.EnableTypes = []string{"anything"}

			err = api.ValidateClusterConfig(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(`log type "anything" (cloudWatch.clusterLogging.enableTypes[0]) is unknown`)))
		})
	})

	type logRetentionEntry struct {
		logging *api.ClusterCloudWatchLogging

		expectedErr string
	}

	DescribeTable("CloudWatch log retention", func(l logRetentionEntry) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.CloudWatch.ClusterLogging = l.logging
		err := api.ValidateClusterConfig(clusterConfig)
		if l.expectedErr != "" {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(l.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())
	},
		Entry("invalid value", logRetentionEntry{
			logging: &api.ClusterCloudWatchLogging{
				LogRetentionInDays: 42,
				EnableTypes:        []string{"api"},
			},
			expectedErr: `invalid value 42 for logRetentionInDays; supported values are [1 3 5 7`,
		}),

		Entry("valid value", logRetentionEntry{
			logging: &api.ClusterCloudWatchLogging{
				LogRetentionInDays: 545,
				EnableTypes:        []string{"api"},
			},
		}),

		Entry("log retention without enableTypes", logRetentionEntry{
			logging: &api.ClusterCloudWatchLogging{
				LogRetentionInDays: 545,
			},
			expectedErr: "cannot set cloudWatch.clusterLogging.logRetentionInDays without enabling log types",
		}),
	)

	Describe("Cluster Endpoint access", func() {
		var cfg *api.ClusterConfig

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
		})

		When("VPC is not set", func() {
			It("should have cluster endpoint access", func() {
				err := api.ValidateClusterConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("VPC is set", func() {
			BeforeEach(func() {
				cfg.VPC = &api.ClusterVPC{}
			})

			When("no cluster endpoint config is set", func() {
				It("should not error", func() {
					err := api.ValidateClusterConfig(cfg)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("cluster endpoint config exists", func() {
				It("should not error on private=true, public=true", func() {
					cfg.VPC.ClusterEndpoints =
						&api.ClusterEndpoints{PrivateAccess: api.Enabled(), PublicAccess: api.Enabled()}
					err := api.ValidateClusterConfig(cfg)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not error on private=false, public=true", func() {
					cfg.VPC.ClusterEndpoints =
						&api.ClusterEndpoints{PrivateAccess: api.Disabled(), PublicAccess: api.Enabled()}
					err := api.ValidateClusterConfig(cfg)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not error on private=true, public=false", func() {
					cfg.VPC.ClusterEndpoints =
						&api.ClusterEndpoints{PrivateAccess: api.Enabled(), PublicAccess: api.Disabled()}
					err := api.ValidateClusterConfig(cfg)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should error on private=false, public=false", func() {
					cfg.VPC.ClusterEndpoints = &api.ClusterEndpoints{PrivateAccess: api.Disabled(), PublicAccess: api.Disabled()}
					err := api.ValidateClusterConfig(cfg)
					Expect(err).To(MatchError(api.ErrClusterEndpointNoAccess))
				})
			})
		})
	})

	type localZonesEntry struct {
		updateClusterConfig func(*api.ClusterConfig)

		expectedErr string
	}

	DescribeTable("AWS Local Zones", func(e localZonesEntry) {
		clusterConfig := api.NewClusterConfig()
		e.updateClusterConfig(clusterConfig)
		clusterConfig.LocalZones = []string{"us-west-2-lax-1a", "us-west-2-lax-1b"}

		err := api.ValidateClusterConfig(clusterConfig)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
	},
		Entry("custom VPC", localZonesEntry{
			updateClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC = &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc-custom",
					},
				}

			},
			expectedErr: "localZones are not supported with a pre-existing VPC",
		}),

		Entry("HighlyAvailable NAT gateway", localZonesEntry{
			updateClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC = &api.ClusterVPC{
					NAT: &api.ClusterNAT{
						Gateway: aws.String("HighlyAvailable"),
					},
				}
			},
			expectedErr: "HighlyAvailable NAT gateway is not supported for localZones",
		}),

		Entry("private cluster", localZonesEntry{
			updateClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.PrivateCluster = &api.PrivateCluster{
					Enabled: true,
				}
			},

			expectedErr: "localZones cannot be used in a fully-private cluster",
		}),

		Entry("IPv6", localZonesEntry{
			updateClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
					IPFamily: api.IPV6Family,
				}
				clusterConfig.Addons = []*api.Addon{
					{
						Name: "vpc-cni",
					},
					{
						Name: "coredns",
					},
					{
						Name: "kube-proxy",
					},
				}
				clusterConfig.IAM.WithOIDC = api.Enabled()
				clusterConfig.VPC.NAT = nil
			},

			expectedErr: "localZones are not supported with IPv6",
		}),
	)

	Describe("ValidatePrivateCluster", func() {
		var (
			cfg *api.ClusterConfig
			vpc *api.ClusterVPC
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			vpc = api.NewClusterVPC(false)
			cfg.VPC = vpc
			cfg.PrivateCluster = &api.PrivateCluster{
				Enabled: true,
			}
		})
		When("private cluster is enabled", func() {
			It("validates the config", func() {
				err := api.ValidateClusterConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("vpc is provided but no private subnets", func() {
			It("fails the validation", func() {
				cfg.VPC.Subnets = &api.ClusterSubnets{}
				cfg.VPC.ID = "id"
				err := api.ValidateClusterConfig(cfg)
				Expect(err).To(MatchError(ContainSubstring("vpc.subnets.private must be specified in a fully-private cluster when a pre-existing VPC is supplied")))
			})
		})
		When("additional endpoints are defined with skip endpoints", func() {
			It("fails the validation", func() {
				cfg.PrivateCluster.AdditionalEndpointServices = []string{api.EndpointServiceCloudFormation}
				cfg.PrivateCluster.SkipEndpointCreation = true
				err := api.ValidateClusterConfig(cfg)
				Expect(err).To(MatchError(ContainSubstring("privateCluster.additionalEndpointServices cannot be set when privateCluster.skipEndpointCreation is true")))
			})
		})
		When("additional endpoints are defined", func() {
			It("validates the endpoint configuration", func() {
				cfg.PrivateCluster.AdditionalEndpointServices = []string{api.EndpointServiceCloudFormation}
				err := api.ValidateClusterConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("additional endpoints are defined incorrectly", func() {
			It("fails the endpoint validation", func() {
				cfg.PrivateCluster.AdditionalEndpointServices = []string{"unknown"}
				err := api.ValidateClusterConfig(cfg)
				Expect(err).To(MatchError(ContainSubstring("invalid value in privateCluster.additionalEndpointServices")))
			})
		})
		When("private cluster is enabled with skip endpoints", func() {
			It("does not fail the validation", func() {
				cfg.PrivateCluster.SkipEndpointCreation = true
				err := api.ValidateClusterConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
	Describe("network config", func() {
		var (
			cfg *api.ClusterConfig
			vpc *api.ClusterVPC
			err error
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			vpc = api.NewClusterVPC(false)
			cfg.VPC = vpc
		})

		Context("ipFamily", func() {
			It("should not error default ipFamily setting", func() {
				err = api.ValidateClusterConfig(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.KubernetesNetworkConfig.IPFamily).To(Equal(api.IPV4Family))
			})

			When("ipFamily isn't IPv4 or IPv6", func() {
				It("returns an error", func() {
					cfg.KubernetesNetworkConfig.IPFamily = "invalid"
					err = api.ValidateClusterConfig(cfg)
					Expect(err).To(MatchError(ContainSubstring(`invalid value "invalid" for ipFamily; allowed are IPv4 and IPv6`)))
				})
			})

			When("ipFamily is empty", func() {
				It("treats it as IPv4 and does not return an error", func() {
					cfg.KubernetesNetworkConfig.IPFamily = ""
					Expect(api.ValidateClusterConfig(cfg)).To(Succeed())
				})
			})

			When("ipFamily is set to IPv6", func() {
				JustBeforeEach(func() {
					cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
				})

				It("accepts that setting", func() {
					cfg.VPC.NAT = nil
					cfg.VPC.IPv6Cidr = "foo"
					cfg.VPC.IPv6Pool = "bar"
					cfg.Addons = append(cfg.Addons,
						&api.Addon{Name: api.KubeProxyAddon},
						&api.Addon{Name: api.CoreDNSAddon},
						&api.Addon{Name: api.VPCCNIAddon},
					)
					cfg.IAM = &api.ClusterIAM{
						WithOIDC: api.Enabled(),
					}
					cfg.Metadata.Version = api.Version1_21
					err = cfg.ValidateVPCConfig()
					Expect(err).ToNot(HaveOccurred())
					cfg.Metadata.Version = "1.31"
					err = cfg.ValidateVPCConfig()
					Expect(err).ToNot(HaveOccurred())
				})

				When("the casing of ipv6 isn't standard", func() {
					It("accepts that setting", func() {
						cfg.KubernetesNetworkConfig.IPFamily = "iPv6"
						cfg.VPC.NAT = nil
						cfg.VPC.IPv6Cidr = "foo"
						cfg.VPC.IPv6Pool = "bar"
						cfg.Addons = append(cfg.Addons,
							&api.Addon{Name: api.KubeProxyAddon},
							&api.Addon{Name: api.CoreDNSAddon},
							&api.Addon{Name: api.VPCCNIAddon},
						)
						cfg.IAM = &api.ClusterIAM{
							WithOIDC: api.Enabled(),
						}
						cfg.Metadata.Version = api.Version1_21
						err = cfg.ValidateVPCConfig()
						Expect(err).ToNot(HaveOccurred())
						cfg.Metadata.Version = "1.31"
						err = cfg.ValidateVPCConfig()
						Expect(err).ToNot(HaveOccurred())
					})
				})

				When("ipFamily is set to IPv6 but version is not or too low", func() {
					It("returns an error", func() {
						cfg.VPC.NAT = nil
						cfg.Addons = append(cfg.Addons,
							&api.Addon{Name: api.KubeProxyAddon},
							&api.Addon{Name: api.CoreDNSAddon},
							&api.Addon{Name: api.VPCCNIAddon},
						)
						cfg.IAM = &api.ClusterIAM{
							WithOIDC: api.Enabled(),
						}
						cfg.Metadata.Version = ""
						err = api.ValidateClusterConfig(cfg)
						Expect(err).To(MatchError(ContainSubstring("failed to convert  cluster version to semver: unable to parse first version")))
						cfg.Metadata.Version = api.Version1_12
						err = api.ValidateClusterConfig(cfg)
						Expect(err).To(MatchError(ContainSubstring("cluster version must be >= 1.21")))
					})
				})

				When("ipFamily is set to IPv6 but no managed addons are provided", func() {
					It("it returns an error including which addons are missing", func() {
						cfg.VPC.NAT = nil
						cfg.IAM = &api.ClusterIAM{
							WithOIDC: api.Enabled(),
						}
						cfg.Addons = append(cfg.Addons, &api.Addon{Name: api.KubeProxyAddon})
						err = api.ValidateClusterConfig(cfg)
						Expect(err).To(MatchError(ContainSubstring("the default core addons must be defined for IPv6; missing addon(s): vpc-cni, coredns")))
					})
				})

				When("the vpc-cni version is configured", func() {
					When("the version of the vpc-cni is too low", func() {
						It("returns an error", func() {
							cfg.Metadata.Version = api.Version1_22
							cfg.IAM = &api.ClusterIAM{
								WithOIDC: api.Enabled(),
							}
							cfg.Addons = append(cfg.Addons,
								&api.Addon{Name: api.KubeProxyAddon},
								&api.Addon{Name: api.CoreDNSAddon},
								&api.Addon{Name: api.VPCCNIAddon, Version: "1.9.0"},
							)
							cfg.VPC.NAT = nil
							err = api.ValidateClusterConfig(cfg)
							Expect(err).To(MatchError(ContainSubstring("vpc-cni version must be at least version 1.10.0 for IPv6")))
						})
					})

					When("the version of the vpc-cni is supported", func() {
						It("does not error", func() {
							cfg.Metadata.Version = api.Version1_22
							cfg.IAM = &api.ClusterIAM{
								WithOIDC: api.Enabled(),
							}
							cfg.Addons = append(cfg.Addons,
								&api.Addon{Name: api.KubeProxyAddon},
								&api.Addon{Name: api.CoreDNSAddon},
								&api.Addon{Name: api.VPCCNIAddon, Version: "1.10"},
							)
							cfg.VPC.NAT = nil
							err = cfg.ValidateVPCConfig()
							Expect(err).NotTo(HaveOccurred())
						})
					})

					When("the version of the vpc-cni is not configured", func() {
						It("does not error", func() {
							cfg.Metadata.Version = api.Version1_22
							cfg.IAM = &api.ClusterIAM{
								WithOIDC: api.Enabled(),
							}
							cfg.Addons = append(cfg.Addons,
								&api.Addon{Name: api.KubeProxyAddon},
								&api.Addon{Name: api.CoreDNSAddon},
								&api.Addon{Name: api.VPCCNIAddon},
							)
							cfg.VPC.NAT = nil
							err = api.ValidateClusterConfig(cfg)
							Expect(err).NotTo(HaveOccurred())
						})
					})

					When("the version of the vpc-cni is latest", func() {
						It("does not error", func() {
							cfg.Metadata.Version = api.Version1_22
							cfg.IAM = &api.ClusterIAM{
								WithOIDC: api.Enabled(),
							}
							cfg.Addons = append(cfg.Addons,
								&api.Addon{Name: api.KubeProxyAddon},
								&api.Addon{Name: api.CoreDNSAddon},
								&api.Addon{Name: api.VPCCNIAddon, Version: "latest"},
							)
							cfg.VPC.NAT = nil
							err = cfg.ValidateVPCConfig()
							Expect(err).NotTo(HaveOccurred())
						})
					})

					When("the version of the vpc-cni is invalid", func() {
						It("it returns an error", func() {
							cfg.Metadata.Version = api.Version1_22
							cfg.IAM = &api.ClusterIAM{
								WithOIDC: api.Enabled(),
							}
							cfg.Addons = append(cfg.Addons,
								&api.Addon{Name: api.KubeProxyAddon},
								&api.Addon{Name: api.CoreDNSAddon},
								&api.Addon{Name: api.VPCCNIAddon, Version: "1.invalid!semver"},
							)
							cfg.VPC.NAT = nil
							err = api.ValidateClusterConfig(cfg)
							Expect(err).To(MatchError(ContainSubstring("failed to parse version")))
						})
					})
				})

				When("iam is not set", func() {
					It("returns an error", func() {
						cfg.Addons = append(cfg.Addons,
							&api.Addon{Name: api.KubeProxyAddon},
							&api.Addon{Name: api.CoreDNSAddon},
							&api.Addon{Name: api.VPCCNIAddon},
						)
						err = api.ValidateClusterConfig(cfg)
						Expect(err).To(MatchError(ContainSubstring("oidc needs to be enabled if IPv6 is set")))
					})
				})

				When("iam is set but OIDC is disabled", func() {
					It("returns an error", func() {
						cfg.IAM = &api.ClusterIAM{
							WithOIDC: api.Disabled(),
						}
						cfg.Addons = append(cfg.Addons,
							&api.Addon{Name: api.KubeProxyAddon},
							&api.Addon{Name: api.CoreDNSAddon},
							&api.Addon{Name: api.VPCCNIAddon},
						)
						err = api.ValidateClusterConfig(cfg)
						Expect(err).To(MatchError(ContainSubstring("oidc needs to be enabled if IPv6 is set")))
					})
				})

				When("ipFamily is set to IPv6 and vpc.NAT is defined", func() {
					It("it returns an error", func() {
						cfg.Metadata.Version = api.Version1_22
						cfg.IAM = &api.ClusterIAM{
							WithOIDC: api.Enabled(),
						}
						cfg.Addons = append(cfg.Addons,
							&api.Addon{Name: api.KubeProxyAddon},
							&api.Addon{Name: api.CoreDNSAddon},
							&api.Addon{Name: api.VPCCNIAddon},
						)
						cfg.VPC.NAT = &api.ClusterNAT{}
						err = cfg.ValidateVPCConfig()
						Expect(err).To(MatchError(ContainSubstring("setting NAT is not supported with IPv6")))
					})
				})

				When("ipFamily is set to IPv6 and serviceIPv4CIDR is not empty", func() {
					It("it returns an error", func() {
						cfg.Metadata.Version = api.Version1_22
						cfg.IAM = &api.ClusterIAM{
							WithOIDC: api.Enabled(),
						}
						cfg.Addons = append(cfg.Addons,
							&api.Addon{Name: api.KubeProxyAddon},
							&api.Addon{Name: api.CoreDNSAddon},
							&api.Addon{Name: api.VPCCNIAddon},
						)
						cfg.KubernetesNetworkConfig.ServiceIPv4CIDR = "192.168.0.0/24"
						cfg.VPC.NAT = nil
						err = api.ValidateClusterConfig(cfg)
						Expect(err).To(MatchError(ContainSubstring("service ipv4 cidr is not supported with IPv6")))
					})
				})

				When("ipFamily is set to IPv6 and AutoAllocateIPv6 is set", func() {
					It("it returns an error", func() {
						cfg.VPC.AutoAllocateIPv6 = api.Enabled()
						cfg.Metadata.Version = api.Version1_22
						cfg.IAM = &api.ClusterIAM{
							WithOIDC: api.Enabled(),
						}
						cfg.Addons = append(cfg.Addons,
							&api.Addon{Name: api.KubeProxyAddon},
							&api.Addon{Name: api.CoreDNSAddon},
							&api.Addon{Name: api.VPCCNIAddon},
						)
						cfg.VPC.NAT = nil
						err = cfg.ValidateVPCConfig()
						Expect(err).To(MatchError(ContainSubstring("auto allocate ipv6 is not supported with IPv6")))
					})
				})
			})
		})

		Context("extraCIDRs", func() {
			It("validates cidrs", func() {
				cfg.VPC.ExtraCIDRs = []string{"192.168.0.0/24"}
				cfg.VPC.PublicAccessCIDRs = []string{"3.48.58.68/24"}
				err = cfg.ValidateVPCConfig()
				Expect(err).ToNot(HaveOccurred())
			})

			When("extraCIDRs has an invalid cidr", func() {
				It("returns an error", func() {
					cfg.VPC.ExtraCIDRs = []string{"not-a-cidr"}
					err = cfg.ValidateVPCConfig()
					Expect(err).To(HaveOccurred())
				})
			})

			When("public access cidrs has an invalid cidr", func() {
				It("returns an error", func() {
					cfg.VPC.PublicAccessCIDRs = []string{"48.58.68/24"}
					err = cfg.ValidateVPCConfig()
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("ipv6 CIDRs", func() {
			When("IPv6Cidr or IPv6CidrPool is provided and ipv6 is not set", func() {
				It("returns an error", func() {
					cfg.KubernetesNetworkConfig.IPFamily = api.IPV4Family
					cfg.VPC.IPv6Cidr = "foo"
					err = cfg.ValidateVPCConfig()
					Expect(err).To(MatchError("Ipv6Cidr and Ipv6CidrPool are only supported when IPFamily is set to IPv6"))

					cfg.KubernetesNetworkConfig.IPFamily = api.IPV4Family
					cfg.VPC.IPv6Cidr = ""
					cfg.VPC.IPv6Pool = "bar"
					err = cfg.ValidateVPCConfig()
					Expect(err).To(MatchError("Ipv6Cidr and Ipv6CidrPool are only supported when IPFamily is set to IPv6"))
				})
			})

			When("only one of IPv6Cidr or IPv6CidrPool is provided and ipv6 is set", func() {
				It("returns an error", func() {
					cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
					cfg.VPC.IPv6Cidr = "foo"
					err = cfg.ValidateVPCConfig()
					Expect(err).To(MatchError("Ipv6Cidr and Ipv6Pool must both be configured to use a custom IPv6 CIDR and address pool"))

					cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
					cfg.VPC.IPv6Cidr = ""
					cfg.VPC.IPv6Pool = "bar"
					err = cfg.ValidateVPCConfig()
					Expect(err).To(MatchError("Ipv6Cidr and Ipv6Pool must both be configured to use a custom IPv6 CIDR and address pool"))
				})
			})

			When("it's set alongside VPC.ID", func() {
				It("returns an error", func() {
					cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
					cfg.VPC.IPv6Cidr = "foo"
					cfg.VPC.IPv6Pool = "bar"
					cfg.VPC.ID = "123"
					err = cfg.ValidateVPCConfig()
					Expect(err).To(MatchError("cannot provide VPC.IPv6Cidr when using a pre-existing VPC.ID"))
				})
			})
		})

		Context("extraIPv6CIDRs", func() {
			It("validates cidrs", func() {
				cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
				cfg.Metadata.Version = api.LatestVersion
				cfg.Addons = append(cfg.Addons,
					&api.Addon{Name: api.KubeProxyAddon},
					&api.Addon{Name: api.CoreDNSAddon},
					&api.Addon{Name: api.VPCCNIAddon},
				)
				cfg.IAM = &api.ClusterIAM{
					WithOIDC: api.Enabled(),
				}
				cfg.VPC.ExtraIPv6CIDRs = []string{"2002::1234:abcd:ffff:c0a8:101/64"}
				cfg.VPC.NAT = nil
				err = cfg.ValidateVPCConfig()
				Expect(err).ToNot(HaveOccurred())
			})

			When("extraIPv6CIDRs has an invalid cidr", func() {
				It("returns an error", func() {
					cfg.VPC.ExtraIPv6CIDRs = []string{"not-a-cidr"}
					cfg.Metadata.Version = api.LatestVersion
					cfg.Addons = append(cfg.Addons,
						&api.Addon{Name: api.KubeProxyAddon},
						&api.Addon{Name: api.CoreDNSAddon},
						&api.Addon{Name: api.VPCCNIAddon},
					)
					cfg.IAM = &api.ClusterIAM{
						WithOIDC: api.Enabled(),
					}
					cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
					err = cfg.ValidateVPCConfig()
					Expect(err).To(HaveOccurred())

					cfg.VPC.ExtraIPv6CIDRs = []string{"2002::1234:abcd:ffff:c0a8:101/644"}
					err = cfg.ValidateVPCConfig()
					Expect(err).To(HaveOccurred())
				})
			})

			When("when ipv4 is configured", func() {
				It("returns an error", func() {
					cfg.KubernetesNetworkConfig.IPFamily = api.IPV4Family
					cfg.Metadata.Version = api.LatestVersion
					cfg.Addons = append(cfg.Addons,
						&api.Addon{Name: api.KubeProxyAddon},
						&api.Addon{Name: api.CoreDNSAddon},
						&api.Addon{Name: api.VPCCNIAddon},
					)
					cfg.IAM = &api.ClusterIAM{
						WithOIDC: api.Enabled(),
					}
					cfg.VPC.ExtraIPv6CIDRs = []string{"2002::1234:abcd:ffff:c0a8:101/644"}
					err = cfg.ValidateVPCConfig()
					Expect(err).To(MatchError("cannot specify vpc.extraIPv6CIDRs with an IPv4 cluster"))
				})
			})

		})
	})

	Describe("ValidatePrivateCluster", func() {
		var (
			cfg *api.ClusterConfig
			vpc *api.ClusterVPC
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			vpc = api.NewClusterVPC(false)
			cfg.VPC = vpc
			cfg.PrivateCluster = &api.PrivateCluster{
				Enabled: true,
			}
		})
		When("private cluster is enabled", func() {
			It("validates the config", func() {
				err := cfg.ValidatePrivateCluster()
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("vpc is provided but no private subnets", func() {
			It("fails the validation", func() {
				cfg.VPC.Subnets = &api.ClusterSubnets{}
				cfg.VPC.ID = "id"
				err := cfg.ValidatePrivateCluster()
				Expect(err).To(MatchError(ContainSubstring("vpc.subnets.private must be specified in a fully-private cluster when a pre-existing VPC is supplied")))
			})
		})
		When("additional endpoints are defined with skip endpoints", func() {
			It("fails the validation", func() {
				cfg.PrivateCluster.AdditionalEndpointServices = []string{api.EndpointServiceCloudFormation}
				cfg.PrivateCluster.SkipEndpointCreation = true
				err := cfg.ValidatePrivateCluster()
				Expect(err).To(MatchError(ContainSubstring("privateCluster.additionalEndpointServices cannot be set when privateCluster.skipEndpointCreation is true")))
			})
		})
		When("additional endpoints are defined", func() {
			It("validates the endpoint configuration", func() {
				cfg.PrivateCluster.AdditionalEndpointServices = []string{api.EndpointServiceCloudFormation}
				err := cfg.ValidatePrivateCluster()
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("additional endpoints are defined incorrectly", func() {
			It("fails the endpoint validation", func() {
				cfg.PrivateCluster.AdditionalEndpointServices = []string{"unknown"}
				err := cfg.ValidatePrivateCluster()
				Expect(err).To(MatchError(ContainSubstring("invalid value in privateCluster.additionalEndpointServices")))
			})
		})
		When("private cluster is enabled with skip endpoints", func() {
			It("does not fail the validation", func() {
				cfg.PrivateCluster.SkipEndpointCreation = true
				err := cfg.ValidatePrivateCluster()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	type cpuCreditsEntry struct {
		modifyNodeGroup func(*api.NodeGroup)
		expectedError   string
	}

	DescribeTable("cpuCredits", func(e cpuCreditsEntry) {
		ng := &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{},
			InstancesDistribution: &api.NodeGroupInstancesDistribution{
				InstanceTypes: []string{"t3.medium", "t3.large"},
			},
			CPUCredits: aws.String("unlimited"),
		}
		if e.modifyNodeGroup != nil {
			e.modifyNodeGroup(ng)
		}

		err := api.ValidateNodeGroup(0, ng)
		if e.expectedError != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedError)))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
	},

		Entry("instanceType is set", cpuCreditsEntry{
			modifyNodeGroup: func(ng *api.NodeGroup) {
				ng.InstanceType = "mixed"
			},
		}),
		Entry("instanceType is not set", cpuCreditsEntry{}),
		Entry("instancesDistribution is not set", cpuCreditsEntry{
			modifyNodeGroup: func(ng *api.NodeGroup) {
				ng.InstancesDistribution = nil
			},
			expectedError: "cpuCredits option set for nodegroup, but it has no t2/t3 instance types",
		}),
		Entry("instancesDistribution.instanceTypes is not set", cpuCreditsEntry{
			modifyNodeGroup: func(ng *api.NodeGroup) {
				ng.InstancesDistribution.InstanceTypes = nil
			},
			expectedError: "at least two instance types have to be specified for mixed nodegroups",
		}),
	)

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
				Expect(err).NotTo(HaveOccurred())
			})

			It("It doesn't fail when instance distribution is enabled and instanceType is \"mixed\"", func() {
				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{"t3.medium"}

				err := api.ValidateNodeGroup(0, ng)
				Expect(err).NotTo(HaveOccurred())
			})

			It("fails in case of arm-gpu distribution instance type", func() {
				ng.InstanceType = "mixed"
				ng.InstancesDistribution.InstanceTypes = []string{"g5g.medium"}
				ng.AMIFamily = api.NodeImageFamilyAmazonLinux2
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("ARM GPU instance types are not supported for unmanaged nodegroups with AMIFamily AmazonLinux2"))
			})

			It("fails in case of arm-gpu instance type", func() {
				ng.InstanceType = "g5g.medium"
				ng.InstancesDistribution.InstanceTypes = nil
				ng.AMIFamily = api.NodeImageFamilyAmazonLinux2
				err := api.ValidateNodeGroup(0, ng)
				Expect(err).To(MatchError("ARM GPU instance types are not supported for unmanaged nodegroups with AMIFamily AmazonLinux2"))
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
				Expect(err).NotTo(HaveOccurred())
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
				Expect(err).NotTo(HaveOccurred())
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
				Expect(err).NotTo(HaveOccurred())
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
				Expect(err).NotTo(HaveOccurred())
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
				Expect(err).NotTo(HaveOccurred())
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
			Expect(err).NotTo(HaveOccurred())
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

		It("mutates the AMIFamily to the correct value when the capitalisation is incorrect", func() {
			ng.AMIFamily = "aMAZONlINUx2"
			Expect(api.ValidateNodeGroup(0, ng)).To(Succeed())
			Expect(ng.AMIFamily).To(Equal("AmazonLinux2"))

			mng := api.NewManagedNodeGroup()
			mng.AMIFamily = "bOTTLEROCKEt"
			Expect(api.ValidateManagedNodeGroup(0, mng)).To(Succeed())
			Expect(mng.AMIFamily).To(Equal("Bottlerocket"))
		})

		It("fails when the AMIFamily is not supported", func() {
			ng.AMIFamily = "SomeTrash"
			err := api.ValidateNodeGroup(0, ng)
			Expect(err).To(MatchError("AMI Family SomeTrash is not supported - use one of: AmazonLinux2, Ubuntu2004, Ubuntu1804, Bottlerocket, WindowsServer2019CoreContainer, WindowsServer2019FullContainer, WindowsServer2004CoreContainer, WindowsServer20H2CoreContainer"))
		})
	})

	Describe("Windows node groups", func() {
		It("returns an error with unsupported fields", func() {
			doc := api.InlineDocument{
				"cgroupDriver": "systemd",
			}

			ngs := map[string]*api.NodeGroup{
				"KubeletExtraConfig": {KubeletExtraConfig: &doc, NodeGroupBase: &api.NodeGroupBase{}},
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

	Describe("Karpenter", func() {
		It("returns an error when OIDC is not set", func() {
			cfg := api.NewClusterConfig()
			cfg.Karpenter = &api.Karpenter{
				Version: "0.6.1",
			}
			Expect(api.ValidateClusterConfig(cfg)).To(MatchError(ContainSubstring("failed to validate karpenter config: iam.withOIDC must be enabled with Karpenter")))
		})

		It("returns an error when version is missing", func() {
			cfg := api.NewClusterConfig()
			cfg.Karpenter = &api.Karpenter{}
			Expect(api.ValidateClusterConfig(cfg)).To(MatchError(ContainSubstring("version field is required if installing Karpenter is enabled")))
		})

		It("returns an error when version is missing", func() {
			cfg := api.NewClusterConfig()
			cfg.IAM.WithOIDC = aws.Bool(true)
			cfg.Karpenter = &api.Karpenter{
				Version: "isitmeeeyourlookingfoorrrr",
			}
			Expect(api.ValidateClusterConfig(cfg)).To(MatchError(ContainSubstring("failed to parse karpenter version")))
		})

		It("returns an error when the version is not supported", func() {
			cfg := api.NewClusterConfig()
			cfg.IAM.WithOIDC = aws.Bool(true)
			cfg.Karpenter = &api.Karpenter{
				Version: "0.10.0",
			}
			Expect(api.ValidateClusterConfig(cfg)).To(MatchError(ContainSubstring("failed to validate karpenter config: maximum supported version is 0.9")))
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
			Expect(err).NotTo(HaveOccurred())
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

	Describe("Availability Zones", func() {
		When("the config file does not specify any AZ", func() {
			It("skips validation", func() {
				Expect(api.ValidateClusterConfig(api.NewClusterConfig())).NotTo(HaveOccurred())
			})
		})

		When("the config file contains too few availability zones", func() {
			It("returns an error", func() {
				cfg := api.NewClusterConfig()
				cfg.AvailabilityZones = append(cfg.AvailabilityZones, "az-1")
				Expect(api.ValidateClusterConfig(cfg)).To(MatchError("only 1 zone(s) specified [az-1], 2 are required (can be non-unique)"))
			})
		})
	})

	Describe("Validate SecretsEncryption", func() {
		var cfg *api.ClusterConfig

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
		})

		When("a key ARN is set", func() {
			When("the key is valid", func() {
				It("does not return an error", func() {
					cfg.SecretsEncryption = &api.SecretsEncryption{
						KeyARN: "arn:aws:kms:us-west-2:000000000000:key/12345-12345",
					}
					err := api.ValidateClusterConfig(cfg)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("the key is invalid", func() {
				It("returns an error", func() {
					cfg.SecretsEncryption = &api.SecretsEncryption{
						KeyARN: "invalid:arn",
					}
					err := api.ValidateClusterConfig(cfg)
					Expect(err).To(MatchError(ContainSubstring("invalid ARN")))
				})
			})
		})

		When("a key ARN is not set", func() {
			It("returns an error", func() {
				cfg.SecretsEncryption = &api.SecretsEncryption{}
				err := api.ValidateClusterConfig(cfg)
				Expect(err).To(MatchError(ContainSubstring("field secretsEncryption.keyARN is required for enabling secrets encryption")))
			})
		})
	})
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
