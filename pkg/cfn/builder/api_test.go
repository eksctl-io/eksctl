package builder_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/weaveworks/goformation/v4"
	"k8s.io/apimachinery/pkg/util/sets"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/cfn/builder"
	bootstrapfakes "github.com/weaveworks/eksctl/pkg/nodebootstrap/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/vpc"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"
)

const (
	clusterName = "ferocious-mushroom-1532594698"
	endpoint    = "https://DE37D8AFB23F7275D2361AD6B2599143.yl4.us-west-2.eks.amazonaws.com"
	caCert      = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNE1EWXdOekExTlRBMU5Wb1hEVEk0TURZd05EQTFOVEExTlZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTWJoCnpvZElYR0drckNSZE1jUmVEN0YvMnB1NFZweTdvd3FEVDgrdk9zeGs2bXFMNWxQd3ZicFhmYkE3R0xzMDVHa0wKaDdqL0ZjcU91cnMwUFZSK3N5REtuQXltdDFORWxGNllGQktSV1dUQ1hNd2lwN1pweW9XMXdoYTlJYUlPUGxCTQpPTEVlckRabFVrVDFVV0dWeVdsMmxPeFgxa2JhV2gvakptWWdkeW5jMXhZZ3kxa2JybmVMSkkwLzVUVTRCajJxClB1emtrYW5Xd3lKbGdXQzhBSXlpWW82WFh2UVZmRzYrM3RISE5XM1F1b3ZoRng2MTFOYnl6RUI3QTdtZGNiNmgKR0ZpWjdOeThHZnFzdjJJSmI2Nk9FVzBSdW9oY1k3UDZPdnZmYnlKREhaU2hqTStRWFkxQXN5b3g4Ri9UelhHSgpQUWpoWUZWWEVhZU1wQmJqNmNFQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFCa2hKRVd4MHk1LzlMSklWdXJ1c1hZbjN6Z2EKRkZ6V0JsQU44WTlqUHB3S2t0Vy9JNFYyUGg3bWY2Z3ZwZ3Jhc2t1Slk1aHZPcDdBQmcxSTFhaHUxNUFpMUI0ZApuMllRaDlOaHdXM2pKMmhuRXk0VElpb0gza2JFdHRnUVB2bWhUQzNEYUJreEpkbmZJSEJCV1RFTTU1czRwRmxUClpzQVJ3aDc1Q3hYbjdScVU0akpKcWNPaTRjeU5qeFVpRDBqR1FaTmNiZWEyMkRCeTJXaEEzUWZnbGNScGtDVGUKRDVPS3NOWlF4MW9MZFAwci9TSmtPT1NPeUdnbVJURTIrODQxN21PRW02Z3RPMCszdWJkbXQ0aENsWEtFTTZYdwpuQWNlK0JxVUNYblVIN2ZNS3p2TDE5UExvMm5KbFU1TnlCbU1nL1pNVHVlUy80eFZmKy94WnpsQ0Q1WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	arn         = "arn:aws:eks:us-west-2:122333:cluster/" + clusterName

	fargatePodExecutionRoleARN = "arn:aws:iam::123:role/" + clusterName + "-FargatePodExecutionRole-XYZ"

	vpcID          = "vpc-0e265ad953062b94b"
	subnetsPublic  = "subnet-0f98135715dfcf55f,subnet-0ade11bad78dced9e,subnet-0e2e63ff1712bf6ef"
	subnetsPrivate = "subnet-0f98135715dfcf55a,subnet-0ade11bad78dced9f,subnet-0e2e63ff1712bf6ea"

	p4InstanceType = "p4d.24xlarge"
)

type Tag struct {
	Key   interface{}
	Value interface{}

	PropagateAtLaunch string
}

type Properties struct {
	Tags []Tag

	Path, RoleName           string
	Roles, ManagedPolicyArns []interface{}
	AssumeRolePolicyDocument interface{}

	PolicyDocument struct {
		Statement []struct {
			Action    []string
			Effect    string
			Resource  interface{}
			Condition map[string]interface{}
		}
	}

	LaunchTemplateData LaunchTemplateData

	VPCZoneIdentifier interface{}

	LoadBalancerNames                 []string
	MetricsCollection                 []map[string]interface{}
	TargetGroupARNs                   []string
	DesiredCapacity, MinSize, MaxSize string

	CidrIP, CidrIpv6, IPProtocol string
	FromPort, ToPort             int

	VpcID, SubnetID                            interface{}
	RouteTableID, AllocationID                 interface{}
	GatewayID, InternetGatewayID, NatGatewayID interface{}
	DestinationCidrBlock                       interface{}

	Ipv6CidrBlock map[string][]interface{}

	AmazonProvidedIpv6CidrBlock         bool
	AvailabilityZone, Domain, CidrBlock string

	Name, Version      string
	RoleArn            interface{}
	ResourcesVpcConfig struct {
		SecurityGroupIds []interface{}
		SubnetIds        []interface{}
	}
	MixedInstancesPolicy *struct {
		LaunchTemplate struct {
			LaunchTemplateSpecification struct {
				LaunchTemplateName map[string]interface{}
				Version            map[string]interface{}
				Overrides          []struct {
					InstanceType string
				}
			}
		}
		InstancesDistribution struct {
			OnDemandBaseCapacity                string
			OnDemandPercentageAboveBaseCapacity string
			SpotMaxPrice                        string
			SpotInstancePools                   string
			SpotAllocationStrategy              string
		}
	}
}

type NetworkInterface struct {
	DeviceIndex              int
	AssociatePublicIPAddress bool
	NetworkCardIndex         int
	InterfaceType            string
}
type LaunchTemplateData struct {
	IamInstanceProfile              struct{ Arn interface{} }
	UserData, InstanceType, ImageID string
	BlockDeviceMappings             []interface{}
	EbsOptimized                    *bool
	NetworkInterfaces               []NetworkInterface
	InstanceMarketOptions           *struct {
		MarketType  string
		SpotOptions struct {
			SpotInstanceType string
			MaxPrice         string
		}
	}
	CreditSpecification *struct {
		CPUCredits string
	}
}

type Template struct {
	Description string
	Resources   map[string]struct {
		Properties   Properties
		DependsOn    []string
		UpdatePolicy map[string]map[string]interface{}
	}
}

var appMeshActions = []string{
	"servicediscovery:CreateService",
	"servicediscovery:DeleteService",
	"servicediscovery:GetService",
	"servicediscovery:GetInstance",
	"servicediscovery:RegisterInstance",
	"servicediscovery:DeregisterInstance",
	"servicediscovery:ListInstances",
	"servicediscovery:ListNamespaces",
	"servicediscovery:ListServices",
	"servicediscovery:GetInstancesHealthStatus",
	"servicediscovery:UpdateInstanceCustomHealthStatus",
	"servicediscovery:GetOperation",
	"route53:GetHealthCheck",
	"route53:CreateHealthCheck",
	"route53:UpdateHealthCheck",
	"route53:ChangeResourceRecordSets",
	"route53:DeleteHealthCheck",
}

