package builder_test

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/builder/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	_ "embed"
)

var _ = Describe("Cluster Template Builder", func() {
	var (
		crs           *builder.ClusterResourceSet
		cfg           *api.ClusterConfig
		provider      *mockprovider.MockProvider
		existingStack *gjson.Result
	)

	BeforeEach(func() {
		provider = mockprovider.NewMockProvider()
		existingStack = nil
		cfg = api.NewClusterConfig()
		cfg.VPC = vpcConfig()
		cfg.AvailabilityZones = []string{"us-west-2a", "us-west-2b"}
		cfg.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
			ServiceIPv4CIDR: "131.10.55.70/18",
			IPFamily:        api.IPV4Family,
		}
	})

	JustBeforeEach(func() {
		crs = builder.NewClusterResourceSet(provider.EC2(), provider.Region(), cfg, existingStack, false)
	})

	Describe("AddAllResources", func() {
		var (
			addErr          error
			clusterTemplate *fakes.FakeTemplate
		)

		JustBeforeEach(func() {
			addErr = crs.AddAllResources(context.Background())
			clusterTemplate = &fakes.FakeTemplate{}
			templateBody, err := crs.RenderJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(json.Unmarshal(templateBody, clusterTemplate)).To(Succeed())
		})

		It("should not error", func() {
			Expect(addErr).NotTo(HaveOccurred())
		})

		It("should add a template description", func() {
			Expect(clusterTemplate.Description).To(Equal("EKS cluster (dedicated VPC: true, dedicated IAM: true) [created and managed by eksctl]"))
		})

		It("should add control plane resources", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlane"))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.Name).To(Equal(cfg.Metadata.Name))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.Version).To(Equal(cfg.Metadata.Version))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.ResourcesVpcConfig.SecurityGroupIds[0]).To(ContainElement("ControlPlaneSecurityGroup"))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.ResourcesVpcConfig.SubnetIds).To(HaveLen(4))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.RoleArn).To(ContainElement([]interface{}{"ServiceRole", "Arn"}))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.EncryptionConfig).To(BeNil())
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.KubernetesNetworkConfig.ServiceIPv4CIDR).To(Equal("131.10.55.70/18"))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.KubernetesNetworkConfig.IPFamily).To(Equal("ipv4"))
		})

		It("should add vpc resources", func() {
			Expect(clusterTemplate.Resources).To(HaveKey(vpcResourceKey))
			Expect(clusterTemplate.Resources).To(HaveKey(igwKey))
			Expect(clusterTemplate.Resources).To(HaveKey(gaKey))

			Expect(clusterTemplate.Resources).To(HaveKey(pubRouteTable))
			Expect(clusterTemplate.Resources).To(HaveKey(pubSubnetRoute))
			Expect(clusterTemplate.Resources).To(HaveKey(rtaPublicB))
			Expect(clusterTemplate.Resources).To(HaveKey(rtaPublicA))
			Expect(clusterTemplate.Resources).To(HaveKey(publicSubnetRef1))
			Expect(clusterTemplate.Resources).To(HaveKey(publicSubnetRef2))

			Expect(clusterTemplate.Resources).To(HaveKey(privRouteTableB))
			Expect(clusterTemplate.Resources).To(HaveKey(privRouteTableA))
			Expect(clusterTemplate.Resources).To(HaveKey(rtaPrivateA))
			Expect(clusterTemplate.Resources).To(HaveKey(rtaPrivateB))
			Expect(clusterTemplate.Resources).To(HaveKey(privateSubnetRef1))
			Expect(clusterTemplate.Resources).To(HaveKey(privateSubnetRef1))
		})

		Context("when ipFamily is set to IPv6", func() {
			BeforeEach(func() {
				cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
			})

			It("should add control plane resources", func() {
				Expect(clusterTemplate.Resources["ControlPlane"].Properties.KubernetesNetworkConfig.IPFamily).To(Equal("ipv6"))
			})

			It("should add IPv6 vpc resources", func() {
				Expect(clusterTemplate.Resources).To(HaveKey(builder.VPCResourceKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.IPv6CIDRBlockKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.IGWKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.GAKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.EgressOnlyInternetGatewayKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.NATGatewayKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.ElasticIPKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PubRouteTableKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PubSubRouteKey))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PubSubIPv6RouteKey))
				privateRouteTableA := builder.PrivateRouteTableKey + azAFormatted
				Expect(clusterTemplate.Resources).To(HaveKey(privateRouteTableA))
				privateRouteTableB := builder.PrivateRouteTableKey + azBFormatted
				Expect(clusterTemplate.Resources).To(HaveKey(privateRouteTableB))
				privateRouteA := builder.PrivateSubnetRouteKey + azAFormatted
				Expect(clusterTemplate.Resources).To(HaveKey(privateRouteA))
				privateRouteB := builder.PrivateSubnetRouteKey + azBFormatted
				Expect(clusterTemplate.Resources).To(HaveKey(privateRouteB))
				privateRouteA = builder.PrivateSubnetIpv6RouteKey + azAFormatted
				Expect(clusterTemplate.Resources).To(HaveKey(privateRouteA))
				privateRouteB = builder.PrivateSubnetIpv6RouteKey + azBFormatted
				Expect(clusterTemplate.Resources).To(HaveKey(privateRouteB))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PublicSubnetKey + azAFormatted))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PublicSubnetKey + azBFormatted))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateSubnetKey + azAFormatted))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateSubnetKey + azBFormatted))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PubRouteTableAssociation + azAFormatted))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PubRouteTableAssociation + azBFormatted))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateRouteTableAssociation + azAFormatted))
				Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateRouteTableAssociation + azBFormatted))
			})
		})

		Context("when AutoAllocateIPv6 is enabled", func() {
			BeforeEach(func() {
				autoAllocated := true
				cfg.VPC.AutoAllocateIPv6 = &autoAllocated
			})

			It("should add AutoAllocatedCIDRv6 vpc resource", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("AutoAllocatedCIDRv6"))
			})
		})

		Context("when NAT is enabled", func() {
			BeforeEach(func() {
				defaultNat := api.ClusterNATDefault
				cfg.VPC.NAT.Gateway = &defaultNat
			})

			It("should add nat vpc resources", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("NATIP"))
				Expect(clusterTemplate.Resources).To(HaveKey("NATGateway"))
			})
		})

		It("should add security group resources", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlaneSecurityGroup"))
			Expect(clusterTemplate.Resources).To(HaveKey("ClusterSharedNodeSecurityGroup"))
			Expect(clusterTemplate.Resources).To(HaveKey("IngressInterNodeGroupSG"))
			Expect(clusterTemplate.Resources).NotTo(HaveKey("IngressDefaultClusterToNodeSG"))
			Expect(clusterTemplate.Resources).NotTo(HaveKey("IngressNodeToDefaultClusterSG"))
			Expect(clusterTemplate.Resources).To(HaveKey("ClusterSharedNodeSecurityGroup"))
		})

		Context("when extraCIDRs are defined", func() {
			BeforeEach(func() {
				cfg.VPC.ExtraCIDRs = []string{"192.168.0.0/24", "192.168.1.0/24"}
			})

			It("should add extra control plane ingress rules", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraCIDR0"))

				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR0"].Properties).To(Equal(fakes.Properties{
					CidrIP:     "192.168.0.0/24",
					IPProtocol: "tcp",
					FromPort:   443,
					ToPort:     443,
					GroupID: map[string]interface{}{
						"Ref": "ControlPlaneSecurityGroup"},
					Description: "Allow Extra CIDR 0 (192.168.0.0/24) to communicate to controlplane",
				}))

				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraCIDR1"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR1"].Properties).To(Equal(fakes.Properties{
					CidrIP:     "192.168.1.0/24",
					IPProtocol: "tcp",
					FromPort:   443,
					ToPort:     443,
					GroupID: map[string]interface{}{
						"Ref": "ControlPlaneSecurityGroup"},
					Description: "Allow Extra CIDR 1 (192.168.1.0/24) to communicate to controlplane",
				}))
			})
		})

		Context("when extraIPv6CIDR is defined", func() {
			BeforeEach(func() {
				cfg.VPC.ExtraIPv6CIDRs = []string{"2002::1234:abcd:ffff:c0a8:101/64", "2003::1234:abcd:ffff:c0a8:101/64"}
			})

			It("should add extra control plane ingress rules", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraIPv6CIDR0"))

				Expect(clusterTemplate.Resources["IngressControlPlaneExtraIPv6CIDR0"].Properties).To(Equal(fakes.Properties{
					CidrIPv6:   "2002::1234:abcd:ffff:c0a8:101/64",
					IPProtocol: "tcp",
					FromPort:   443,
					ToPort:     443,
					GroupID: map[string]interface{}{
						"Ref": "ControlPlaneSecurityGroup"},
					Description: "Allow Extra IPv6 CIDR 0 (2002::1234:abcd:ffff:c0a8:101/64) to communicate to controlplane",
				}))

				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraIPv6CIDR1"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraIPv6CIDR1"].Properties).To(Equal(fakes.Properties{
					CidrIPv6:   "2003::1234:abcd:ffff:c0a8:101/64",
					IPProtocol: "tcp",
					FromPort:   443,
					ToPort:     443,
					GroupID: map[string]interface{}{
						"Ref": "ControlPlaneSecurityGroup"},
					Description: "Allow Extra IPv6 CIDR 1 (2003::1234:abcd:ffff:c0a8:101/64) to communicate to controlplane",
				}))
			})

			Context("when managed nodegroups are configured is true", func() {
				BeforeEach(func() {
					enabled := true
					cfg.VPC.ManageSharedNodeSecurityGroupRules = &enabled
				})

				It("sets IngressDefaultClusterToNodeSG and IngressNodeToDefaultClusterSG resources", func() {
					Expect(clusterTemplate.Resources).To(HaveKey("IngressDefaultClusterToNodeSG"))
					Expect(clusterTemplate.Resources).To(HaveKey("IngressNodeToDefaultClusterSG"))
				})
			})
		})

		Context("if SharedNodeSecurityGroup is set", func() {
			BeforeEach(func() {
				cfg.VPC.SharedNodeSecurityGroup = "foo"
			})

			It("should not add various shared security group resources", func() {
				Expect(clusterTemplate.Resources).NotTo(HaveKey("ClusterSharedNodeSecurityGroup"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("IngressInterNodeGroupSG"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("IngressDefaultClusterToNodeSG"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("IngressNodeToDefaultClusterSG"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("ClusterSharedNodeSecurityGroup"))
			})
		})

		Context("if the control plane SecurityGroup is set", func() {
			BeforeEach(func() {
				cfg.VPC.SecurityGroup = "foo"
			})

			It("should not add the ControlPlaneSecurityGroup resources", func() {
				Expect(clusterTemplate.Resources).NotTo(HaveKey("ControlPlaneSecurityGroup"))
				Expect(clusterTemplate.Resources["ControlPlane"].Properties.ResourcesVpcConfig.SecurityGroupIds).To(ContainElement("foo"))
			})
		})

		It("should add iam resources and policies", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("ServiceRole"))
			Expect(clusterTemplate.Resources).To(HaveKey("PolicyELBPermissions"))
			Expect(clusterTemplate.Resources).To(HaveKey("PolicyCloudWatchMetrics"))
		})

		It("should add the correct policies and references to the ServiceRole ARN", func() {
			Expect(clusterTemplate.Resources["ServiceRole"].Properties.ManagedPolicyArns).To(HaveLen(2))
			Expect(clusterTemplate.Resources["ServiceRole"].Properties.ManagedPolicyArns).To(ContainElements(makePolicyARNRef("AmazonEKSClusterPolicy"), makePolicyARNRef("AmazonEKSVPCResourceController")))

			cwPolicy := clusterTemplate.Resources["PolicyCloudWatchMetrics"].Properties
			Expect(isRefTo(cwPolicy.Roles[0], "ServiceRole")).To(BeTrue())
			elbPolicy := clusterTemplate.Resources["PolicyELBPermissions"].Properties
			Expect(isRefTo(elbPolicy.Roles[0], "ServiceRole")).To(BeTrue())
		})

		It("should add iam outputs", func() {
			Expect(clusterTemplate.Outputs).To(HaveKey("ServiceRoleARN"))
		})

		Context("when ServiceRoleARN is set", func() {
			BeforeEach(func() {
				role := "foo"
				cfg.IAM.ServiceRoleARN = &role
			})

			It("should not add other iam resources", func() {
				Expect(clusterTemplate.Resources).NotTo(HaveKey("ServiceRole"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("PolicyELBPermissions"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("PolicyCloudWatchMetrics"))
			})
		})

		Context("when ServiceRolePermissionsBoundary is set", func() {
			BeforeEach(func() {
				pb := "foo"
				cfg.IAM.ServiceRolePermissionsBoundary = &pb
			})

			It("adds the permissions boundary to the service role", func() {
				Expect(clusterTemplate.Resources["ServiceRole"].Properties.PermissionsBoundary).To(Equal("foo"))
			})
		})

		Context("when VPCResourceControllerPolicy is disabled", func() {
			BeforeEach(func() {
				policy := false
				cfg.IAM.VPCResourceControllerPolicy = &policy
			})

			It("only adds the AmazonEKSClusterPolicy to the service role policy arn", func() {
				Expect(clusterTemplate.Resources["ServiceRole"].Properties.ManagedPolicyArns).To(HaveLen(1))
				Expect(clusterTemplate.Resources["ServiceRole"].Properties.ManagedPolicyArns[0]).To(Equal(makePolicyARNRef("AmazonEKSClusterPolicy")))
			})
		})

		When("SecretsEncryption is configured", func() {
			BeforeEach(func() {
				cfg.SecretsEncryption = &api.SecretsEncryption{
					KeyARN: "key-thing",
				}
			})

			It("should add the key arn to the control plane resource", func() {
				Expect(clusterTemplate.Resources["ControlPlane"].Properties.EncryptionConfig[0].Provider.KeyARN).To(Equal("key-thing"))
				Expect(clusterTemplate.Resources["ControlPlane"].Properties.EncryptionConfig[0].Resources[0]).To(Equal("secrets"))
			})
		})

		It("should add cluster stack outputs", func() {
			Expect(clusterTemplate.Outputs).To(HaveLen(12))
			Expect(clusterTemplate.Outputs).To(HaveKey("ARN"))
			Expect(clusterTemplate.Outputs).To(HaveKey("ClusterStackName"))
			Expect(clusterTemplate.Outputs).To(HaveKey("SecurityGroup"))
			Expect(clusterTemplate.Outputs).To(HaveKey("SharedNodeSecurityGroup"))
			Expect(clusterTemplate.Outputs).To(HaveKey("SubnetsPrivate"))
			Expect(clusterTemplate.Outputs).To(HaveKey("SubnetsPublic"))
			Expect(clusterTemplate.Outputs).To(HaveKey(vpcResourceKey))
			Expect(clusterTemplate.Outputs).To(HaveKey("CertificateAuthorityData"))
			Expect(clusterTemplate.Outputs).To(HaveKey("Endpoint"))
			Expect(clusterTemplate.Outputs).To(HaveKey("FeatureNATMode"))
			Expect(clusterTemplate.Outputs).To(HaveKey("ServiceRoleARN"))
			Expect(clusterTemplate.Outputs).To(HaveKey("ClusterSecurityGroupId"))
		})

		It("should add partition mappings", func() {
			Expect(clusterTemplate.Mappings["ServicePrincipalPartitionMap"]).NotTo(BeNil())
		})

		Context("when private networking is set", func() {
			BeforeEach(func() {
				cfg.PrivateCluster = &api.PrivateCluster{Enabled: true}

				detailsJSON := serviceDetailsJSON
				var output *ec2.DescribeVpcEndpointServicesOutput
				Expect(json.Unmarshal([]byte(detailsJSON), &output)).To(Succeed())
				provider.MockEC2().On("DescribeVpcEndpointServices", mock.Anything, mock.MatchedBy(func(e *ec2.DescribeVpcEndpointServicesInput) bool {
					return len(e.ServiceNames) == 5
				})).Return(output, nil)
			})

			It("the correct VPC endpoint resources are added", func() {
				Expect(clusterTemplate.Resources).To(HaveKey(ContainSubstring("VPCEndpoint")))
			})

			It("adds the ClusterFullyPrivate output", func() {
				Expect(clusterTemplate.Outputs).To(HaveKey("ClusterFullyPrivate"))
			})

			It("no NAT resources are set", func() {
				Expect(clusterTemplate.Resources).NotTo(HaveKey("NATIP"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("NATGateway"))
			})

			It("does not set public networking", func() {
				Expect(clusterTemplate.Resources).NotTo(HaveKey("PublicSubnetRoute"))
				Expect(clusterTemplate.Resources).NotTo(HaveKey("PublicSubnetRoute"))
				Expect(clusterTemplate.Resources).To(HaveKey(ContainSubstring("PrivateRouteTable")))
			})
			When("skip endpoint creation is set", func() {
				BeforeEach(func() {
					cfg.PrivateCluster = &api.PrivateCluster{
						Enabled:              true,
						SkipEndpointCreation: true,
					}
				})
				It("will skip creating all of the endpoints", func() {
					Expect(clusterTemplate.Resources).NotTo(HaveKey(ContainSubstring("VPCEndpoint")))
				})
			})

			When("ipv6 cluster is enabled", func() {
				BeforeEach(func() {
					cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
				})
				It("should only add private IPv6 vpc resources", func() {
					Expect(clusterTemplate.Resources).To(HaveKey(builder.VPCResourceKey))
					Expect(clusterTemplate.Resources).To(HaveKey(builder.IPv6CIDRBlockKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.IGWKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.GAKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.EgressOnlyInternetGatewayKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.NATGatewayKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.ElasticIPKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.PubRouteTableKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.PubSubRouteKey))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.PubSubIPv6RouteKey))
					privateRouteTableA := builder.PrivateRouteTableKey + azAFormatted
					Expect(clusterTemplate.Resources).To(HaveKey(privateRouteTableA))
					privateRouteTableB := builder.PrivateRouteTableKey + azBFormatted
					Expect(clusterTemplate.Resources).To(HaveKey(privateRouteTableB))
					privateRouteA := builder.PrivateSubnetRouteKey + azAFormatted
					Expect(clusterTemplate.Resources).NotTo(HaveKey(privateRouteA))
					privateRouteB := builder.PrivateSubnetRouteKey + azBFormatted
					Expect(clusterTemplate.Resources).NotTo(HaveKey(privateRouteB))
					privateRouteA = builder.PrivateSubnetIpv6RouteKey + azAFormatted
					Expect(clusterTemplate.Resources).NotTo(HaveKey(privateRouteA))
					privateRouteB = builder.PrivateSubnetIpv6RouteKey + azBFormatted
					Expect(clusterTemplate.Resources).NotTo(HaveKey(privateRouteB))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.PublicSubnetKey + azAFormatted))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.PublicSubnetKey + azBFormatted))
					Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateSubnetKey + azAFormatted))
					Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateSubnetKey + azBFormatted))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.PubRouteTableAssociation + azAFormatted))
					Expect(clusterTemplate.Resources).NotTo(HaveKey(builder.PubRouteTableAssociation + azBFormatted))
					Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateRouteTableAssociation + azAFormatted))
					Expect(clusterTemplate.Resources).To(HaveKey(builder.PrivateRouteTableAssociation + azBFormatted))
				})
			})

			When("skip endpoint creation is set", func() {
				BeforeEach(func() {
					cfg.PrivateCluster = &api.PrivateCluster{
						Enabled:              true,
						SkipEndpointCreation: true,
					}
				})
				It("will skip creating all of the endpoints", func() {
					Expect(clusterTemplate.Resources).NotTo(HaveKey(ContainSubstring("VPCEndpoint")))
				})
			})
		})

		Context("when default config is used", func() {
			It("should not enable non-default features", func() {
				cluster := clusterTemplate.Resources["ControlPlane"].Properties
				By("ensuring logging is not enabled")
				Expect(cluster.Logging.ClusterLogging.EnabledTypes).To(BeEmpty())

				vpcResources := cluster.ResourcesVpcConfig
				By("ensuring public access is enabled but private access is not")
				Expect(vpcResources.EndpointPublicAccess).To(BeTrue())
				Expect(vpcResources.EndpointPrivateAccess).To(BeFalse())

				By("ensuring publicAccessCIDRs is not enabled")
				Expect(vpcResources.PublicAccessCidrs).To(BeEmpty())
			})
		})

		Context("clusterEndpoints.privateAccess is enabled", func() {
			BeforeEach(func() {
				cfg.VPC.ClusterEndpoints.PrivateAccess = api.Enabled()
				cfg.VPC.ClusterEndpoints.PublicAccess = api.Disabled()
			})
			It("should enable privateAccess in the stack", func() {
				vpcResources := clusterTemplate.Resources["ControlPlane"].Properties.ResourcesVpcConfig
				Expect(vpcResources.EndpointPrivateAccess).To(BeTrue())
				Expect(vpcResources.EndpointPublicAccess).To(BeFalse())
			})
		})

		Context("vpc.publicAccessCIDRs is set", func() {
			BeforeEach(func() {
				cfg.VPC.PublicAccessCIDRs = []string{"17.0.0.0/8", "73.0.0.0/8"}
			})
			It("should set the supplied CIDRs in the stack", func() {
				vpcResources := clusterTemplate.Resources["ControlPlane"].Properties.ResourcesVpcConfig
				Expect(vpcResources.PublicAccessCidrs).To(Equal([]string{"17.0.0.0/8", "73.0.0.0/8"}))
			})
		})

		Context("cluster tags are set", func() {
			BeforeEach(func() {
				cfg.Metadata.Tags = map[string]string{
					"type": "production",
					"key":  "value",
				}
			})
			It("should set tags in the stack", func() {
				Expect(clusterTemplate.Resources["ControlPlane"].Properties.Tags).To(ConsistOf([]fakes.Tag{
					{
						Key:   "type",
						Value: "production",
					},
					{
						Key:   "key",
						Value: "value",
					},
					{
						Key: "Name",
						Value: map[string]interface{}{
							"Fn::Sub": "${AWS::StackName}/ControlPlane",
						},
					},
				}))
			})
		})

		Context("cluster logging is enabled", func() {
			BeforeEach(func() {
				cfg.CloudWatch.ClusterLogging = &api.ClusterCloudWatchLogging{
					EnableTypes: []string{"api", "audit", "scheduler", "controllerManager"},
				}
			})

			It("should have logging enabled in the stack", func() {
				Expect(clusterTemplate.Resources["ControlPlane"].Properties.Logging.ClusterLogging.EnabledTypes).To(Equal([]fakes.ClusterLoggingType{
					{
						Type: "api",
					},
					{
						Type: "audit",
					},
					{
						Type: "scheduler",
					},
					{
						Type: "controllerManager",
					},
				}))
			})
		})

		Context("when adding vpc endpoint resources fails", func() {
			BeforeEach(func() {
				cfg.PrivateCluster = &api.PrivateCluster{Enabled: true}
				provider.MockEC2().On("DescribeVpcEndpointServices", mock.Anything, mock.Anything).Return(nil, errors.New("o-noes"))
			})

			It("should return the error", func() {
				Expect(addErr).To(MatchError(ContainSubstring("error describing VPC endpoint services")))
			})
		})

		Context("when fargate profiles are configured", func() {
			BeforeEach(func() {
				cfg.Metadata.AccountID = "111122223333"
				cfg.FargateProfiles = []*api.FargateProfile{{
					Name: "fp-default",
					Selectors: []api.FargateProfileSelector{
						{Namespace: "default"},
					},
				}}
			})

			It("should add resources for fargate", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("FargatePodExecutionRole"))
			})
		})

		Context("when the spec has insufficient subnets", func() {
			BeforeEach(func() {
				cfg.VPC.Subnets = &api.ClusterSubnets{}
			})

			It("should fail", func() {
				Expect(addErr).To(MatchError(HaveSuffix("insufficient number of subnets, at least 2x public and/or 2x private subnets are required")))
			})
		})

		Context("[Outposts] when the spec has insufficient subnets", func() {
			BeforeEach(func() {
				cfg.VPC.Subnets = &api.ClusterSubnets{}
				cfg.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			})

			It("should fail", func() {
				Expect(addErr).To(MatchError(HaveSuffix("insufficient number of subnets, at least 1x public and/or 1x private subnets are required for Outposts")))
			})
		})

		Context("[Outposts] when the cluster is fully-private", func() {
			BeforeEach(func() {
				cfg.PrivateCluster = &api.PrivateCluster{
					Enabled: true,
				}
				az := cfg.AvailabilityZones[0]
				cfg.VPC.ManageSharedNodeSecurityGroupRules = api.Enabled()
				subnet := cfg.VPC.Subnets.Private[az]
				subnet.ID = ""

				cfg.VPC.Subnets.Private = api.AZSubnetMapping{
					az: subnet,
				}
				cfg.Outpost = &api.Outpost{
					ControlPlaneOutpostARN:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
					ControlPlaneInstanceType: "m5.large",
				}

				var output *ec2.DescribeVpcEndpointServicesOutput
				Expect(json.Unmarshal(serviceDetailsOutpostsJSON, &output)).To(Succeed())
				provider.MockEC2().On("DescribeVpcEndpointServices", mock.Anything, mock.MatchedBy(func(e *ec2.DescribeVpcEndpointServicesInput) bool {
					return reflect.DeepEqual(e.ServiceNames, output.ServiceNames)
				})).Return(output, nil)
			})

			It("should create a security group for ingress traffic from private subnet CIDRs", func() {
				const ingressRuleKey = "IngressPrivateSubnetUSWEST2A"
				Expect(clusterTemplate.Resources).To(HaveKey(ingressRuleKey))
			})
		})

		Context("when adding vpc resources fails", func() {
			BeforeEach(func() {
				cfg.VPC = &api.ClusterVPC{}
				cfg.VPC.ID = "kitten"
			})

			It("should return the error", func() {
				Expect(addErr).To(MatchError(ContainSubstring("insufficient number of subnets")))
			})
		})
	})

	Describe("GetAllOutputs", func() {
		It("should not error", func() {
			// the actual work gets done right the way down in outputs where there is currently no interface
			// so there is little value here right now
			Expect(crs.GetAllOutputs(types.Stack{})).To(Succeed())
		})
	})

	Describe("RenderJSON", func() {
		It("returns the template rendered as JSON", func() {
			// the work actually gets done on the internal resource set
			Expect(crs.AddAllResources(context.Background())).To(Succeed())
			result, err := crs.RenderJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(vpcResourceKey))
		})
	})

	Describe("Template", func() {
		It("returns the template from the inner resource set", func() {
			// the work actually gets done on the internal resource set
			Expect(crs.AddAllResources(context.Background())).To(Succeed())
			clusterTemplate := crs.Template()
			Expect(clusterTemplate.Resources).To(HaveKey(vpcResourceKey))
		})
	})
})

func makePolicyARNRef(policy string) map[string]interface{} {
	return map[string]interface{}{
		"Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/" + policy,
	}
}

//go:embed testdata/service_details_outposts.json
var serviceDetailsOutpostsJSON []byte

var serviceDetailsJSON = `
{
  "ServiceNames": [ "com.amazonaws.us-west-2.ec2" ],
  "ServiceDetails": [
    {
      "ServiceType": [ { "ServiceType": "Interface" } ],
      "ServiceName": "com.amazonaws.us-west-2.ec2",
      "BaseEndpointDnsNames": [ "ec2.us-west-2.vpce.amazonaws.com" ]
    }
  ]
}
`
