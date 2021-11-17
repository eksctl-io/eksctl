package builder_test

import (
	"encoding/json"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/builder/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Cluster Template Builder", func() {
	var (
		crs                  *builder.ClusterResourceSet
		cfg                  *api.ClusterConfig
		provider             *mockprovider.MockProvider
		supportsManagedNodes bool
		existingStack        *gjson.Result
	)

	BeforeEach(func() {
		provider = mockprovider.NewMockProvider()
		supportsManagedNodes = false
		existingStack = nil
		cfg = api.NewClusterConfig()
		cfg.VPC = vpcConfig()
		cfg.AvailabilityZones = []string{"us-west-2a", "us-west-2b"}
	})

	JustBeforeEach(func() {
		crs = builder.NewClusterResourceSet(provider.EC2(), provider.Region(), cfg, supportsManagedNodes, existingStack)
	})

	Describe("AddAllResources", func() {
		var (
			addErr          error
			clusterTemplate *fakes.FakeTemplate
		)

		JustBeforeEach(func() {
			addErr = crs.AddAllResources()
			clusterTemplate = &fakes.FakeTemplate{}
			templateBody, err := crs.RenderJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(json.Unmarshal(templateBody, clusterTemplate)).To(Succeed())
		})

		It("should not error", func() {
			Expect(addErr).NotTo(HaveOccurred())
		})

		It("should add a template description", func() {
			Expect(clusterTemplate.Description).To(Equal("EKS cluster (dedicated VPC: true, dedicated IAM: true) [created and managed by eksctl]"))
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

		Context("when 1 extraCIDR is defined", func() {
			BeforeEach(func() {
				oneExtraCIDR := []string{"192.168.0.0/24"}
				cfg.VPC.ExtraCIDRs = oneExtraCIDR
			})

			It("should add 1 extra control plane ingress rule", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraCIDR0"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR0"].Properties.CidrIP).To(Equal("192.168.0.0/24"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR0"].Properties.IPProtocol).To(Equal("tcp"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR0"].Properties.FromPort).To(Equal(443))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR0"].Properties.ToPort).To(Equal(443))
			})
		})

		Context("when 3 extraCIDRs are defined", func() {
			BeforeEach(func() {
				threeExtraCIDRs := []string{"192.168.0.0/24", "192.168.1.0/24", "192.168.2.0/24"}
				cfg.VPC.ExtraCIDRs = threeExtraCIDRs
			})

			It("should add 3 extra control plane ingress rules", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraCIDR0"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR0"].Properties.CidrIP).To(Equal("192.168.0.0/24"))
				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraCIDR1"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR1"].Properties.CidrIP).To(Equal("192.168.1.0/24"))
				Expect(clusterTemplate.Resources).To(HaveKey("IngressControlPlaneExtraCIDR2"))
				Expect(clusterTemplate.Resources["IngressControlPlaneExtraCIDR2"].Properties.CidrIP).To(Equal("192.168.2.0/24"))
			})
		})

		Context("when supportsManagedNodes is true", func() {
			BeforeEach(func() {
				supportsManagedNodes = true
				enabled := true
				cfg.VPC.ManageSharedNodeSecurityGroupRules = &enabled
			})

			It("sets IngressDefaultClusterToNodeSG and IngressNodeToDefaultClusterSG resources", func() {
				Expect(clusterTemplate.Resources).To(HaveKey("IngressDefaultClusterToNodeSG"))
				Expect(clusterTemplate.Resources).To(HaveKey("IngressNodeToDefaultClusterSG"))
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
			Expect(clusterTemplate.Resources["ServiceRole"].Properties.ManagedPolicyArns[0]).To(Equal(makePolicyARNRef("AmazonEKSClusterPolicy")))
			Expect(clusterTemplate.Resources["ServiceRole"].Properties.ManagedPolicyArns[1]).To(Equal(makePolicyARNRef("AmazonEKSVPCResourceController")))

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

		It("should add control plane resources", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlane"))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.Name).To(Equal(cfg.Metadata.Name))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.Version).To(Equal(cfg.Metadata.Version))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.ResourcesVpcConfig.SecurityGroupIds[0]).To(ContainElement("ControlPlaneSecurityGroup"))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.ResourcesVpcConfig.SubnetIds).To(HaveLen(4))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.RoleArn).To(ContainElement([]interface{}{"ServiceRole", "Arn"}))
			Expect(clusterTemplate.Resources["ControlPlane"].Properties.EncryptionConfig).To(BeNil())
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
			Expect(clusterTemplate.Outputs).To(HaveLen(11))
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
				provider.MockEC2().On("DescribeVpcEndpointServices", mock.MatchedBy(func(e *ec2.DescribeVpcEndpointServicesInput) bool {
					return len(e.ServiceNames) == 5
				})).Return(output, nil)
			})

			It("the correct vpc endpoint resources are added", func() {
				Expect(clusterTemplate.Resources).To(HaveKey(ContainSubstring("VPCEndpoint")))
			})

			It("adds the ClusterFullyPrivate output", func() {
				Expect(clusterTemplate.Outputs).To(HaveKey("ClusterFullyPrivate"))
			})

			It("no nat resources are set", func() {
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
		})

		Context("when adding vpc endpoint resources fails", func() {
			BeforeEach(func() {
				cfg.PrivateCluster = &api.PrivateCluster{Enabled: true}
				provider.MockEC2().On("DescribeVpcEndpointServices", mock.Anything).Return(nil, errors.New("o-noes"))
			})

			It("should return the error", func() {
				Expect(addErr).To(MatchError(ContainSubstring("error describing VPC endpoint services")))
			})
		})

		Context("when fargate profiles are configured", func() {
			BeforeEach(func() {
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
				Expect(addErr).To(MatchError(ContainSubstring("insufficient number of subnets")))
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
			Expect(crs.GetAllOutputs(cfn.Stack{})).To(Succeed())
		})
	})

	Describe("RenderJSON", func() {
		It("returns the template rendered as JSON", func() {
			// the work actually gets done on the internal resource set
			result, err := crs.RenderJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(vpcResourceKey))
		})
	})

	Describe("Template", func() {
		It("returns the template from the inner resource set", func() {
			// the work actually gets done on the internal resource set
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