func testVPC() *api.ClusterVPC {
	disable := api.ClusterDisableNAT
	return &api.ClusterVPC{
		Network: api.Network{
			ID: vpcID,
			CIDR: &ipnet.IPNet{
				IPNet: net.IPNet{
					IP:   []byte{192, 168, 0, 0},
					Mask: []byte{255, 255, 0, 0},
				},
			},
		},
		NAT: &api.ClusterNAT{
			Gateway: &disable,
		},
		SecurityGroup:                      "sg-0b44c48bcba5b7362",
		SharedNodeSecurityGroup:            "sg-shared",
		ManageSharedNodeSecurityGroupRules: api.Enabled(),
		AutoAllocateIPv6:                   api.Disabled(),
		Subnets: &api.ClusterSubnets{
			Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				"us-west-2b": {
					ID: "subnet-0f98135715dfcf55f",
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 0, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				"us-west-2a": {
					ID: "subnet-0ade11bad78dced9e",
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 32, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				"us-west-2c": {
					ID: "subnet-0e2e63ff1712bf6ef",
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 64, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
			}),
			Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				"us-west-2b": {
					ID: "subnet-0f98135715dfcf55a",
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 96, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				"us-west-2a": {
					ID: "subnet-0ade11bad78dced9f",
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 128, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				"us-west-2c": {
					ID: "subnet-0e2e63ff1712bf6ea",
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 160, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
			}),
		},
		ClusterEndpoints: &api.ClusterEndpoints{},
	}
}

var subnetLists = map[api.SubnetTopology]string{
	"Public":  subnetsPublic,
	"Private": subnetsPrivate,
}

func newStackWithOutputs(outputs map[string]string) cfn.Stack {
	s := cfn.Stack{}
	for k, v := range outputs {
		func(k, v string) {
			s.Outputs = append(s.Outputs,
				&cfn.Output{
					OutputKey:   &k,
					OutputValue: &v,
				})
		}(k, v)
	}
	return s
}

var _ = Describe("CloudFormation template builder API", func() {
	var (
		crs  *ClusterResourceSet
		ngrs *NodeGroupResourceSet

		clusterTemplate, ngTemplate *Template

		err error

		caCertData []byte
	)

	Describe("should decode CA data", func() {
		caCertData, err = base64.StdEncoding.DecodeString(caCert)
		It("should not error", func() { Expect(err).ShouldNot(HaveOccurred()) })
	})

	testAZs := []string{"us-west-2b", "us-west-2a", "us-west-2c"}

	newClusterConfigAndNodegroup := func(withFullVPC bool) (*api.ClusterConfig, *api.NodeGroup) {
		cfg := api.NewClusterConfig()
		ng := cfg.NewNodeGroup()

		cfg.Metadata.Region = "us-west-2"
		cfg.Metadata.Name = clusterName

		cfg.Status = &api.ClusterStatus{
			CertificateAuthorityData: caCertData,
			Endpoint:                 endpoint,
		}

		cfg.AvailabilityZones = testAZs
		ng.Name = "ng-abcd1234"
		ng.InstanceType = "t2.medium"
		ng.AMIFamily = "AmazonLinux2"
		ng.VolumeSize = new(int)
		*ng.VolumeSize = 2
		ng.VolumeType = new(string)
		*ng.VolumeType = api.NodeVolumeTypeSC1
		ng.VolumeName = new(string)
		*ng.VolumeName = "/dev/xvda"
		ng.VolumeEncrypted = api.Disabled()
		ng.VolumeKmsKeyID = new(string)

		if withFullVPC {
			cfg.VPC = testVPC()
		} else {
			*cfg.VPC.CIDR = api.DefaultCIDR()
		}

		return cfg, ng
	}

	newSimpleClusterConfig := func() *api.ClusterConfig {
		cfg, _ := newClusterConfigAndNodegroup(false)
		return cfg
	}

	p := mockprovider.NewMockProvider()

	{
		joinCompare := func(input *ec2.DescribeSubnetsInput, compare string) bool {
			ids := make([]string, len(input.SubnetIds))
			for x, id := range input.SubnetIds {
				ids[x] = *id
			}
			return strings.Join(ids, ",") == compare
		}
		testVPC := testVPC()

		p.MockEC2().On("DescribeInstanceTypes",
			&ec2.DescribeInstanceTypesInput{
				InstanceTypes: aws.StringSlice([]string{p4InstanceType}),
			},
		).Return(
			&ec2.DescribeInstanceTypesOutput{
				InstanceTypes: []*ec2.InstanceTypeInfo{
					{
						InstanceType: aws.String(p4InstanceType),
						NetworkInfo: &ec2.NetworkInfo{
							EfaSupported:        aws.Bool(true),
							MaximumNetworkCards: aws.Int64(4),
						},
					},
				},
			}, nil,
		)

		p.MockEC2().On("DescribeVpcs", mock.MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
			return *input.VpcIds[0] == vpcID
		})).Return(&ec2.DescribeVpcsOutput{
			Vpcs: []*ec2.Vpc{{
				VpcId:     aws.String(vpcID),
				CidrBlock: aws.String("192.168.0.0/16"),
			}},
		}, nil)

		for t := range subnetLists {
			fn := func(list string, subnetsByAz map[string]api.AZSubnetSpec) {
				subnets := strings.Split(list, ",")

				output := &ec2.DescribeSubnetsOutput{
					Subnets: make([]*ec2.Subnet, len(subnets)),
				}

				for i := range subnets {
					subnet := &ec2.Subnet{}
					subnet.SetSubnetId(subnets[i])
					subnet.SetVpcId(vpcID)
					for az := range subnetsByAz {
						if subnetsByAz[az].ID == subnets[i] {
							subnet.SetAvailabilityZone(az)
							subnet.SetCidrBlock(subnetsByAz[az].CIDR.String())
						}
					}
					output.Subnets[i] = subnet

				}
				p.MockEC2().On("DescribeSubnets", mock.MatchedBy(func(input *ec2.DescribeSubnetsInput) bool {
					_, _ = fmt.Fprintf(GinkgoWriter, "%s subnets = %#v\n", t, output)
					return joinCompare(input, list)
				})).Return(output, nil)
			}
			switch t {
			case "Private":
				fn(subnetLists[t], testVPC.Subnets.Private)
			case "Public":
				fn(subnetLists[t], testVPC.Subnets.Public)
			}

		}
	}

	Describe("GetAllOutputsFromClusterStack", func() {
		expected := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Region:  "us-west-2",
				Name:    clusterName,
				Version: "1.19",
			},
			Status: &api.ClusterStatus{
				Endpoint:                 endpoint,
				CertificateAuthorityData: caCertData,
				ARN:                      arn,
			},
			AvailabilityZones: testAZs,
			VPC:               testVPC(),
			IAM: &api.ClusterIAM{
				WithOIDC:       api.Disabled(),
				ServiceRoleARN: aws.String(arn),
			},
			CloudWatch: &api.ClusterCloudWatch{
				ClusterLogging: &api.ClusterCloudWatchLogging{},
			},
			PrivateCluster: &api.PrivateCluster{},
			NodeGroups: []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						AMIFamily:         "AmazonLinux2",
						InstanceType:      "t2.medium",
						Name:              "ng-abcd1234",
						PrivateNetworking: false,
						VolumeSize:        aws.Int(2),
						IAM: &api.NodeGroupIAM{
							WithAddonPolicies: api.NodeGroupIAMAddonPolicies{
								ImageBuilder:              api.Disabled(),
								AutoScaler:                api.Disabled(),
								ExternalDNS:               api.Disabled(),
								CertManager:               api.Disabled(),
								AppMesh:                   api.Disabled(),
								AppMeshPreview:            api.Disabled(),
								EBS:                       api.Disabled(),
								FSX:                       api.Disabled(),
								EFS:                       api.Disabled(),
								AWSLoadBalancerController: api.Disabled(),
								XRay:                      api.Disabled(),
								CloudWatch:                api.Disabled(),
							},
						},
						ScalingConfig: &api.ScalingConfig{},
						SSH: &api.NodeGroupSSH{
							Allow:         api.Disabled(),
							PublicKeyPath: &api.DefaultNodeSSHPublicKeyPath,
						},
						AMI: "",
						SecurityGroups: &api.NodeGroupSGs{
							WithLocal:  api.Enabled(),
							WithShared: api.Enabled(),
							AttachIDs:  []string{},
						},
						VolumeType:       aws.String(api.NodeVolumeTypeSC1),
						VolumeName:       aws.String("/dev/xvda"),
						VolumeEncrypted:  api.Disabled(),
						VolumeKmsKeyID:   aws.String(""),
						DisableIMDSv1:    api.Disabled(),
						DisablePodIMDS:   api.Disabled(),
						InstanceSelector: &api.InstanceSelector{},
					},
				},
			},
		}

		cfg := newSimpleClusterConfig()

		setSubnets(cfg)

		sampleOutputs := map[string]string{
			"SecurityGroup":              "sg-0b44c48bcba5b7362",
			"SubnetsPublic":              subnetsPublic,
			"SubnetsPrivate":             subnetsPrivate,
			"VPC":                        vpcID,
			"Endpoint":                   endpoint,
			"CertificateAuthorityData":   caCert,
			"ARN":                        arn,
			"ClusterStackName":           "",
			"SharedNodeSecurityGroup":    "sg-shared",
			"ServiceRoleARN":             arn,
			"FargatePodExecutionRoleARN": fargatePodExecutionRoleARN,
			"FeatureNATMode":             "Single",
			"ClusterSecurityGroupId":     "sg-09ef4509a37f28b4c",
		}

		It("should add all resources and collect outputs without errors", func() {
			crs = NewClusterResourceSet(p.EC2(), p.Region(), cfg, false, nil)
			err := crs.AddAllResources()
			Expect(err).ShouldNot(HaveOccurred())
			sampleStack := newStackWithOutputs(sampleOutputs)
			err = crs.GetAllOutputs(sampleStack)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("resulting config should match what is expected", func() {
			cfgData, err := json.Marshal(cfg)
			Expect(err).ShouldNot(HaveOccurred())
			expectedData, err := json.Marshal(expected)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfgData).To(MatchJSON(expectedData))
		})
	})

	assertBuildChecks := func(cfg *api.ClusterConfig, clusterStackName string, ng *api.NodeGroup, managedNodesSupport bool) {
		It("should add all resources without errors", func() {
			crs = NewClusterResourceSet(p.EC2(), p.Region(), cfg, managedNodesSupport, nil)
			err = crs.AddAllResources()
			Expect(err).ShouldNot(HaveOccurred())

			// TODO: yes I know this is terrible, I will improve this whole disaster
			// of a test file in another PR
			fakeVPCImporter := new(vpcfakes.FakeImporter)
			fakeVPCImporter.ControlPlaneSecurityGroupReturns(gfnt.MakeFnImportValueString(clusterStackName + "::SecurityGroup"))
			fakeVPCImporter.SharedNodeSecurityGroupReturns(gfnt.MakeFnImportValueString(clusterStackName + "::SecurityGroup"))
			fakeVPCImporter.SubnetsPrivateReturns(gfnt.MakeFnSplit(",", gfnt.MakeFnImportValueString(clusterStackName+"::SubnetsPrivate")))
			fakeVPCImporter.SubnetsPublicReturns(gfnt.MakeFnSplit(",", gfnt.MakeFnImportValueString(clusterStackName+"::SubnetsPublic")))

			// TODO see note above. worst testfile ever
			fakeBootstrapper := new(bootstrapfakes.FakeBootstrapper)
			fakeBootstrapper.UserDataReturns("lovely data right here", nil)

			ngrs = NewNodeGroupResourceSet(p.EC2(), p.IAM(), cfg, ng, managedNodesSupport, false, fakeVPCImporter)
			ngrs.SetBootstrapper(fakeBootstrapper)

			err = ngrs.AddAllResources()
			Expect(err).ShouldNot(HaveOccurred())

			t := ngrs.Template()
			Expect(t.Resources).Should(HaveKey("NodeGroup"))

			templateBody, err := t.JSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(templateBody).ShouldNot(BeEmpty())

			outputs := crs.Template().Outputs
			_, hasClusterSG := outputs["ClusterSecurityGroupId"]

			Expect(hasClusterSG).To(Equal(managedNodesSupport))
		})
	}

	build := func(cfg *api.ClusterConfig, name string, ng *api.NodeGroup) {
		assertBuildChecks(cfg, name, ng, false)
	}

	roundtrip := func() {
		It("should serialise JSON without errors, and parse the template", func() {
			ngTemplate = &Template{}
			{
				templateBody, err := ngrs.RenderJSON()
				Expect(err).ShouldNot(HaveOccurred())
				err = json.Unmarshal(templateBody, ngTemplate)
				Expect(err).ShouldNot(HaveOccurred())
			}
			clusterTemplate = &Template{}
			{
				templateBody, err := crs.RenderJSON()
				Expect(err).ShouldNot(HaveOccurred())
				err = json.Unmarshal(templateBody, clusterTemplate)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	}

	Context("Security group for managed nodes", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)
		assertBuildChecks(cfg, "managed-cluster", ng, true)
	})

	Context("AutoNameTag", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtrip()

		It("SG should have correct tags", func() {
			Expect(ngTemplate.Resources).ToNot(BeNil())
			Expect(ngTemplate.Resources).To(HaveLen(8))
			Expect(ngTemplate.Resources["SG"].Properties.Tags).To(HaveLen(2))
			Expect(ngTemplate.Resources["SG"].Properties.Tags[0].Key).To(Equal("kubernetes.io/cluster/" + clusterName))
			Expect(ngTemplate.Resources["SG"].Properties.Tags[0].Value).To(Equal("owned"))
			Expect(ngTemplate.Resources["SG"].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(ngTemplate.Resources["SG"].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SG",
			}))
		})
	})

	Context("NodeGroupTags", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "t2.medium"
		ng.Name = "ng-abcd1234"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtrip()

		It("should have correct tags", func() {
			expectedTags := []Tag{
				{
					Key:               "Name",
					Value:             clusterName + "-ng-abcd1234-Node",
					PropagateAtLaunch: "true",
				},
				{
					Key:               "kubernetes.io/cluster/" + clusterName,
					Value:             "owned",
					PropagateAtLaunch: "true",
				},
			}

			ngProps := getNodeGroupProperties(ngTemplate)

			Expect(ngProps.Tags).ToNot(BeNil())
			Expect(ngProps.Tags).To(Equal(expectedTags))
		})
	})

	Context("NodeGroup DesiredCapacity=nil MaxSize=nil MinSize=nil", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = nil
		ng.MaxSize = nil
		ng.MinSize = nil

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test1-cluster", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(ngTemplate).MaxSize).To(Equal("2"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("2"))

		})
	})

	Context("NodeGroup DesiredCapacity=10 MaxSize=nil MinSize=nil", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = nil
		ng.MaxSize = nil
		ng.MinSize = nil

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test2-cluster", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(ngTemplate).MaxSize).To(Equal("2"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("2"))

		})
	})

	Context("NodeGroup DesiredCapacity=nil MaxSize=30 MinSize=nil", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = nil
		ng.MaxSize = new(int)
		*ng.MaxSize = 30
		ng.MinSize = nil

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test3-cluster", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(ngTemplate).MaxSize).To(Equal("30"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("2"))
		})
	})

	Context("NodeGroup DesiredCapacity=nil MaxSize=nil MinSize=90", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = nil
		ng.MaxSize = nil
		ng.MinSize = new(int)
		*ng.MinSize = 90

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test4-cluster", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(ngTemplate).MaxSize).To(Equal("90"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("90"))
		})
	})

	Context("NodeGroup DesiredCapacity=nil MaxSize=91 MinSize=61", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = nil
		ng.MaxSize = new(int)
		*ng.MaxSize = 91
		ng.MinSize = new(int)
		*ng.MinSize = 61

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test5-cluster", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(ngTemplate).MaxSize).To(Equal("91"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("61"))
		})
	})

	Context("NodeGroup DesiredCapacity=32 MaxSize=92 MinSize=nil", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = new(int)
		*ng.DesiredCapacity = 32
		ng.MaxSize = new(int)
		*ng.MaxSize = 92
		ng.MinSize = nil

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test6-cluster", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).DesiredCapacity).To(Equal("32"))
			Expect(getNodeGroupProperties(ngTemplate).MaxSize).To(Equal("92"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("32"))
		})
	})

	Context("NodeGroup DesiredCapacity=33 MaxSize=nil MinSize=31", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = new(int)
		*ng.DesiredCapacity = 33
		ng.MaxSize = nil
		ng.MinSize = new(int)
		*ng.MinSize = 31

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test7-cluster", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).DesiredCapacity).To(Equal("33"))
			Expect(getNodeGroupProperties(ngTemplate).MaxSize).To(Equal("33"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("31"))
		})
	})

	Context("NodeGroupAutoScaling", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)
		ng.ClassicLoadBalancerNames = []string{"clb-1", "clb-2"}
		ng.TargetGroupARNs = []string{"tg-arn-1", "tg-arn-2"}

		ng.MinSize = new(int)
		*ng.MinSize = 10
		ng.InstanceType = "m5.2xlarge"

		ng.IAM.InstanceRoleName = "a-named-role"
		ng.IAM.WithAddonPolicies.AutoScaler = api.Enabled()

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtrip()

		It("should have correct instance type and min size", func() {
			Expect(getLaunchTemplateData(ngTemplate).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(ngTemplate).MinSize).To(Equal("10"))
		})

		It("should have correct instance role and profile", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeInstanceRole"))

			role := ngTemplate.Resources["NodeInstanceRole"].Properties

			Expect(role.Path).To(Equal("/"))
			Expect(role.RoleName).To(Equal("a-named-role"))
			Expect(role.ManagedPolicyArns).To(ConsistOf(makePolicyARNRef("AmazonEKSWorkerNodePolicy",
				"AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly")...))

			checkARPD([]string{"EC2"}, role.AssumeRolePolicyDocument)

			Expect(ngTemplate.Resources).To(HaveKey("NodeInstanceProfile"))

			profile := ngTemplate.Resources["NodeInstanceProfile"].Properties

			Expect(profile.Path).To(Equal("/"))
			Expect(profile.Roles).To(HaveLen(1))
			isRefTo(profile.Roles[0], "NodeInstanceRole")

			isFnGetAttOf(getLaunchTemplateData(ngTemplate).IamInstanceProfile.Arn, "NodeInstanceProfile", "Arn")
		})

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())
			Expect(ngTemplate.Resources).To(HaveKey("PolicyAutoScaling"))

			policy := ngTemplate.Resources["PolicyAutoScaling"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"autoscaling:DescribeAutoScalingGroups",
				"autoscaling:DescribeAutoScalingInstances",
				"autoscaling:DescribeLaunchConfigurations",
				"autoscaling:DescribeTags",
				"autoscaling:SetDesiredCapacity",
				"autoscaling:TerminateInstanceInAutoScalingGroup",
				"ec2:DescribeLaunchTemplateVersions",
			}))
		})

		It("should have auto-discovery tags", func() {
			expectedTags := []Tag{
				{
					Key:               "Name",
					Value:             fmt.Sprintf("%s-%s-Node", clusterName, "ng-abcd1234"),
					PropagateAtLaunch: "true",
				},
				{
					Key:               "kubernetes.io/cluster/" + clusterName,
					Value:             "owned",
					PropagateAtLaunch: "true",
				},
				{
					Key:               "k8s.io/cluster-autoscaler/enabled",
					Value:             "true",
					PropagateAtLaunch: "true",
				},
				{
					Key:               "k8s.io/cluster-autoscaler/" + clusterName,
					Value:             "owned",
					PropagateAtLaunch: "true",
				},
			}

			ngProps := getNodeGroupProperties(ngTemplate)

			Expect(ngProps.Tags).ToNot(BeNil())
			Expect(ngProps.Tags).To(Equal(expectedTags))
		})

		It("should have classic load balancer names set", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
			ng := ngTemplate.Resources["NodeGroup"]
			Expect(ng).ToNot(BeNil())
			Expect(ng.Properties).ToNot(BeNil())

			Expect(ng.Properties.LoadBalancerNames).To(Equal([]string{"clb-1", "clb-2"}))
		})

		It("should have target groups ARNs set", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
			ng := ngTemplate.Resources["NodeGroup"]
			Expect(ng).ToNot(BeNil())
			Expect(ng.Properties).ToNot(BeNil())

			Expect(ng.Properties.TargetGroupARNs).To(Equal([]string{"tg-arn-1", "tg-arn-2"}))
		})

		It("should have target groups ARNs set", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
			ng := ngTemplate.Resources["NodeGroup"]
			Expect(ng).ToNot(BeNil())
			Expect(ng.Properties).ToNot(BeNil())

			Expect(ng.Properties.TargetGroupARNs).To(Equal([]string{"tg-arn-1", "tg-arn-2"}))
		})
	})

	Context("NodeGroupAutoScaling with metrics collection", func() {
		Context("with all details", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)
			ng.ASGMetricsCollection = []api.MetricsCollection{
				{
					Granularity: "1Minute",
					Metrics: []string{
						"GroupMinSize",
						"GroupMaxSize",
					},
				},
			}
			build(cfg, "eksctl-test-123-with-metrics", ng)
			roundtrip()

			It("should have both Granularity and Metrics details", func() {
				Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
				ng := ngTemplate.Resources["NodeGroup"]
				Expect(ng).ToNot(BeNil())
				Expect(ng.Properties).ToNot(BeNil())

				Expect(ng.Properties.MetricsCollection).To(HaveLen(1))
				var metricsCollection = ng.Properties.MetricsCollection[0]
				Expect(metricsCollection).To(HaveKey("Granularity"))
				Expect(metricsCollection).To(HaveKey("Metrics"))
				Expect(metricsCollection["Granularity"]).To(Equal("1Minute"))
				Expect(metricsCollection["Metrics"]).To(ContainElement("GroupMinSize"))
				Expect(metricsCollection["Metrics"]).To(ContainElement("GroupMaxSize"))
			})
		})

		Context("without metrics details", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)
			ng.ASGMetricsCollection = []api.MetricsCollection{
				{
					Granularity: "1Minute",
				},
			}
			build(cfg, "eksctl-test-123-cluster", ng)
			roundtrip()

			It("should have only Granularity", func() {
				Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
				ng := ngTemplate.Resources["NodeGroup"]
				Expect(ng).ToNot(BeNil())
				Expect(ng.Properties).ToNot(BeNil())

				Expect(ng.Properties.MetricsCollection).To(HaveLen(1))
				var metricsCollection = ng.Properties.MetricsCollection[0]
				Expect(metricsCollection).To(HaveKey("Granularity"))
				Expect(metricsCollection["Granularity"]).To(Equal("1Minute"))
			})
		})
	})

	Context("NodeGroupCertManagerExternalDNS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.CertManager = api.Enabled()
		ng.IAM.WithAddonPolicies.ExternalDNS = api.Enabled()

		build(cfg, "eksctl-test-cert-manager-external-dns-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources).To(HaveKey("PolicyCertManagerChangeSet"))

			policy1 := ngTemplate.Resources["PolicyCertManagerChangeSet"].Properties

			Expect(policy1.Roles).To(HaveLen(1))
			isRefTo(policy1.Roles[0], "NodeInstanceRole")

			Expect(policy1.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy1.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy1.PolicyDocument.Statement[0].Resource).To(Equal(map[string]interface{}{
				"Fn::Sub": "arn:${AWS::Partition}:route53:::hostedzone/*",
			}))
			Expect(policy1.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ChangeResourceRecordSets",
			}))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyCertManagerHostedZones"))

			policy2 := ngTemplate.Resources["PolicyCertManagerHostedZones"].Properties

			Expect(policy2.Roles).To(HaveLen(1))
			isRefTo(policy2.Roles[0], "NodeInstanceRole")

			Expect(policy2.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy2.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy2.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy2.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ListResourceRecordSets",
				"route53:ListHostedZonesByName",
			}))

			policy3 := ngTemplate.Resources["PolicyExternalDNSHostedZones"].Properties

			Expect(policy3.Roles).To(HaveLen(1))
			isRefTo(policy3.Roles[0], "NodeInstanceRole")
			Expect(policy3.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy3.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy3.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy3.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
				"route53:ListTagsForResource",
			}))

			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAppMesh"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAppMeshPreview"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyEBS"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyFSX"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyServiceLinkRole"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyEFS"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyEFSEC2"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAWSLoadBalancerController"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyXRay"))
		})

	})

	Context("NodeGroupAppMeshExternalDNS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.AppMesh = api.Enabled()
		ng.IAM.WithAddonPolicies.ExternalDNS = api.Enabled()

		build(cfg, "eksctl-test-megaapps-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources).To(HaveKey("PolicyExternalDNSChangeSet"))

			policy1 := ngTemplate.Resources["PolicyExternalDNSChangeSet"].Properties

			Expect(policy1.Roles).To(HaveLen(1))
			isRefTo(policy1.Roles[0], "NodeInstanceRole")

			Expect(policy1.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy1.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy1.PolicyDocument.Statement[0].Resource).To(Equal(map[string]interface{}{
				"Fn::Sub": "arn:${AWS::Partition}:route53:::hostedzone/*",
			}))
			Expect(policy1.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ChangeResourceRecordSets",
			}))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyExternalDNSHostedZones"))

			policy2 := ngTemplate.Resources["PolicyExternalDNSHostedZones"].Properties

			Expect(policy2.Roles).To(HaveLen(1))
			isRefTo(policy2.Roles[0], "NodeInstanceRole")

			Expect(policy2.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy2.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy2.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy2.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
				"route53:ListTagsForResource",
			}))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyAppMesh"))

			policy3 := ngTemplate.Resources["PolicyAppMesh"].Properties

			Expect(policy3.Roles).To(HaveLen(1))
			isRefTo(policy3.Roles[0], "NodeInstanceRole")

			Expect(policy3.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy3.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy3.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy3.PolicyDocument.Statement[0].Action).To(Equal(append(appMeshActions, "appmesh:*")))

			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyEBS"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAutoScaling"))
		})

	})

	Context("NodeGroupAppMeshPreview", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.AppMeshPreview = api.Enabled()

		build(cfg, "eksctl-test-appmesh-preview", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).To(HaveKey("PolicyAppMeshPreview"))

			policy3 := ngTemplate.Resources["PolicyAppMeshPreview"].Properties

			Expect(policy3.Roles).To(HaveLen(1))
			isRefTo(policy3.Roles[0], "NodeInstanceRole")

			Expect(policy3.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy3.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy3.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy3.PolicyDocument.Statement[0].Action).To(Equal(append(appMeshActions, "appmesh-preview:*")))
		})
	})

	Context("NodeGroupAppCertManager", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.CertManager = api.Enabled()

		build(cfg, "eksctl-test-cert-manager-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources).To(HaveKey("PolicyCertManagerChangeSet"))

			policy1 := ngTemplate.Resources["PolicyCertManagerChangeSet"].Properties

			Expect(policy1.Roles).To(HaveLen(1))
			isRefTo(policy1.Roles[0], "NodeInstanceRole")

			Expect(policy1.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy1.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy1.PolicyDocument.Statement[0].Resource).To(Equal(map[string]interface{}{
				"Fn::Sub": "arn:${AWS::Partition}:route53:::hostedzone/*",
			}))
			Expect(policy1.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ChangeResourceRecordSets",
			}))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyCertManagerHostedZones"))

			policy2 := ngTemplate.Resources["PolicyCertManagerHostedZones"].Properties

			Expect(policy2.Roles).To(HaveLen(1))
			isRefTo(policy2.Roles[0], "NodeInstanceRole")

			Expect(policy2.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy2.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy2.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy2.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ListResourceRecordSets",
				"route53:ListHostedZonesByName",
			}))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyCertManagerGetChange"))

			policy3 := ngTemplate.Resources["PolicyCertManagerGetChange"].Properties

			Expect(policy3.Roles).To(HaveLen(1))
			isRefTo(policy3.Roles[0], "NodeInstanceRole")

			Expect(policy3.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy3.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy3.PolicyDocument.Statement[0].Resource).To(Equal(map[string]interface{}{
				"Fn::Sub": "arn:${AWS::Partition}:route53:::change/*",
			}))
			Expect(policy3.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:GetChange",
			}))
		})

	})

	Context("NodeGroupAppExternalDNS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.ExternalDNS = api.Enabled()

		build(cfg, "eksctl-test-external-dns-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources).To(HaveKey("PolicyExternalDNSChangeSet"))

			policy1 := ngTemplate.Resources["PolicyExternalDNSChangeSet"].Properties

			Expect(policy1.Roles).To(HaveLen(1))
			isRefTo(policy1.Roles[0], "NodeInstanceRole")

			Expect(policy1.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy1.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy1.PolicyDocument.Statement[0].Resource).To(Equal(map[string]interface{}{
				"Fn::Sub": "arn:${AWS::Partition}:route53:::hostedzone/*",
			}))
			Expect(policy1.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ChangeResourceRecordSets",
			}))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyExternalDNSHostedZones"))
			policy2 := ngTemplate.Resources["PolicyExternalDNSHostedZones"].Properties

			Expect(policy2.Roles).To(HaveLen(1))
			isRefTo(policy2.Roles[0], "NodeInstanceRole")

			Expect(policy2.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy2.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy2.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy2.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
				"route53:ListTagsForResource",
			}))

			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyCertManagerGetChange"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyCertManagerHostedZones"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAppMesh"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyEBS"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyFSX"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyServiceLinkRole"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyEFS"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyEFSEC2"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAWSLoadBalancerController"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyXRay"))
		})
	})

	Context("NodeGroupAWSLoadBalancerController", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.AWSLoadBalancerController = api.Enabled()

		build(cfg, "eksctl-test-megaapps-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources).To(HaveKey("PolicyAWSLoadBalancerController"))

			policy := ngTemplate.Resources["PolicyAWSLoadBalancerController"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(7))

			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal(
				map[string]interface{}{
					"Fn::Sub": "arn:${AWS::Partition}:ec2:*:*:security-group/*",
				},
			))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{"ec2:CreateTags"}))
			Expect(policy.PolicyDocument.Statement[0].Condition).To(HaveLen(2))

			Expect(policy.PolicyDocument.Statement[1].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[1].Resource).To(Equal(
				map[string]interface{}{
					"Fn::Sub": "arn:${AWS::Partition}:ec2:*:*:security-group/*",
				},
			))
			Expect(policy.PolicyDocument.Statement[1].Action).To(Equal([]string{"ec2:CreateTags", "ec2:DeleteTags"}))
			Expect(policy.PolicyDocument.Statement[1].Condition).To(HaveLen(1))

			Expect(policy.PolicyDocument.Statement[2].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[2].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[2].Action).To(Equal([]string{
				"elasticloadbalancing:CreateLoadBalancer",
				"elasticloadbalancing:CreateTargetGroup",
			}))
			Expect(policy.PolicyDocument.Statement[2].Condition).To(HaveLen(1))

			Expect(policy.PolicyDocument.Statement[3].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[3].Resource).To(Equal([]interface{}{
				map[string]interface{}{
					"Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:targetgroup/*/*",
				},
				map[string]interface{}{
					"Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:loadbalancer/net/*/*",
				},
				map[string]interface{}{
					"Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:loadbalancer/app/*/*",
				},
			}))
			Expect(policy.PolicyDocument.Statement[3].Action).To(Equal([]string{
				"elasticloadbalancing:AddTags",
				"elasticloadbalancing:RemoveTags",
			}))
			Expect(policy.PolicyDocument.Statement[3].Condition).To(HaveLen(1))

			Expect(policy.PolicyDocument.Statement[4].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[4].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[4].Action).To(Equal([]string{
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:RevokeSecurityGroupIngress",
				"ec2:DeleteSecurityGroup",
				"elasticloadbalancing:ModifyLoadBalancerAttributes",
				"elasticloadbalancing:SetIpAddressType",
				"elasticloadbalancing:SetSecurityGroups",
				"elasticloadbalancing:SetSubnets",
				"elasticloadbalancing:DeleteLoadBalancer",
				"elasticloadbalancing:ModifyTargetGroup",
				"elasticloadbalancing:ModifyTargetGroupAttributes",
				"elasticloadbalancing:DeleteTargetGroup",
			}))
			Expect(policy.PolicyDocument.Statement[4].Condition).To(HaveLen(1))

			Expect(policy.PolicyDocument.Statement[5].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[5].Resource).To(Equal(
				map[string]interface{}{
					"Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:targetgroup/*/*",
				},
			))
			Expect(policy.PolicyDocument.Statement[5].Action).To(Equal([]string{
				"elasticloadbalancing:RegisterTargets",
				"elasticloadbalancing:DeregisterTargets",
			}))
			Expect(policy.PolicyDocument.Statement[5].Condition).To(HaveLen(0))

			Expect(policy.PolicyDocument.Statement[6].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[6].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[6].Action).To(Equal([]string{
				"iam:CreateServiceLinkedRole",
				"ec2:DescribeAccountAttributes",
				"ec2:DescribeAddresses",
				"ec2:DescribeInternetGateways",
				"ec2:DescribeVpcs",
				"ec2:DescribeSubnets",
				"ec2:DescribeSecurityGroups",
				"ec2:DescribeInstances",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DescribeTags",
				"elasticloadbalancing:DescribeLoadBalancers",
				"elasticloadbalancing:DescribeLoadBalancerAttributes",
				"elasticloadbalancing:DescribeListeners",
				"elasticloadbalancing:DescribeListenerCertificates",
				"elasticloadbalancing:DescribeSSLPolicies",
				"elasticloadbalancing:DescribeRules",
				"elasticloadbalancing:DescribeTargetGroups",
				"elasticloadbalancing:DescribeTargetGroupAttributes",
				"elasticloadbalancing:DescribeTargetHealth",
				"elasticloadbalancing:DescribeTags",
				"cognito-idp:DescribeUserPoolClient",
				"acm:ListCertificates",
				"acm:DescribeCertificate",
				"iam:ListServerCertificates",
				"iam:GetServerCertificate",
				"waf-regional:GetWebACL",
				"waf-regional:GetWebACLForResource",
				"waf-regional:AssociateWebACL",
				"waf-regional:DisassociateWebACL",
				"wafv2:GetWebACL",
				"wafv2:GetWebACLForResource",
				"wafv2:AssociateWebACL",
				"wafv2:DisassociateWebACL",
				"shield:GetSubscriptionState",
				"shield:DescribeProtection",
				"shield:CreateProtection",
				"shield:DeleteProtection",
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:RevokeSecurityGroupIngress",
				"ec2:CreateSecurityGroup",
				"elasticloadbalancing:CreateListener",
				"elasticloadbalancing:DeleteListener",
				"elasticloadbalancing:CreateRule",
				"elasticloadbalancing:DeleteRule",
				"elasticloadbalancing:SetWebAcl",
				"elasticloadbalancing:ModifyListener",
				"elasticloadbalancing:AddListenerCertificates",
				"elasticloadbalancing:RemoveListenerCertificates",
				"elasticloadbalancing:ModifyRule",
			}))
			Expect(policy.PolicyDocument.Statement[6].Condition).To(HaveLen(0))
		})
	})

	Context("NodeGroupXRay", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.XRay = api.Enabled()

		build(cfg, "eksctl-test-megaapps-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources).To(HaveKey("PolicyXRay"))

			policy := ngTemplate.Resources["PolicyXRay"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"xray:PutTraceSegments",
				"xray:PutTelemetryRecords",
				"xray:GetSamplingRules",
				"xray:GetSamplingTargets",
				"xray:GetSamplingStatisticSummaries",
			}))
		})

	})

	Context("NodeGroupCloudWatch", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.CloudWatch = api.Enabled()

		build(cfg, "eksctl-test-cwenabled-cluster", ng)

		roundtrip()

		It("should have correct managed profile", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeInstanceRole"))

			role := ngTemplate.Resources["NodeInstanceRole"].Properties

			Expect(role.ManagedPolicyArns).To(ConsistOf(makePolicyARNRef("AmazonEKSWorkerNodePolicy",
				"AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "CloudWatchAgentServerPolicy")))
		})
	})

	Context("NodeGroupEBS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeSize = nil
		ng.IAM.WithAddonPolicies.EBS = api.Enabled()

		build(cfg, "eksctl-test-ebs-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(getLaunchTemplateData(ngTemplate).BlockDeviceMappings).To(HaveLen(0))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyEBS"))

			policy := ngTemplate.Resources["PolicyEBS"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"ec2:AttachVolume",
				"ec2:CreateSnapshot",
				"ec2:CreateTags",
				"ec2:CreateVolume",
				"ec2:DeleteSnapshot",
				"ec2:DeleteTags",
				"ec2:DeleteVolume",
				"ec2:DescribeAvailabilityZones",
				"ec2:DescribeInstances",
				"ec2:DescribeSnapshots",
				"ec2:DescribeTags",
				"ec2:DescribeVolumes",
				"ec2:DescribeVolumesModifications",
				"ec2:DetachVolume",
				"ec2:ModifyVolume",
			}))

			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyExternalDNSChangeSet"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyExternalDNSHostedZones"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAppMesh"))
		})
	})

	Context("NodeGroupFSX", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeSize = nil
		ng.IAM.WithAddonPolicies.FSX = api.Enabled()

		build(cfg, "eksctl-test-fsx-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(getLaunchTemplateData(ngTemplate).BlockDeviceMappings).To(HaveLen(0))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyFSX"))

			policy := ngTemplate.Resources["PolicyFSX"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"fsx:*",
			}))

			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyExternalDNSChangeSet"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyExternalDNSHostedZones"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAppMesh"))
		})
	})

	Context("NodeGroupEFS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeSize = nil
		ng.IAM.WithAddonPolicies.EFS = api.Enabled()

		build(cfg, "eksctl-test-efs-cluster", ng)

		roundtrip()

		It("should have correct policies", func() {
			Expect(getLaunchTemplateData(ngTemplate).BlockDeviceMappings).To(HaveLen(0))

			Expect(ngTemplate.Resources).To(HaveKey("PolicyEFS"))

			policy := ngTemplate.Resources["PolicyEFS"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"elasticfilesystem:*",
			}))

			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyExternalDNSChangeSet"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyExternalDNSHostedZones"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("PolicyAppMesh"))
		})
	})

	Context("NodeGroup with custom role and profile", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.InstanceRoleARN = "arn:role"
		ng.IAM.InstanceProfileARN = "arn:profile"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtrip()

		It("should have correct instance role and profile", func() {
			Expect(ngTemplate.Resources).ToNot(HaveKey("NodeInstanceRole"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("NodeInstanceProfile"))

			Expect(getLaunchTemplateData(ngTemplate).IamInstanceProfile.Arn).To(Equal("arn:profile"))
		})
	})

	Context("NodeGroup with cutom role", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.InstanceRoleARN = "arn:role"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtrip()

		It("should have correct instance role and profile", func() {
			Expect(ngTemplate.Resources).ToNot(HaveKey("NodeInstanceRole"))

			Expect(ngTemplate.Resources).To(HaveKey("NodeInstanceProfile"))

			profile := ngTemplate.Resources["NodeInstanceProfile"].Properties

			Expect(profile.Path).To(Equal("/"))
			Expect(profile.Roles).To(HaveLen(1))
			Expect(profile.Roles[0]).To(Equal("arn:role"))

			isFnGetAttOf(getLaunchTemplateData(ngTemplate).IamInstanceProfile.Arn, "NodeInstanceProfile", "Arn")
		})
	})

	Context("NodeGroup with custom role containing a deep resource path is normalized", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.InstanceRoleARN = "arn:aws:iam::1234567890:role/foo/bar/baz/custom-eks-role"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtrip()

		It("should have correct instance role and profile", func() {
			Expect(ngTemplate.Resources).ToNot(HaveKey("NodeInstanceRole"))
			Expect(ngTemplate.Resources).To(HaveKey("NodeInstanceProfile"))

			profile := ngTemplate.Resources["NodeInstanceProfile"].Properties

			Expect(profile.Path).To(Equal("/"))
			Expect(profile.Roles).To(HaveLen(1))
			Expect(profile.Roles[0]).To(Equal("custom-eks-role"))

			isFnGetAttOf(getLaunchTemplateData(ngTemplate).IamInstanceProfile.Arn, "NodeInstanceProfile", "Arn")
		})
	})

	Context("NodeGroup with cutom profile", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.InstanceProfileARN = "arn:profile"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtrip()

		It("should have correct instance role and profile", func() {
			Expect(ngTemplate.Resources).ToNot(HaveKey("NodeInstanceRole"))
			Expect(ngTemplate.Resources).ToNot(HaveKey("NodeInstanceProfile"))

			Expect(getLaunchTemplateData(ngTemplate).IamInstanceProfile.Arn).To(Equal("arn:profile"))
		})
	})

	Context("Nodegroup encrypted volume using default key, or encrypted AMI", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeEncrypted = api.Enabled()

		build(cfg, "eksctl-test-private-ng", ng)

		roundtrip()

		It("should have correct resources and attributes", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			ltd := getLaunchTemplateData(ngTemplate)
			Expect(ltd.BlockDeviceMappings).To(HaveLen(1))

			rootVolume := ltd.BlockDeviceMappings[0].(map[string]interface{})
			Expect(rootVolume).To(HaveKey("Ebs"))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("Encrypted", true))
		})
	})

	Context("Nodegroup encrypted volume using CMK", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeEncrypted = api.Enabled()
		*ng.VolumeKmsKeyID = "36c0b54e-64ed-4f2d-a1c7-96558764311e"

		build(cfg, "eksctl-test-private-ng", ng)

		roundtrip()

		It("should have correct resources and attributes", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			ltd := getLaunchTemplateData(ngTemplate)
			Expect(ltd.BlockDeviceMappings).To(HaveLen(1))

			rootVolume := ltd.BlockDeviceMappings[0].(map[string]interface{})
			Expect(rootVolume).To(HaveKey("Ebs"))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("Encrypted", true))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("KmsKeyId", "36c0b54e-64ed-4f2d-a1c7-96558764311e"))
		})
	})

	Context("Nodegroup{VolumeType=sc1 VolumeSize=2.0}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		build(cfg, "eksctl-test-private-ng", ng)

		roundtrip()

		It("should have correct resources and attributes", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			ltd := getLaunchTemplateData(ngTemplate)
			Expect(ltd.BlockDeviceMappings).To(HaveLen(1))

			rootVolume := ltd.BlockDeviceMappings[0].(map[string]interface{})
			Expect(rootVolume).To(HaveKey("Ebs"))
			Expect(rootVolume).To(HaveKeyWithValue("DeviceName", "/dev/xvda"))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("VolumeType", "sc1"))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("VolumeSize", 2.0))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("Encrypted", false))
		})
	})

	assertSSHRules := func(expectedIngressRules string) {
		bytes, err := ngrs.RenderJSON()
		Expect(err).ToNot(HaveOccurred())
		template, err := goformation.ParseJSON(bytes)
		Expect(err).ToNot(HaveOccurred())

		securityGroup, err := template.GetEC2SecurityGroupWithName("SG")
		Expect(err).ToNot(HaveOccurred())
		ingressRules, err := json.Marshal(securityGroup.SecurityGroupIngress)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(ingressRules)).To(MatchJSON(expectedIngressRules))
	}

	Context("NodeGroup{PrivateNetworking=true SSH.Allow=true}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.SSH.Allow = api.Enabled()
		keyName := ""
		ng.SSH.PublicKeyName = &keyName
		ng.InstanceType = "t2.medium"
		ng.PrivateNetworking = true
		ng.AMIFamily = "AmazonLinux2"

		build(cfg, "eksctl-test-private-ng", ng)

		roundtrip()

		It("should have correct description", func() {
			Expect(ngTemplate.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(ngTemplate.Description).To(ContainSubstring("SSH access: true"))
			Expect(ngTemplate.Description).To(ContainSubstring("private networking: true"))
		})

		It("should have correct resources and attributes", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
			ng := ngTemplate.Resources["NodeGroup"].Properties
			Expect(ng.VPCZoneIdentifier).ToNot(BeNil())
			x, ok := ng.VPCZoneIdentifier.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(x).To(HaveLen(1))
			refSubnets := map[string]interface{}{
				"Fn::Split": []interface{}{
					",",
					map[string]interface{}{
						"Fn::ImportValue": "eksctl-test-private-ng::SubnetsPrivate",
					},
				},
			}
			Expect(x).To(Equal(refSubnets))

			ltd := getLaunchTemplateData(ngTemplate)

			isFnGetAttOf(ltd.IamInstanceProfile.Arn, "NodeInstanceProfile", "Arn")

			Expect(ltd.InstanceType).To(Equal("t2.medium"))

			Expect(ltd.NetworkInterfaces).To(HaveLen(1))
			Expect(ltd.NetworkInterfaces[0].DeviceIndex).To(Equal(0))
			Expect(ltd.NetworkInterfaces[0].AssociatePublicIPAddress).To(BeFalse())

			expectedIngressRules := `[
  {
    "Description": "[IngressInterCluster] Allow worker nodes in group ng-abcd1234 to communicate with control plane (kubelet and workload TCP ports)",
    "FromPort": 1025,
    "IpProtocol": "tcp",
    "SourceSecurityGroupId": {
      "Fn::ImportValue": "eksctl-test-private-ng::SecurityGroup"
    },
    "ToPort": 65535
  },
  {
    "Description": "[IngressInterClusterAPI] Allow worker nodes in group ng-abcd1234 to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)",
    "FromPort": 443,
    "IpProtocol": "tcp",
    "SourceSecurityGroupId": {
      "Fn::ImportValue": "eksctl-test-private-ng::SecurityGroup"
    },
    "ToPort": 443
  },
  {
    "CidrIp": "192.168.0.0/16",
    "Description": "Allow SSH access to worker nodes in group ng-abcd1234 (private, only inside VPC)",
    "FromPort": 22,
    "IpProtocol": "tcp",
    "ToPort": 22
  }
]`
			assertSSHRules(expectedIngressRules)
		})
	})

	Context("NodeGroup{PrivateNetworking=false SSH.Allow=true}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.SSH.Allow = api.Enabled()
		keyName := ""
		ng.SSH.PublicKeyName = &keyName
		ng.InstanceType = "t2.large"
		ng.PrivateNetworking = false
		ng.AMIFamily = "AmazonLinux2"

		build(cfg, "eksctl-test-public-ng", ng)

		roundtrip()

		It("should have correct description", func() {
			Expect(ngTemplate.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(ngTemplate.Description).To(ContainSubstring("SSH access: true"))
			Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have correct resources and attributes", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())
			Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))

			Expect(ngTemplate.Resources["NodeGroup"].Properties.VPCZoneIdentifier).ToNot(BeNil())
			x, ok := ngTemplate.Resources["NodeGroup"].Properties.VPCZoneIdentifier.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(x).To(HaveLen(1))
			refSubnets := map[string]interface{}{
				"Fn::Split": []interface{}{
					",",
					map[string]interface{}{
						"Fn::ImportValue": "eksctl-test-public-ng::SubnetsPublic",
					},
				},
			}
			Expect(x).To(Equal(refSubnets))

			ltd := getLaunchTemplateData(ngTemplate)

			Expect(ltd.InstanceType).To(Equal("t2.large"))

			Expect(ltd.NetworkInterfaces).To(HaveLen(1))
			Expect(ltd.NetworkInterfaces[0].DeviceIndex).To(Equal(0))
			Expect(ltd.NetworkInterfaces[0].AssociatePublicIPAddress).To(BeFalse())

			expectedIngressRules := `[
  {
    "Description": "[IngressInterCluster] Allow worker nodes in group ng-abcd1234 to communicate with control plane (kubelet and workload TCP ports)",
    "FromPort": 1025,
    "IpProtocol": "tcp",
    "SourceSecurityGroupId": {
      "Fn::ImportValue": "eksctl-test-public-ng::SecurityGroup"
    },
    "ToPort": 65535
  },
  {
    "Description": "[IngressInterClusterAPI] Allow worker nodes in group ng-abcd1234 to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)",
    "FromPort": 443,
    "IpProtocol": "tcp",
    "SourceSecurityGroupId": {
      "Fn::ImportValue": "eksctl-test-public-ng::SecurityGroup"
    },
    "ToPort": 443
  },
  {
    "CidrIp": "0.0.0.0/0",
    "Description": "Allow SSH access to worker nodes in group ng-abcd1234",
    "FromPort": 22,
    "IpProtocol": "tcp",
    "ToPort": 22
  },
  {
    "CidrIpv6": "::/0",
    "Description": "Allow SSH access to worker nodes in group ng-abcd1234",
    "FromPort": 22,
    "IpProtocol": "tcp",
    "ToPort": 22
  }
]`

			assertSSHRules(expectedIngressRules)
		})
	})

	Context("NodeGroup{PrivateNetworking=false SSH.Allow=false}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)
		disable := api.ClusterDisableNAT

		cfg.VPC = &api.ClusterVPC{
			Network: api.Network{
				ID: vpcID,
			},
			NAT: &api.ClusterNAT{
				Gateway: &disable,
			},
			SecurityGroup: "sg-0b44c48bcba5b7362",
			Subnets: &api.ClusterSubnets{
				Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"us-west-2b": {
						ID: "subnet-0f98135715dfcf55f",
					},
					"us-west-2a": {
						ID: "subnet-0ade11bad78dced9e",
					},
					"us-west-2c": {
						ID: "subnet-0e2e63ff1712bf6ef",
					},
				}),
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"us-west-2b": {
						ID: "subnet-0f98135715dfcf55a",
					},
					"us-west-2a": {
						ID: "subnet-0ade11bad78dced9f",
					},
					"us-west-2c": {
						ID: "subnet-0e2e63ff1712bf6ea",
					},
				}),
			},
		}

		ng.AvailabilityZones = []string{testAZs[1]}
		ng.SSH.Allow = api.Disabled()
		ng.InstanceType = "t2.medium"
		ng.PrivateNetworking = false
		ng.AMIFamily = "AmazonLinux2"

		It("should have 1 AZ for the nodegroup", func() {
			Expect(ng.AvailabilityZones).To(Equal([]string{"us-west-2a"}))
		})

		build(cfg, "eksctl-test-public-ng", ng)

		roundtrip()

		It("should have correct description", func() {
			Expect(ngTemplate.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
			Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have correct resources and attributes", func() {
			Expect(ngTemplate.Resources).ToNot(BeEmpty())

			Expect(ngTemplate.Resources["NodeGroup"].Properties.VPCZoneIdentifier).ToNot(BeNil())
			x, ok := ngTemplate.Resources["NodeGroup"].Properties.VPCZoneIdentifier.([]interface{})
			Expect(ok).To(BeTrue())
			refSubnets := []interface{}{
				cfg.VPC.Subnets.Public["us-west-2a"].ID,
			}
			Expect(x).To(Equal(refSubnets))

			ltd := getLaunchTemplateData(ngTemplate)

			Expect(ltd.InstanceType).To(Equal("t2.medium"))

			Expect(ltd.NetworkInterfaces).To(HaveLen(1))
			Expect(ltd.NetworkInterfaces[0].DeviceIndex).To(Equal(0))
			Expect(ltd.NetworkInterfaces[0].AssociatePublicIPAddress).To(BeFalse())

			assertSSHRules(`[
  {
    "Description": "[IngressInterCluster] Allow worker nodes in group ng-abcd1234 to communicate with control plane (kubelet and workload TCP ports)",
    "FromPort": 1025,
    "IpProtocol": "tcp",
    "SourceSecurityGroupId": {
      "Fn::ImportValue": "eksctl-test-public-ng::SecurityGroup"
    },
    "ToPort": 65535
  },
  {
    "Description": "[IngressInterClusterAPI] Allow worker nodes in group ng-abcd1234 to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)",
    "FromPort": 443,
    "IpProtocol": "tcp",
    "SourceSecurityGroupId": {
      "Fn::ImportValue": "eksctl-test-public-ng::SecurityGroup"
    },
    "ToPort": 443
  }
]`)
		})
	})

	Context("NodeGroup{EBSOptimized=nil}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		build(cfg, "eksctl-test-ebs-optimized", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).EbsOptimized).To(BeNil())
		})
	})

	Context("NodeGroup{EBSOptimized=false}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.EBSOptimized = api.Disabled()

		build(cfg, "eksctl-test-ebs-optimized", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).EbsOptimized).ToNot(BeNil())
			Expect(*getLaunchTemplateData(ngTemplate).EbsOptimized).To(BeFalse())
		})
	})

	Context("NodeGroup{EBSOptimized=true}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.EBSOptimized = api.Enabled()

		build(cfg, "eksctl-test-ebs-optimized", ng)

		roundtrip()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(ngTemplate).EbsOptimized).ToNot(BeNil())
			Expect(*getLaunchTemplateData(ngTemplate).EbsOptimized).To(BeTrue())
		})
	})

	Context("UserData", func() {
		When("ami family is AmazonLinux2", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.AMIFamily = "AmazonLinux2"

			build(cfg, "eksctl-test-123-cluster", ng)

			roundtrip()

			It("userdata should not be empty", func() {
				Expect(getLaunchTemplateData(ngTemplate).UserData).ToNot(BeEmpty())
			})

			It("should have correct description", func() {
				Expect(ngTemplate.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
				Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
				Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
			})
		})

		When("ami family is Ubuntu1804", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.AMIFamily = "Ubuntu1804"

			build(cfg, "eksctl-test-123-cluster", ng)

			roundtrip()

			It("userdata should not be empty", func() {
				Expect(getLaunchTemplateData(ngTemplate).UserData).ToNot(BeEmpty())
			})

			It("should have correct description", func() {
				Expect(ngTemplate.Description).To(ContainSubstring("AMI family: Ubuntu1804"))
				Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
				Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
			})
		})

		When("ami family is Ubuntu2004", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.AMIFamily = "Ubuntu2004"

			build(cfg, "eksctl-test-123-cluster", ng)

			roundtrip()

			It("userdata should not be empty", func() {
				Expect(getLaunchTemplateData(ngTemplate).UserData).ToNot(BeEmpty())
			})

			It("should have correct description", func() {
				Expect(ngTemplate.Description).To(ContainSubstring("AMI family: Ubuntu2004"))
				Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
				Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
			})
		})

		When("ami family is Bottlerocket", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.AMIFamily = "Bottlerocket"

			build(cfg, "eksctl-test-123-cluster", ng)

			roundtrip()

			It("userdata should not be empty", func() {
				Expect(getLaunchTemplateData(ngTemplate).UserData).ToNot(BeEmpty())
			})

			It("should have correct description", func() {
				Expect(ngTemplate.Description).To(ContainSubstring("AMI family: Bottlerocket"))
				Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
				Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
			})
		})

		When("ami family is WindowsServer2019CoreContainer", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.AMIFamily = "WindowsServer2019CoreContainer"

			build(cfg, "eksctl-test-123-cluster", ng)

			roundtrip()

			It("userdata should not be empty", func() {
				Expect(getLaunchTemplateData(ngTemplate).UserData).ToNot(BeEmpty())
			})

			It("should have correct description", func() {
				Expect(ngTemplate.Description).To(ContainSubstring("AMI family: WindowsServer2019CoreContainer"))
				Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
				Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
			})
		})

		When("ami family is WindowsServer2019FullContainer", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.AMIFamily = "WindowsServer2019FullContainer"

			build(cfg, "eksctl-test-123-cluster", ng)

			roundtrip()

			It("userdata should not be empty", func() {
				Expect(getLaunchTemplateData(ngTemplate).UserData).ToNot(BeEmpty())
			})

			It("should have correct description", func() {
				Expect(ngTemplate.Description).To(ContainSubstring("AMI family: WindowsServer2019FullContainer"))
				Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
				Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
			})
		})

		When("ami family is WindowsServer2004CoreContainer", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.AMIFamily = "WindowsServer2004CoreContainer"

			build(cfg, "eksctl-test-123-cluster", ng)

			roundtrip()

			It("userdata should not be empty", func() {
				Expect(getLaunchTemplateData(ngTemplate).UserData).ToNot(BeEmpty())
			})

			It("should have correct description", func() {
				Expect(ngTemplate.Description).To(ContainSubstring("AMI family: WindowsServer2004CoreContainer"))
				Expect(ngTemplate.Description).To(ContainSubstring("SSH access: false"))
				Expect(ngTemplate.Description).To(ContainSubstring("private networking: false"))
			})
		})
	})

	Context("with Fargate profiles", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)
		name := "test-fargate-profile"
		cfg.Metadata.Name = name
		cfg.FargateProfiles = []*api.FargateProfile{
			// default fargate profile
			{
				Name: "fp-default",
				Selectors: []api.FargateProfileSelector{
					{Namespace: "default"},
					{Namespace: "kube-system"},
				},
			},
		}
		build(cfg, fmt.Sprintf("eksctl-%s-cluster", name), ng)
		roundtrip()

		It("should have the Fargate pod execution role", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlane"))
			Expect(clusterTemplate.Resources).To(HaveKey("ServiceRole"))
			Expect(clusterTemplate.Resources).To(HaveKey("PolicyCloudWatchMetrics"))
			Expect(clusterTemplate.Resources).To(HaveKey("FargatePodExecutionRole"))
			Expect(clusterTemplate.Resources).To(HaveLen(5))
		})
	})

	Context("without VPC and IAM", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		cfg.Metadata.Name = "test-1"

		role1 := "role-1"

		cfg.IAM = &api.ClusterIAM{
			ServiceRoleARN: &role1,
		}

		build(cfg, "eksctl-test-1-cluster", ng)

		roundtrip()

		It("should only have EKS resources", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlane"))
			Expect(clusterTemplate.Resources).To(HaveLen(1))

			cp := clusterTemplate.Resources["ControlPlane"].Properties

			Expect(cp.Name).To(Equal(cfg.Metadata.Name))

			Expect(cp.RoleArn).To(Equal(role1))

			Expect(cp.ResourcesVpcConfig.SecurityGroupIds).To(HaveLen(1))
			Expect(cp.ResourcesVpcConfig.SecurityGroupIds[0]).To(Equal(cfg.VPC.SecurityGroup))

			Expect(cp.ResourcesVpcConfig.SubnetIds).To(HaveLen(6))
		})

	})

	Context("without VPC", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		cfg.Metadata.Name = "test-SharedVPC"

		build(cfg, "eksctl-test-SharedVPC-cluster", ng)

		roundtrip()

		It("should have EKS and IAM resources", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlane"))
			Expect(clusterTemplate.Resources).To(HaveKey("ServiceRole"))
			Expect(clusterTemplate.Resources).To(HaveKey("PolicyCloudWatchMetrics"))
			Expect(clusterTemplate.Resources).To(HaveLen(4))
		})

		It("should have correct own IAM resources", func() {
			Expect(clusterTemplate.Resources["ServiceRole"].Properties).ToNot(BeNil())

			Expect(clusterTemplate.Resources["ServiceRole"].Properties.ManagedPolicyArns).To(Equal(
				makePolicyARNRef("AmazonEKSClusterPolicy", "AmazonEKSVPCResourceController")),
			)

			checkARPD([]string{"EKS"}, clusterTemplate.Resources["ServiceRole"].Properties.AssumeRolePolicyDocument)

			policy1 := clusterTemplate.Resources["PolicyCloudWatchMetrics"].Properties

			Expect(policy1).ToNot(BeNil())
			isRefTo(policy1.Roles[0], "ServiceRole")

			Expect(policy1.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy1.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy1.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy1.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"cloudwatch:PutMetricData",
			}))
		})

		It("should use own IAM role and given VPC and subnets", func() {
			cp := clusterTemplate.Resources["ControlPlane"].Properties

			Expect(cp.Name).To(Equal(cfg.Metadata.Name))

			isFnGetAttOf(cp.RoleArn, "ServiceRole", "Arn")

			Expect(cp.ResourcesVpcConfig.SecurityGroupIds).To(HaveLen(1))
			Expect(cp.ResourcesVpcConfig.SecurityGroupIds[0]).To(Equal(cfg.VPC.SecurityGroup))

			Expect(cp.ResourcesVpcConfig.SubnetIds).To(HaveLen(6))
			subnetIDs := sets.NewString()
			for _, subnet := range cp.ResourcesVpcConfig.SubnetIds {
				subnetIDs.Insert(subnet.(string))
			}
			Expect(subnetIDs.HasAll(strings.Split(subnetsPublic, ",")...)).To(BeTrue())
			Expect(subnetIDs.HasAll(strings.Split(subnetsPrivate, ",")...)).To(BeTrue())
		})
	})

	Context("VPC with default CIDR", func() {
		cfg, ng := newClusterConfigAndNodegroup(false)

		cfg.Metadata.Name = "test-OwnVPC"

		setSubnets(cfg)

		build(cfg, "eksctl-test-OwnVPC-cluster", ng)

		roundtrip()

		It("should have correct own VPC resources and properties", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("VPC"))

			Expect(clusterTemplate.Resources).ToNot(HaveKey("AutoAllocatedCIDRv6"))

			Expect(clusterTemplate.Resources).To(HaveKey("InternetGateway"))
			Expect(clusterTemplate.Resources).To(HaveKey("NATIP"))
			Expect(clusterTemplate.Resources).To(HaveKey("NATGateway"))
			Expect(clusterTemplate.Resources).To(HaveKey("VPCGatewayAttachment"))

			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlaneSecurityGroup"))
			Expect(clusterTemplate.Resources).To(HaveKey("ClusterSharedNodeSecurityGroup"))

			Expect(clusterTemplate.Resources).To(HaveKey("PublicRouteTable"))
			Expect(clusterTemplate.Resources).To(HaveKey("PublicSubnetRoute"))

			for _, suffix1 := range []string{"PrivateUSWEST2", "PublicUSWEST2"} {
				for _, suffix2 := range []string{"A", "B", "C"} {
					suffix := suffix1 + suffix2
					Expect(clusterTemplate.Resources).To(HaveKey("Subnet" + suffix))
					Expect(clusterTemplate.Resources).To(HaveKey("RouteTableAssociation" + suffix))
					isRefTo(clusterTemplate.Resources["RouteTableAssociation"+suffix].Properties.SubnetID, "Subnet"+suffix)
					Expect(clusterTemplate.Resources).ToNot(HaveKey(suffix + "CIDRv6"))
				}
			}

			Expect(len(clusterTemplate.Resources)).To(Equal(32))
		})

		It("should use own VPC and subnets", func() {
			cp := clusterTemplate.Resources["ControlPlane"].Properties

			Expect(cp.Name).To(Equal(cfg.Metadata.Name))

			Expect(cp.ResourcesVpcConfig.SecurityGroupIds).To(HaveLen(1))
			isRefTo(cp.ResourcesVpcConfig.SecurityGroupIds[0], "ControlPlaneSecurityGroup")

			Expect(cp.ResourcesVpcConfig.SubnetIds).To(HaveLen(6))
			subnetRefs := sets.NewString()
			for _, subnet := range cp.ResourcesVpcConfig.SubnetIds {
				Expect(subnet.(map[string]interface{})).To(HaveKey("Ref"))
				subnetRefs.Insert(subnet.(map[string]interface{})["Ref"].(string))
			}
			for _, suffix1 := range []string{"PrivateUSWEST2", "PublicUSWEST2"} {
				for _, suffix2 := range []string{"A", "B", "C"} {
					Expect(subnetRefs.Has("Subnet" + suffix1 + suffix2)).To(BeTrue())
					subnet := clusterTemplate.Resources["Subnet"+suffix1+suffix2].Properties
					Expect(subnet.Tags).To(HaveLen(2))
					isRefTo(subnet.VpcID, "VPC")
					Expect(subnet.AvailabilityZone).To(HavePrefix("us-west-2"))
					Expect(subnet.CidrBlock).To(HavePrefix("192.168."))
					Expect(subnet.CidrBlock).To(HaveSuffix(".0/19"))
				}
			}
		})

		It("should route Internet traffic from private subnets through the single NAT gateway", func() {
			zones := []string{"A", "B", "C"}
			region := "USWEST2"

			for _, zone := range zones {
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.NatGatewayID, "NATGateway")
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.SubnetID, "SubnetPrivate"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
			}
		})
	})

	Context("VPC with custom CIDR and IPv6", func() {
		cfg, ng := newClusterConfigAndNodegroup(false)

		cfg.VPC.CIDR, _ = ipnet.ParseCIDR("10.2.0.0/16")

		cfg.VPC.AutoAllocateIPv6 = api.Enabled()

		setSubnets(cfg)

		build(cfg, "eksctl-test-VPCIPv6-cluster", ng)

		roundtrip()

		It("should have correct own VPC resources and properties", func() {
			Expect(clusterTemplate.Resources).To(HaveKey("VPC"))

			Expect(clusterTemplate.Resources).To(HaveKey("AutoAllocatedCIDRv6"))
			isRefTo(clusterTemplate.Resources["AutoAllocatedCIDRv6"].Properties.VpcID, "VPC")
			Expect(clusterTemplate.Resources["AutoAllocatedCIDRv6"].Properties.AmazonProvidedIpv6CidrBlock).To(BeTrue())

			Expect(clusterTemplate.Resources).To(HaveKey("InternetGateway"))
			Expect(clusterTemplate.Resources).To(HaveKey("NATIP"))
			Expect(clusterTemplate.Resources).To(HaveKey("NATGateway"))
			Expect(clusterTemplate.Resources).To(HaveKey("VPCGatewayAttachment"))

			Expect(clusterTemplate.Resources).To(HaveKey("ControlPlaneSecurityGroup"))
			Expect(clusterTemplate.Resources).To(HaveKey("ClusterSharedNodeSecurityGroup"))

			Expect(clusterTemplate.Resources).To(HaveKey("PublicRouteTable"))
			Expect(clusterTemplate.Resources).To(HaveKey("PublicSubnetRoute"))
			Expect(clusterTemplate.Resources["PublicSubnetRoute"].DependsOn).To(
				BeEquivalentTo([]string{"VPCGatewayAttachment"}),
			)

			expectedFnCIDR := `{ "Fn::Cidr": [{ "Fn::Select": [ 0, { "Fn::GetAtt": ["VPC", "Ipv6CidrBlocks"] }]}, 8, 64 ]}`

			for _, suffix1 := range []string{"PrivateUSWEST2", "PublicUSWEST2"} {
				for _, suffix2 := range []string{"A", "B", "C"} {
					suffix := suffix1 + suffix2
					Expect(clusterTemplate.Resources).To(HaveKey("Subnet" + suffix))
					Expect(clusterTemplate.Resources).To(HaveKey("RouteTableAssociation" + suffix))
					isRefTo(clusterTemplate.Resources["RouteTableAssociation"+suffix].Properties.SubnetID, "Subnet"+suffix)
					Expect(clusterTemplate.Resources).To(HaveKey(suffix + "CIDRv6"))

					cidr := clusterTemplate.Resources[suffix+"CIDRv6"].Properties
					isRefTo(cidr.SubnetID, "Subnet"+suffix)
					Expect(cidr.Ipv6CidrBlock["Fn::Select"]).To(HaveLen(2))
					Expect(cidr.Ipv6CidrBlock["Fn::Select"][0].(float64) >= 0)
					Expect(cidr.Ipv6CidrBlock["Fn::Select"][0].(float64) < 8)

					actualFnCIDR, _ := json.Marshal(cidr.Ipv6CidrBlock["Fn::Select"][1])
					Expect(actualFnCIDR).To(MatchJSON([]byte(expectedFnCIDR)))
				}
			}

			Expect(len(clusterTemplate.Resources)).To(Equal(39))
		})

		It("should use own VPC and subnets", func() {
			cp := clusterTemplate.Resources["ControlPlane"].Properties

			Expect(cp.Name).To(Equal(cfg.Metadata.Name))

			Expect(cp.ResourcesVpcConfig.SecurityGroupIds).To(HaveLen(1))
			isRefTo(cp.ResourcesVpcConfig.SecurityGroupIds[0], "ControlPlaneSecurityGroup")

			Expect(cp.ResourcesVpcConfig.SubnetIds).To(HaveLen(6))
			subnetRefs := sets.NewString()
			for _, subnet := range cp.ResourcesVpcConfig.SubnetIds {
				Expect(subnet.(map[string]interface{})).To(HaveKey("Ref"))
				subnetRefs.Insert(subnet.(map[string]interface{})["Ref"].(string))
			}
			for _, suffix1 := range []string{"PrivateUSWEST2", "PublicUSWEST2"} {
				for _, suffix2 := range []string{"A", "B", "C"} {
					Expect(subnetRefs.Has("Subnet" + suffix1 + suffix2)).To(BeTrue())
					subnet := clusterTemplate.Resources["Subnet"+suffix1+suffix2].Properties
					Expect(subnet.Tags).To(HaveLen(2))
					isRefTo(subnet.VpcID, "VPC")
					Expect(subnet.AvailabilityZone).To(HavePrefix("us-west-2"))
					Expect(subnet.CidrBlock).To(HavePrefix("10.2."))
					Expect(subnet.CidrBlock).To(HaveSuffix(".0/19"))
				}
			}
		})

		It("should route Internet traffic from private subnets through the single NAT gateway", func() {
			zones := []string{"A", "B", "C"}
			region := "USWEST2"

			for _, zone := range zones {
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.NatGatewayID, "NATGateway")
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.SubnetID, "SubnetPrivate"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
			}
		})

	})

	Context("VPC with highly available NAT gateways", func() {

		zones := []string{"A", "B", "C"}
		region := "USWEST2"

		cfg, ng := newClusterConfigAndNodegroup(false)

		cfg.Metadata.Name = "test-HA-NAT-VPC"

		highly := api.ClusterHighlyAvailableNAT
		cfg.VPC.NAT = &api.ClusterNAT{
			Gateway: &highly,
		}

		setSubnets(cfg)

		build(cfg, "eksctl-test-HA-NAT-VPC-cluster", ng)

		roundtrip()

		It("should have correct HA NAT VPC resources", func() {

			Expect(clusterTemplate.Resources).To(HaveKey("VPC"))

			for _, zone := range zones {
				Expect(clusterTemplate.Resources).To(HaveKey("NATIP" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("NATGateway" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("SubnetPrivate" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("PrivateRouteTable" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("NATPrivateSubnetRoute" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("RouteTableAssociationPrivate" + region + zone))
			}

			Expect(len(clusterTemplate.Resources)).To(Equal(36))
		})

		It("should route Internet traffic from private subnets through their corresponding NAT gateways", func() {
			for _, zone := range zones {
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.NatGatewayID, "NATGateway"+region+zone)
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.SubnetID, "SubnetPrivate"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
			}
		})
	})

	Context("VPC with single NAT gateway", func() {

		zones := []string{"A", "B", "C"}
		region := "USWEST2"

		cfg, ng := newClusterConfigAndNodegroup(false)

		cfg.Metadata.Name = "test-single-NAT-VPC"

		single := api.ClusterSingleNAT
		cfg.VPC.NAT = &api.ClusterNAT{
			Gateway: &single,
		}

		setSubnets(cfg)

		build(cfg, "eksctl-test-single-NAT-VPC-cluster", ng)

		roundtrip()

		It("should have correct single NAT VPC resources", func() {

			Expect(clusterTemplate.Resources).To(HaveKey("VPC"))

			Expect(clusterTemplate.Resources).To(HaveKey("NATIP"))
			Expect(clusterTemplate.Resources).To(HaveKey("NATGateway"))

			for _, zone := range zones {
				Expect(clusterTemplate.Resources).To(HaveKey("SubnetPrivate" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("PrivateRouteTable" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("RouteTableAssociationPrivate" + region + zone))
			}

			Expect(len(clusterTemplate.Resources)).To(Equal(32))
		})

		It("should route Internet traffic from private subnets through the single NAT gateway", func() {
			for _, zone := range zones {
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.NatGatewayID, "NATGateway")
				isRefTo(clusterTemplate.Resources["NATPrivateSubnetRoute"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.SubnetID, "SubnetPrivate"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
			}
		})
	})

	Context("VPC with no NAT gateway", func() {

		zones := []string{"A", "B", "C"}
		region := "USWEST2"

		cfg, ng := newClusterConfigAndNodegroup(false)

		cfg.Metadata.Name = "test-single-NAT-VPC"

		disable := api.ClusterDisableNAT
		cfg.VPC.NAT = &api.ClusterNAT{
			Gateway: &disable,
		}

		setSubnets(cfg)

		build(cfg, "eksctl-test-disabled-NAT-VPC-cluster", ng)

		roundtrip()

		It("should have correct disabled NAT VPC resources and properties", func() {

			Expect(clusterTemplate.Resources).To(HaveKey("VPC"))

			for _, zone := range zones {
				Expect(clusterTemplate.Resources).To(HaveKey("SubnetPrivate" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("PrivateRouteTable" + region + zone))
				Expect(clusterTemplate.Resources).To(HaveKey("RouteTableAssociationPrivate" + region + zone))
			}

			Expect(len(clusterTemplate.Resources)).To(Equal(27))

		})

		It("should not route Internet traffic from private subnets", func() {

			for _, zone := range zones {
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.SubnetID, "SubnetPrivate"+region+zone)
				isRefTo(clusterTemplate.Resources["RouteTableAssociationPrivate"+region+zone].Properties.RouteTableID, "PrivateRouteTable"+region+zone)
			}
		})

	})

	maxSpotPrice := 0.045
	baseCap := 40
	percentageOnDemand := 20
	pools := 3
	spotAllocationStrategy := "lowest-price"
	zero := 0
	cpuCreditsUnlimited := "unlimited"
	cpuCreditsStandard := "standard"

	Context("Nodegroup with Mixed instances", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "mixed"
		ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
			MaxPrice:                            &maxSpotPrice,
			InstanceTypes:                       []string{"m5.large", "m5a.xlarge"},
			OnDemandBaseCapacity:                &baseCap,
			OnDemandPercentageAboveBaseCapacity: &percentageOnDemand,
			SpotInstancePools:                   &pools,
			SpotAllocationStrategy:              &spotAllocationStrategy,
		}

		ng.MinSize = &zero
		ng.MaxSize = &zero

		build(cfg, "eksctl-test-spot-cluster", ng)

		roundtrip()

		It("should have mixed instances with correct max price", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeGroupLaunchTemplate"))

			launchTemplateData := getLaunchTemplateData(ngTemplate)
			Expect(launchTemplateData.InstanceType).To(Equal("m5.large"))
			Expect(launchTemplateData.InstanceMarketOptions).To(BeNil())

			nodeGroupProperties := getNodeGroupProperties(ngTemplate)
			Expect(nodeGroupProperties.MinSize).To(Equal("0"))
			Expect(nodeGroupProperties.MaxSize).To(Equal("0"))
			Expect(nodeGroupProperties.DesiredCapacity).To(Equal(""))

			Expect(nodeGroupProperties.MixedInstancesPolicy).To(Not(BeNil()))
			Expect(nodeGroupProperties.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.LaunchTemplateName["Fn::Sub"]).To(Equal("${AWS::StackName}"))
			Expect(nodeGroupProperties.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.Version["Fn::GetAtt"]).To(Equal(
				[]interface{}{"NodeGroupLaunchTemplate", "LatestVersionNumber"}),
			)
			Expect(nodeGroupProperties.MixedInstancesPolicy.LaunchTemplate).To(Not(BeNil()))

			Expect(nodeGroupProperties.MixedInstancesPolicy.InstancesDistribution).To(Not(BeNil()))
			Expect(nodeGroupProperties.MixedInstancesPolicy.InstancesDistribution.OnDemandBaseCapacity).To(Equal("40"))
			Expect(nodeGroupProperties.MixedInstancesPolicy.InstancesDistribution.OnDemandPercentageAboveBaseCapacity).To(Equal("20"))
			Expect(nodeGroupProperties.MixedInstancesPolicy.InstancesDistribution.SpotInstancePools).To(Equal("3"))
			Expect(nodeGroupProperties.MixedInstancesPolicy.InstancesDistribution.SpotMaxPrice).To(Equal("0.045000"))
			Expect(nodeGroupProperties.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy).To(Equal("lowest-price"))

		})
	})

	Context("NodeGroup{CPUCredits=nil}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		build(cfg, "eksctl-test-t3-unlimited", ng)

		roundtrip()

		It("should have correct resources and attributes", func() {
			Expect(getLaunchTemplateData(ngTemplate).CreditSpecification).To(BeNil())
		})
	})

	Context("NodeGroup{CPUCredits=standard InstancesDistribution.InstanceTypes=t3.medium,t3a.medium}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "mixed"
		ng.CPUCredits = &cpuCreditsStandard
		ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
			MaxPrice:                            &maxSpotPrice,
			InstanceTypes:                       []string{"t3.medium", "t3a.medium"},
			OnDemandBaseCapacity:                &baseCap,
			OnDemandPercentageAboveBaseCapacity: &percentageOnDemand,
			SpotInstancePools:                   &pools,
			SpotAllocationStrategy:              &spotAllocationStrategy,
		}

		build(cfg, "eksctl-test-t3-unlimited", ng)

		roundtrip()

		It("should have correct resources and attributes", func() {
			Expect(getLaunchTemplateData(ngTemplate).CreditSpecification).ToNot(BeNil())
			Expect(getLaunchTemplateData(ngTemplate).CreditSpecification.CPUCredits).ToNot(BeNil())
			Expect(getLaunchTemplateData(ngTemplate).CreditSpecification.CPUCredits).To(Equal("standard"))
		})
	})

	Context("NodeGroup{CPUCredits=unlimited InstancesDistribution.InstanceTypes=t3.medium,t3a.medium}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "mixed"
		ng.CPUCredits = &cpuCreditsUnlimited
		ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
			MaxPrice:                            &maxSpotPrice,
			InstanceTypes:                       []string{"t3.medium", "t3a.medium"},
			OnDemandBaseCapacity:                &baseCap,
			OnDemandPercentageAboveBaseCapacity: &percentageOnDemand,
			SpotInstancePools:                   &pools,
			SpotAllocationStrategy:              &spotAllocationStrategy,
		}

		build(cfg, "eksctl-test-t3-unlimited", ng)

		roundtrip()

		It("should have correct resources and attributes", func() {
			Expect(getLaunchTemplateData(ngTemplate).CreditSpecification).ToNot(BeNil())
			Expect(getLaunchTemplateData(ngTemplate).CreditSpecification.CPUCredits).ToNot(BeNil())
			Expect(getLaunchTemplateData(ngTemplate).CreditSpecification.CPUCredits).To(Equal("unlimited"))
		})
	})

	Context("NodeGroup with asgSuspendProcesses", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.ASGSuspendProcesses = []string{"Launch", "InstanceRefresh"}
		build(cfg, "eksctl-test-asgSuspendProcesses", ng)

		roundtrip()

		It("should have correct resources and attributes", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
			ngResource := ngTemplate.Resources["NodeGroup"]
			Expect(ng).ToNot(BeNil())
			Expect(ngResource.UpdatePolicy).To(HaveKey("AutoScalingRollingUpdate"))
			Expect(ngResource.UpdatePolicy["AutoScalingRollingUpdate"]).To(
				HaveKeyWithValue("SuspendProcesses", []interface{}{"Launch", "InstanceRefresh"}),
			)
		})
		Context("empty asgSuspendProcesses", func() {
			cfg, ng := newClusterConfigAndNodegroup(true)

			ng.ASGSuspendProcesses = []string{}
			build(cfg, "eksctl-test-asgSuspendProcesses-empty", ng)

			roundtrip()

			It("shouldn't be included in resource", func() {
				Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
				ngResource := ngTemplate.Resources["NodeGroup"]
				Expect(ng).ToNot(BeNil())
				Expect(ngResource.UpdatePolicy).To(HaveKey("AutoScalingRollingUpdate"))
				Expect(ngResource.UpdatePolicy["AutoScalingRollingUpdate"]).ToNot(
					HaveKey("SuspendProcesses"),
				)
			})
		})
	})

	Context("p4 NodeGroup with EFA enabled", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.EFAEnabled = aws.Bool(true)
		ng.InstanceType = p4InstanceType
		build(cfg, "eksctl-test-nodegroup-p4-efa", ng)

		roundtrip()

		It("should have correct interfaces", func() {
			Expect(ngTemplate.Resources).To(HaveKey("NodeGroup"))
			launchTemplate := ngTemplate.Resources["NodeGroupLaunchTemplate"]
			Expect(ng).ToNot(BeNil())
			Expect(launchTemplate.Properties.LaunchTemplateData.NetworkInterfaces).To(Equal(
				[]NetworkInterface{
					{
						InterfaceType:    "efa",
						DeviceIndex:      0,
						NetworkCardIndex: 0,
					},
					{
						InterfaceType:    "efa",
						DeviceIndex:      1,
						NetworkCardIndex: 1,
					},
					{
						InterfaceType:    "efa",
						DeviceIndex:      2,
						NetworkCardIndex: 2,
					},
					{
						InterfaceType:    "efa",
						DeviceIndex:      3,
						NetworkCardIndex: 3,
					},
				},
			))
		})
	})

})

func setSubnets(cfg *api.ClusterConfig) {
	It("should not error when calling SetSubnets", func() {
		err := vpc.SetSubnets(cfg.VPC, cfg.AvailabilityZones)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should have public and private subnets", func() {
		Expect(cfg.VPC.Subnets.Private).To(HaveLen(3))
		Expect(cfg.VPC.Subnets.Public).To(HaveLen(3))
	})
}

func getNodeGroupProperties(obj *Template) Properties {
	Expect(obj.Resources).To(HaveKey("NodeGroup"))
	ng := obj.Resources["NodeGroup"]
	Expect(ng).ToNot(BeNil())
	Expect(ng.Properties).ToNot(BeNil())
	return ng.Properties
}

func isRefTo(obj interface{}, value string) {
	Expect(obj).ToNot(BeEmpty())
	o, ok := obj.(map[string]interface{})
	Expect(ok).To(BeTrue())
	Expect(o).To(HaveKey(gfnt.Ref))
	Expect(o[gfnt.Ref]).To(Equal(value))
}

func isFnGetAttOf(obj interface{}, logicalName, attr string) {
	Expect(obj).ToNot(BeEmpty())
	o, ok := obj.(map[string]interface{})
	Expect(ok).To(BeTrue())
	Expect(o).To(HaveKey(gfnt.FnGetAtt))
	Expect(o[gfnt.FnGetAtt]).To(Equal([]interface{}{logicalName, attr}))
}

func getLaunchTemplateData(obj *Template) LaunchTemplateData {
	Expect(obj.Resources).ToNot(BeEmpty())
	Expect(obj.Resources).To(HaveKey("NodeGroupLaunchTemplate"))
	return obj.Resources["NodeGroupLaunchTemplate"].Properties.LaunchTemplateData
}

func checkARPD(services []string, arpd interface{}) {
	var serviceRefs []*gfnt.Value
	for _, service := range services {
		serviceRefs = append(serviceRefs, MakeServiceRef(service))
	}
	servicesJSON, err := json.Marshal(serviceRefs)
	Expect(err).ToNot(HaveOccurred())

	expectedARPD := `{
		"Version": "2012-10-17",
		"Statement": [{
						"Action": ["sts:AssumeRole"],
						"Effect": "Allow",
						"Principal": {
								"Service": ` + string(servicesJSON) + `
						}
		}]
	}`
	actualARPD, _ := json.Marshal(arpd)
	Expect(actualARPD).To(MatchJSON([]byte(expectedARPD)))
}

func makePolicyARNRef(policies ...string) []interface{} {
	var values []interface{}
	for _, p := range policies {
		values = append(values, map[string]interface{}{
			"Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/" + p,
		})
	}
	return values
}
