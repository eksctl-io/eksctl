package builder_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"

	gfn "github.com/awslabs/goformation/cloudformation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	clusterName = "ferocious-mushroom-1532594698"
	endpoint    = "https://DE37D8AFB23F7275D2361AD6B2599143.yl4.us-west-2.eks.amazonaws.com"
	caCert      = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNE1EWXdOekExTlRBMU5Wb1hEVEk0TURZd05EQTFOVEExTlZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTWJoCnpvZElYR0drckNSZE1jUmVEN0YvMnB1NFZweTdvd3FEVDgrdk9zeGs2bXFMNWxQd3ZicFhmYkE3R0xzMDVHa0wKaDdqL0ZjcU91cnMwUFZSK3N5REtuQXltdDFORWxGNllGQktSV1dUQ1hNd2lwN1pweW9XMXdoYTlJYUlPUGxCTQpPTEVlckRabFVrVDFVV0dWeVdsMmxPeFgxa2JhV2gvakptWWdkeW5jMXhZZ3kxa2JybmVMSkkwLzVUVTRCajJxClB1emtrYW5Xd3lKbGdXQzhBSXlpWW82WFh2UVZmRzYrM3RISE5XM1F1b3ZoRng2MTFOYnl6RUI3QTdtZGNiNmgKR0ZpWjdOeThHZnFzdjJJSmI2Nk9FVzBSdW9oY1k3UDZPdnZmYnlKREhaU2hqTStRWFkxQXN5b3g4Ri9UelhHSgpQUWpoWUZWWEVhZU1wQmJqNmNFQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFCa2hKRVd4MHk1LzlMSklWdXJ1c1hZbjN6Z2EKRkZ6V0JsQU44WTlqUHB3S2t0Vy9JNFYyUGg3bWY2Z3ZwZ3Jhc2t1Slk1aHZPcDdBQmcxSTFhaHUxNUFpMUI0ZApuMllRaDlOaHdXM2pKMmhuRXk0VElpb0gza2JFdHRnUVB2bWhUQzNEYUJreEpkbmZJSEJCV1RFTTU1czRwRmxUClpzQVJ3aDc1Q3hYbjdScVU0akpKcWNPaTRjeU5qeFVpRDBqR1FaTmNiZWEyMkRCeTJXaEEzUWZnbGNScGtDVGUKRDVPS3NOWlF4MW9MZFAwci9TSmtPT1NPeUdnbVJURTIrODQxN21PRW02Z3RPMCszdWJkbXQ0aENsWEtFTTZYdwpuQWNlK0JxVUNYblVIN2ZNS3p2TDE5UExvMm5KbFU1TnlCbU1nL1pNVHVlUy80eFZmKy94WnpsQ0Q1WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	arn         = "arn:aws:eks:us-west-2:122333:cluster/" + clusterName

	vpcID          = "vpc-0e265ad953062b94b"
	subnetsPublic  = "subnet-0f98135715dfcf55f,subnet-0ade11bad78dced9e,subnet-0e2e63ff1712bf6ef"
	subnetsPrivate = "subnet-0f98135715dfcf55a,subnet-0ade11bad78dced9f,subnet-0e2e63ff1712bf6ea"
)

var overrideBootstrapCommand = "echo foo > /etc/test_foo; echo bar > /etc/test_bar; poweroff -fn;"

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
			Action   []string
			Effect   string
			Resource interface{}
		}
	}

	LaunchTemplateData LaunchTemplateData

	VPCZoneIdentifier interface{}

	TargetGroupARNs                   []string
	DesiredCapacity, MinSize, MaxSize string

	CidrIp, CidrIpv6, IpProtocol string
	FromPort, ToPort             int
}

type LaunchTemplateData struct {
	IamInstanceProfile              struct{ Arn interface{} }
	UserData, InstanceType, ImageId string
	BlockDeviceMappings             []interface{}
	NetworkInterfaces               []struct {
		DeviceIndex              int
		AssociatePublicIpAddress bool
	}
}

type Template struct {
	Description string
	Resources   map[string]struct{ Properties Properties }
}

func kubeconfigBody(authenticator string) string {
	return `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /etc/eksctl/ca.crt
    server: ` + endpoint + `
  name: ` + clusterName + `.us-west-2.eksctl.io
contexts:
- context:
    cluster: ` + clusterName + `.us-west-2.eksctl.io
    user: kubelet@` + clusterName + `.us-west-2.eksctl.io
  name: kubelet@` + clusterName + `.us-west-2.eksctl.io
current-context: kubelet@` + clusterName + `.us-west-2.eksctl.io
kind: Config
preferences: {}
users:
- name: kubelet@` + clusterName + `.us-west-2.eksctl.io
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - ` + clusterName + `
      command: ` + authenticator + `
      env: null
`
}

func testVPC() *api.ClusterVPC {
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
		SecurityGroup:           "sg-0b44c48bcba5b7362",
		SharedNodeSecurityGroup: "sg-shared",
		Subnets: &api.ClusterSubnets{
			Public: map[string]api.Network{
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
			},
			Private: map[string]api.Network{
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
			},
		},
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
		cc   *cloudconfig.CloudConfig
		crs  *ClusterResourceSet
		ngrs *NodeGroupResourceSet
		obj  *Template
		err  error

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
		*ng.VolumeType = api.NodeVolumeTypeIO1
		ng.VolumeName = new(string)
		*ng.VolumeName = "/dev/xvda"
		ng.AutoScalerEnabled = api.Disabled()

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

		p.MockEC2().On("DescribeVpcs", mock.MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
			return *input.VpcIds[0] == vpcID
		})).Return(&ec2.DescribeVpcsOutput{
			Vpcs: []*ec2.Vpc{{
				VpcId:     aws.String(vpcID),
				CidrBlock: aws.String("192.168.0.0/16"),
			}},
		}, nil)

		for t := range subnetLists {
			fn := func(list string, subnetsByAz map[string]api.Network) {
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
					fmt.Fprintf(GinkgoWriter, "%s subnets = %#v\n", t, output)
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
				Version: "1.12",
			},
			Status: &api.ClusterStatus{
				Endpoint:                 endpoint,
				CertificateAuthorityData: caCertData,
				ARN:                      arn,
			},
			AvailabilityZones: testAZs,
			VPC:               testVPC(),
			IAM: api.ClusterIAM{
				ServiceRoleARN: arn,
			},
			NodeGroups: []*api.NodeGroup{
				{
					AMI:               "",
					AMIFamily:         "AmazonLinux2",
					InstanceType:      "t2.medium",
					Name:              "ng-abcd1234",
					PrivateNetworking: false,
					AutoScalerEnabled: api.Disabled(),
					SecurityGroups: &api.NodeGroupSGs{
						WithLocal:  api.Enabled(),
						WithShared: api.Enabled(),
						AttachIDs:  []string{},
					},
					DesiredCapacity: nil,
					VolumeSize:      aws.Int(2),
					VolumeType:      aws.String(api.NodeVolumeTypeIO1),
					VolumeName:      aws.String("/dev/xvda"),
					IAM: &api.NodeGroupIAM{
						WithAddonPolicies: api.NodeGroupIAMAddonPolicies{
							ImageBuilder: api.Disabled(),
							AutoScaler:   api.Disabled(),
							ExternalDNS:  api.Disabled(),
							AppMesh:      api.Disabled(),
							EBS:          api.Disabled(),
							FSX:          api.Disabled(),
							EFS:          api.Disabled(),
							ALBIngress:   api.Disabled(),
						},
					},
					SSH: &api.NodeGroupSSH{
						Allow:         api.Disabled(),
						PublicKeyPath: &api.DefaultNodeSSHPublicKeyPath,
					},
				},
			},
		}

		cfg := newSimpleClusterConfig()

		It("should not error when calling SetSubnets", func() {
			err := vpc.SetSubnets(cfg)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have public and private subnets", func() {
			Expect(cfg.VPC.Subnets.Private).To(HaveLen(3))
			Expect(cfg.VPC.Subnets.Public).To(HaveLen(3))
		})

		sampleOutputs := map[string]string{
			"SecurityGroup":            "sg-0b44c48bcba5b7362",
			"SubnetsPublic":            subnetsPublic,
			"SubnetsPrivate":           subnetsPrivate,
			"VPC":                      vpcID,
			"Endpoint":                 endpoint,
			"CertificateAuthorityData": caCert,
			"ARN":                      arn,
			"ClusterStackName":         "",
			"SharedNodeSecurityGroup":  "sg-shared",
			"ServiceRoleARN":           arn,
		}

		It("should add all resources and collect outputs without errors", func() {
			crs = NewClusterResourceSet(p, cfg)
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

	build := func(cfg *api.ClusterConfig, name string, ng *api.NodeGroup) {
		It("should add all resources without errors", func() {
			crs = NewClusterResourceSet(p, cfg)
			err = crs.AddAllResources()
			Expect(err).ShouldNot(HaveOccurred())

			ngrs = NewNodeGroupResourceSet(p, cfg, name, ng)
			err = ngrs.AddAllResources()
			Expect(err).ShouldNot(HaveOccurred())

			t := ngrs.Template()
			Expect(t.Resources).Should(HaveKey("NodeGroup"))

			templateBody, err := t.JSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(templateBody).ShouldNot(BeEmpty())
		})
	}

	roundtript := func() {
		It("should serialise JSON without errors, and parse the teamplate", func() {
			obj = &Template{}
			templateBody, err := ngrs.RenderJSON()
			Expect(err).ShouldNot(HaveOccurred())
			err = json.Unmarshal(templateBody, obj)
			Expect(err).ShouldNot(HaveOccurred())
		})
	}

	extractCloudConfig := func() {
		It("should extract valid cloud-config using our implementation", func() {
			userData := getLaunchTemplateData(obj).UserData
			Expect(userData).ToNot(BeEmpty())
			cc, err = cloudconfig.DecodeCloudConfig(userData)
			Expect(err).ShouldNot(HaveOccurred())

		})
	}

	Context("AutoNameTag", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		It("SG should have correct tags", func() {
			Expect(obj.Resources).ToNot(BeNil())
			Expect(obj.Resources).To(HaveLen(10))
			Expect(obj.Resources["SG"].Properties.Tags).To(HaveLen(2))
			Expect(obj.Resources["SG"].Properties.Tags[0].Key).To(Equal("kubernetes.io/cluster/" + clusterName))
			Expect(obj.Resources["SG"].Properties.Tags[0].Value).To(Equal("owned"))
			Expect(obj.Resources["SG"].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(obj.Resources["SG"].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SG",
			}))
		})
	})

	Context("NodeGroupTags", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "t2.medium"
		ng.Name = "ng-abcd1234"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

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

			ngProps := getNodeGroupProperties(obj)

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

		roundtript()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(obj).MaxSize).To(Equal("2"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("2"))

		})
	})

	Context("NodeGroup DesiredCapacity=10 MaxSize=nil MinSize=nil", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.DesiredCapacity = nil
		ng.MaxSize = nil
		ng.MinSize = nil

		ng.InstanceType = "m5.2xlarge"

		build(cfg, "eksctl-test2-cluster", ng)

		roundtript()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(obj).MaxSize).To(Equal("2"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("2"))

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

		roundtript()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(obj).MaxSize).To(Equal("30"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("2"))
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

		roundtript()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(obj).MaxSize).To(Equal("90"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("90"))
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

		roundtript()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).DesiredCapacity).To(BeEmpty())
			Expect(getNodeGroupProperties(obj).MaxSize).To(Equal("91"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("61"))
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

		roundtript()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).DesiredCapacity).To(Equal("32"))
			Expect(getNodeGroupProperties(obj).MaxSize).To(Equal("92"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("32"))
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

		roundtript()

		It("should have correct instance type and sizes", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).DesiredCapacity).To(Equal("33"))
			Expect(getNodeGroupProperties(obj).MaxSize).To(Equal("33"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("31"))
		})
	})

	Context("NodeGroupAutoScaling", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.TargetGroupARNs = []string{"tg-arn-1", "tg-arn-2"}

		ng.MinSize = new(int)
		*ng.MinSize = 10
		ng.InstanceType = "m5.2xlarge"

		ng.IAM.InstanceRoleName = "a-named-role"
		ng.IAM.WithAddonPolicies.AutoScaler = api.Enabled()

		ng.AutoScalerEnabled = api.Enabled()

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		It("should have correct instance type and min size", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.2xlarge"))
			Expect(getNodeGroupProperties(obj).MinSize).To(Equal("10"))
		})

		It("should have correct instance role and profile", func() {
			Expect(obj.Resources).To(HaveKey("NodeInstanceRole"))

			role := obj.Resources["NodeInstanceRole"].Properties

			Expect(role.Path).To(Equal("/"))
			Expect(role.RoleName).To(Equal("a-named-role"))
			Expect(role.ManagedPolicyArns).To(HaveLen(3))
			Expect(role.ManagedPolicyArns[0]).To(Equal("arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"))
			Expect(role.ManagedPolicyArns[1]).To(Equal("arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"))
			Expect(role.ManagedPolicyArns[2]).To(Equal("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"))

			expectedARPD := `{
				"Version": "2012-10-17",
				"Statement": [{
								"Action": ["sts:AssumeRole"],
								"Effect": "Allow",
								"Principal": {
										"Service": ["ec2.amazonaws.com"]
								}
				}]
			}`
			actualARPD, _ := json.Marshal(role.AssumeRolePolicyDocument)
			Expect(actualARPD).To(MatchJSON([]byte(expectedARPD)))

			Expect(obj.Resources).To(HaveKey("NodeInstanceProfile"))

			profile := obj.Resources["NodeInstanceProfile"].Properties

			Expect(profile.Path).To(Equal("/"))
			Expect(profile.Roles).To(HaveLen(1))
			isRefTo(profile.Roles[0], "NodeInstanceRole")

			isFnGetAttOf(getLaunchTemplateData(obj).IamInstanceProfile.Arn, "NodeInstanceProfile.Arn")
		})

		It("should have correct policies", func() {
			Expect(obj.Resources).ToNot(BeEmpty())
			Expect(obj.Resources).To(HaveKey("PolicyAutoScaling"))

			policy := obj.Resources["PolicyAutoScaling"].Properties

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

			ngProps := getNodeGroupProperties(obj)

			Expect(ngProps.Tags).ToNot(BeNil())
			Expect(ngProps.Tags).To(Equal(expectedTags))
		})

		It("should have target groups ARNs set", func() {
			Expect(obj.Resources).To(HaveKey("NodeGroup"))
			ng := obj.Resources["NodeGroup"]
			Expect(ng).ToNot(BeNil())
			Expect(ng.Properties).ToNot(BeNil())

			Expect(ng.Properties.TargetGroupARNs).To(Equal([]string{"tg-arn-1", "tg-arn-2"}))
		})

		It("should have target groups ARNs set", func() {
			Expect(obj.Resources).To(HaveKey("NodeGroup"))
			ng := obj.Resources["NodeGroup"]
			Expect(ng).ToNot(BeNil())
			Expect(ng.Properties).ToNot(BeNil())

			Expect(ng.Properties.TargetGroupARNs).To(Equal([]string{"tg-arn-1", "tg-arn-2"}))
		})
	})

	Context("NodeGroupAppMeshExternalDNS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.AppMesh = api.Enabled()
		ng.IAM.WithAddonPolicies.ExternalDNS = api.Enabled()

		build(cfg, "eksctl-test-megaapps-cluster", ng)

		roundtript()

		It("should have correct policies", func() {
			Expect(obj.Resources).ToNot(BeEmpty())

			Expect(obj.Resources).To(HaveKey("PolicyExternalDNSChangeSet"))

			policy1 := obj.Resources["PolicyExternalDNSChangeSet"].Properties

			Expect(policy1.Roles).To(HaveLen(1))
			isRefTo(policy1.Roles[0], "NodeInstanceRole")

			Expect(policy1.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy1.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy1.PolicyDocument.Statement[0].Resource).To(Equal("arn:aws:route53:::hostedzone/*"))
			Expect(policy1.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ChangeResourceRecordSets",
			}))

			Expect(obj.Resources).To(HaveKey("PolicyExternalDNSHostedZones"))

			policy2 := obj.Resources["PolicyExternalDNSHostedZones"].Properties

			Expect(policy2.Roles).To(HaveLen(1))
			isRefTo(policy2.Roles[0], "NodeInstanceRole")

			Expect(policy2.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy2.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy2.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy2.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
			}))

			Expect(obj.Resources).To(HaveKey("PolicyAppMesh"))

			policy3 := obj.Resources["PolicyAppMesh"].Properties

			Expect(policy3.Roles).To(HaveLen(1))
			isRefTo(policy3.Roles[0], "NodeInstanceRole")

			Expect(policy3.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy3.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy3.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy3.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"appmesh:*",
			}))

			Expect(obj.Resources).ToNot(HaveKey("PolicyEBS"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyAutoScaling"))
		})

	})

	Context("NodeGroupALBIngress", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.WithAddonPolicies.ALBIngress = api.Enabled()

		build(cfg, "eksctl-test-megaapps-cluster", ng)

		roundtript()

		It("should have correct policies", func() {
			Expect(obj.Resources).ToNot(BeEmpty())

			Expect(obj.Resources).To(HaveKey("PolicyALBIngress"))

			policy := obj.Resources["PolicyALBIngress"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"acm:DescribeCertificate",
				"acm:ListCertificates",
				"acm:GetCertificate",
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:CreateSecurityGroup",
				"ec2:CreateTags",
				"ec2:DeleteTags",
				"ec2:DeleteSecurityGroup",
				"ec2:DescribeAccountAttributes",
				"ec2:DescribeAddresses",
				"ec2:DescribeInstances",
				"ec2:DescribeInstanceStatus",
				"ec2:DescribeInternetGateways",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DescribeSecurityGroups",
				"ec2:DescribeSubnets",
				"ec2:DescribeTags",
				"ec2:DescribeVpcs",
				"ec2:ModifyInstanceAttribute",
				"ec2:ModifyNetworkInterfaceAttribute",
				"ec2:RevokeSecurityGroupIngress",
				"elasticloadbalancing:AddListenerCertificates",
				"elasticloadbalancing:AddTags",
				"elasticloadbalancing:CreateListener",
				"elasticloadbalancing:CreateLoadBalancer",
				"elasticloadbalancing:CreateRule",
				"elasticloadbalancing:CreateTargetGroup",
				"elasticloadbalancing:DeleteListener",
				"elasticloadbalancing:DeleteLoadBalancer",
				"elasticloadbalancing:DeleteRule",
				"elasticloadbalancing:DeleteTargetGroup",
				"elasticloadbalancing:DeregisterTargets",
				"elasticloadbalancing:DescribeListenerCertificates",
				"elasticloadbalancing:DescribeListeners",
				"elasticloadbalancing:DescribeLoadBalancers",
				"elasticloadbalancing:DescribeLoadBalancerAttributes",
				"elasticloadbalancing:DescribeRules",
				"elasticloadbalancing:DescribeSSLPolicies",
				"elasticloadbalancing:DescribeTags",
				"elasticloadbalancing:DescribeTargetGroups",
				"elasticloadbalancing:DescribeTargetGroupAttributes",
				"elasticloadbalancing:DescribeTargetHealth",
				"elasticloadbalancing:ModifyListener",
				"elasticloadbalancing:ModifyLoadBalancerAttributes",
				"elasticloadbalancing:ModifyRule",
				"elasticloadbalancing:ModifyTargetGroup",
				"elasticloadbalancing:ModifyTargetGroupAttributes",
				"elasticloadbalancing:RegisterTargets",
				"elasticloadbalancing:RemoveListenerCertificates",
				"elasticloadbalancing:RemoveTags",
				"elasticloadbalancing:SetIpAddressType",
				"elasticloadbalancing:SetSecurityGroups",
				"elasticloadbalancing:SetSubnets",
				"elasticloadbalancing:SetWebACL",
				"iam:CreateServiceLinkedRole",
				"iam:GetServerCertificate",
				"iam:ListServerCertificates",
				"waf-regional:GetWebACLForResource",
				"waf-regional:GetWebACL",
				"waf-regional:AssociateWebACL",
				"waf-regional:DisassociateWebACL",
				"tag:GetResources",
				"tag:TagResources",
				"waf:GetWebACL",
			}))
		})

	})

	Context("NodeGroupEBS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeSize = nil
		ng.IAM.WithAddonPolicies.EBS = api.Enabled()

		build(cfg, "eksctl-test-ebs-cluster", ng)

		roundtript()

		It("should have correct policies", func() {
			Expect(getLaunchTemplateData(obj).BlockDeviceMappings).To(HaveLen(0))

			Expect(obj.Resources).To(HaveKey("PolicyEBS"))

			policy := obj.Resources["PolicyEBS"].Properties

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
				"ec2:DescribeInstances",
				"ec2:DescribeSnapshots",
				"ec2:DescribeTags",
				"ec2:DescribeVolumes",
				"ec2:DetachVolume",
			}))

			Expect(obj.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyExternalDNSChangeSet"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyExternalDNSHostedZones"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyAppMesh"))
		})
	})

	Context("NodeGroupFSX", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeSize = nil
		ng.IAM.WithAddonPolicies.FSX = api.Enabled()

		build(cfg, "eksctl-test-fsx-cluster", ng)

		roundtript()

		It("should have correct policies", func() {
			Expect(getLaunchTemplateData(obj).BlockDeviceMappings).To(HaveLen(0))

			Expect(obj.Resources).To(HaveKey("PolicyFSX"))

			policy := obj.Resources["PolicyFSX"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"fsx:*",
			}))

			Expect(obj.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyExternalDNSChangeSet"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyExternalDNSHostedZones"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyAppMesh"))
		})
	})

	Context("NodeGroupEFS", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.VolumeSize = nil
		ng.IAM.WithAddonPolicies.EFS = api.Enabled()

		build(cfg, "eksctl-test-efs-cluster", ng)

		roundtript()

		It("should have correct policies", func() {
			Expect(getLaunchTemplateData(obj).BlockDeviceMappings).To(HaveLen(0))

			Expect(obj.Resources).To(HaveKey("PolicyEFS"))

			policy := obj.Resources["PolicyEFS"].Properties

			Expect(policy.Roles).To(HaveLen(1))
			isRefTo(policy.Roles[0], "NodeInstanceRole")

			Expect(policy.PolicyDocument.Statement).To(HaveLen(1))
			Expect(policy.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(policy.PolicyDocument.Statement[0].Resource).To(Equal("arn:aws:elasticfilesystem:us-west-2:123456789012:file-system/*"))
			Expect(policy.PolicyDocument.Statement[0].Action).To(Equal([]string{
				"elasticfilesystem:*",
			}))

			Expect(obj.Resources).ToNot(HaveKey("PolicyAutoScaling"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyExternalDNSChangeSet"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyExternalDNSHostedZones"))
			Expect(obj.Resources).ToNot(HaveKey("PolicyAppMesh"))
		})
	})

	Context("NodeGroup with cutom role and profile", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.InstanceRoleARN = "arn:role"
		ng.IAM.InstanceProfileARN = "arn:profile"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		It("should have correct instance role and profile", func() {
			Expect(obj.Resources).ToNot(HaveKey("NodeInstanceRole"))
			Expect(obj.Resources).ToNot(HaveKey("NodeInstanceProfile"))

			Expect(getLaunchTemplateData(obj).IamInstanceProfile.Arn).To(Equal("arn:profile"))
		})
	})

	Context("NodeGroup with cutom role", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.InstanceRoleARN = "arn:role"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		It("should have correct instance role and profile", func() {
			Expect(obj.Resources).ToNot(HaveKey("NodeInstanceRole"))

			Expect(obj.Resources).To(HaveKey("NodeInstanceProfile"))

			profile := obj.Resources["NodeInstanceProfile"].Properties

			Expect(profile.Path).To(Equal("/"))
			Expect(profile.Roles).To(HaveLen(1))
			Expect(profile.Roles[0]).To(Equal("arn:role"))

			isFnGetAttOf(getLaunchTemplateData(obj).IamInstanceProfile.Arn, "NodeInstanceProfile.Arn")
		})
	})

	Context("NodeGroup with cutom profile", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.IAM.InstanceProfileARN = "arn:profile"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		It("should have correct instance role and profile", func() {
			Expect(obj.Resources).ToNot(HaveKey("NodeInstanceRole"))
			Expect(obj.Resources).ToNot(HaveKey("NodeInstanceProfile"))

			Expect(getLaunchTemplateData(obj).IamInstanceProfile.Arn).To(Equal("arn:profile"))
		})
	})

	Context("NodeGroup{PrivateNetworking=true SSH.Allow=true}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.SSH.Allow = api.Enabled()
		keyName := ""
		ng.SSH.PublicKeyName = &keyName
		ng.InstanceType = "t2.medium"
		ng.PrivateNetworking = true
		ng.AMIFamily = "AmazonLinux2"

		build(cfg, "eksctl-test-private-ng", ng)

		roundtript()

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(obj.Description).To(ContainSubstring("SSH access: true"))
			Expect(obj.Description).To(ContainSubstring("private networking: true"))
		})

		It("should have correct resources and attributes", func() {
			Expect(obj.Resources).ToNot(BeEmpty())

			Expect(obj.Resources).To(HaveKey("NodeGroup"))
			ng := obj.Resources["NodeGroup"].Properties
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

			ltd := getLaunchTemplateData(obj)

			isFnGetAttOf(ltd.IamInstanceProfile.Arn, "NodeInstanceProfile.Arn")

			Expect(ltd.BlockDeviceMappings).To(HaveLen(1))

			rootVolume := ltd.BlockDeviceMappings[0].(map[string]interface{})

			Expect(rootVolume).To(HaveKeyWithValue("DeviceName", "/dev/xvda"))
			Expect(rootVolume).To(HaveKey("Ebs"))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("VolumeType", "io1"))
			Expect(rootVolume["Ebs"].(map[string]interface{})).To(HaveKeyWithValue("VolumeSize", 2.0))

			Expect(ltd.InstanceType).To(Equal("t2.medium"))

			Expect(ltd.NetworkInterfaces).To(HaveLen(1))
			Expect(ltd.NetworkInterfaces[0].DeviceIndex).To(Equal(0))
			Expect(ltd.NetworkInterfaces[0].AssociatePublicIpAddress).To(BeFalse())

			Expect(obj.Resources["SSHIPv4"].Properties.CidrIp).To(Equal("192.168.0.0/16"))
			Expect(obj.Resources["SSHIPv4"].Properties.FromPort).To(Equal(22))
			Expect(obj.Resources["SSHIPv4"].Properties.ToPort).To(Equal(22))

			Expect(obj.Resources).ToNot(HaveKey("SSHIPv6"))
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

		roundtript()

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(obj.Description).To(ContainSubstring("SSH access: true"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have correct resources and attributes", func() {
			Expect(obj.Resources).ToNot(BeEmpty())
			Expect(obj.Resources).To(HaveKey("NodeGroup"))

			Expect(obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier).ToNot(BeNil())
			x, ok := obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier.(map[string]interface{})
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

			ltd := getLaunchTemplateData(obj)

			Expect(ltd.InstanceType).To(Equal("t2.large"))

			Expect(ltd.NetworkInterfaces).To(HaveLen(1))
			Expect(ltd.NetworkInterfaces[0].DeviceIndex).To(Equal(0))
			Expect(ltd.NetworkInterfaces[0].AssociatePublicIpAddress).To(BeTrue())

			Expect(obj.Resources["SSHIPv4"].Properties.CidrIp).To(Equal("0.0.0.0/0"))
			Expect(obj.Resources["SSHIPv4"].Properties.FromPort).To(Equal(22))
			Expect(obj.Resources["SSHIPv4"].Properties.ToPort).To(Equal(22))

			Expect(obj.Resources["SSHIPv6"].Properties.CidrIpv6).To(Equal("::/0"))
			Expect(obj.Resources["SSHIPv6"].Properties.FromPort).To(Equal(22))
			Expect(obj.Resources["SSHIPv6"].Properties.ToPort).To(Equal(22))
		})
	})

	Context("NodeGroup{PrivateNetworking=false SSH.Allow=false}", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		cfg.VPC = &api.ClusterVPC{
			Network: api.Network{
				ID: vpcID,
			},
			SecurityGroup: "sg-0b44c48bcba5b7362",
			Subnets: &api.ClusterSubnets{
				Public: map[string]api.Network{
					"us-west-2b": {
						ID: "subnet-0f98135715dfcf55f",
					},
					"us-west-2a": {
						ID: "subnet-0ade11bad78dced9e",
					},
					"us-west-2c": {
						ID: "subnet-0e2e63ff1712bf6ef",
					},
				},
				Private: map[string]api.Network{
					"us-west-2b": {
						ID: "subnet-0f98135715dfcf55a",
					},
					"us-west-2a": {
						ID: "subnet-0ade11bad78dced9f",
					},
					"us-west-2c": {
						ID: "subnet-0e2e63ff1712bf6ea",
					},
				},
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

		roundtript()

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have correct resources and attributes", func() {
			Expect(obj.Resources).ToNot(BeEmpty())

			Expect(obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier).ToNot(BeNil())
			x, ok := obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier.([]interface{})
			Expect(ok).To(BeTrue())
			refSubnets := []interface{}{
				cfg.VPC.Subnets.Public["us-west-2a"].ID,
			}
			Expect(x).To((Equal(refSubnets)))

			ltd := getLaunchTemplateData(obj)

			Expect(ltd.InstanceType).To(Equal("t2.medium"))

			Expect(ltd.NetworkInterfaces).To(HaveLen(1))
			Expect(ltd.NetworkInterfaces[0].DeviceIndex).To(Equal(0))
			Expect(ltd.NetworkInterfaces[0].AssociatePublicIpAddress).To(BeTrue())

			Expect(obj.Resources).ToNot(HaveKey("SSHIPv4"))

			Expect(obj.Resources).ToNot(HaveKey("SSHIPv6"))

		})
	})

	checkAsset := func(name, expectedContent string) {
		assetContent, err := nodebootstrap.Asset(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(assetContent)).ToNot(BeEmpty())
		Expect(expectedContent).To(Equal(string(assetContent)))
	}

	getFile := func(c *cloudconfig.CloudConfig, p string) *cloudconfig.File {
		for _, f := range c.WriteFiles {
			if f.Path == p {
				return &f
			}
		}
		return nil
	}

	checkScript := func(c *cloudconfig.CloudConfig, p string, assetContent bool) {
		script := getFile(c, p)
		Expect(script).ToNot(BeNil())
		Expect(script.Permissions).To(Equal("0755"))
		scriptRuns := false
		for _, s := range c.Commands {
			if s.([]interface{})[0] == script.Path {
				scriptRuns = true
			}
		}
		Expect(scriptRuns).To(BeTrue())
		if assetContent {
			checkAsset(filepath.Base(p), script.Content)
		}
	}

	Context("UserData - AmazonLinux2", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "m5.large"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have correct instance type", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.large"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=",
				"NODE_TAINTS=",
			}))

			kubeletDropInUnit := getFile(cc, "/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf")
			Expect(kubeletDropInUnit).ToNot(BeNil())
			Expect(kubeletDropInUnit.Permissions).To(Equal("0644"))
			checkAsset("10-eksclt.al2.conf", kubeletDropInUnit.Content)

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("aws-iam-authenticator")))

			kubeletConfigAssetContent, err := nodebootstrap.Asset("kubelet.yaml")
			Expect(err).ToNot(HaveOccurred())

			kubeletConfigAssetContentString := string(kubeletConfigAssetContent) +
				"\n" +
				"maxPods: 29\n" +
				"clusterDNS: [10.100.0.10]\n"

			kubeletConfig := getFile(cc, "/etc/eksctl/kubelet.yaml")
			Expect(kubeletConfig).ToNot(BeNil())
			Expect(kubeletConfig.Permissions).To(Equal("0644"))

			Expect(kubeletConfig.Content).To(MatchYAML(kubeletConfigAssetContentString))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			checkScript(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.al2.sh", true)
		})
	})

	Context("UserData - AmazonLinux2 (custom pre-bootstrap)", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "m5.xlarge"

		ng.PreBootstrapCommands = []string{
			"touch /tmp/test",
			"rm /tmp/test",
		}

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have correct instance type", func() {
			Expect(getLaunchTemplateData(obj).InstanceType).To(Equal("m5.xlarge"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=",
				"NODE_TAINTS=",
			}))

			kubeletDropInUnit := getFile(cc, "/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf")
			Expect(kubeletDropInUnit).ToNot(BeNil())
			Expect(kubeletDropInUnit.Permissions).To(Equal("0644"))
			checkAsset("10-eksclt.al2.conf", kubeletDropInUnit.Content)

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("aws-iam-authenticator")))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			checkScript(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.al2.sh", true)

			Expect(cc.Commands).To(HaveLen(len(ng.PreBootstrapCommands) + 1))
			for i, cmd := range ng.PreBootstrapCommands {
				c := cc.Commands[i].([]interface{})
				Expect(c[0]).To(Equal("/bin/bash"))
				Expect(c[1]).To(Equal("-c"))
				Expect(c[2]).To(Equal(cmd))
			}
		})
	})

	Context("UserData - AmazonLinux2 (custom bootstrap)", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "m5.large"
		ng.Labels = map[string]string{
			"os": "al2",
		}
		ng.Taints = map[string]string{
			"key1": "value1:NoSchedule",
		}

		ng.OverrideBootstrapCommand = &overrideBootstrapCommand

		ng.ClusterDNS = "169.254.20.10"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=os=al2",
				"NODE_TAINTS=key1=value1:NoSchedule",
			}))

			kubeletDropInUnit := getFile(cc, "/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf")
			Expect(kubeletDropInUnit).ToNot(BeNil())
			Expect(kubeletDropInUnit.Permissions).To(Equal("0644"))
			checkAsset("10-eksclt.al2.conf", kubeletDropInUnit.Content)

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("aws-iam-authenticator")))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			script := getFile(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.al2.sh")
			Expect(script).To(BeNil())

			Expect(cc.Commands).To(HaveLen(1))
			Expect(cc.Commands[0]).To(HaveLen(3))
			c := cc.Commands[0].([]interface{})
			Expect(c[0]).To(Equal("/bin/bash"))
			Expect(c[1]).To(Equal("-c"))
			Expect(c[2]).To(Equal(overrideBootstrapCommand))
		})
	})

	Context("UserData - AmazonLinux2 (custom bootstrap and pre-bootstrap)", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		ng.InstanceType = "m5.large"
		ng.Labels = map[string]string{
			"os": "al2",
		}

		ng.PreBootstrapCommands = []string{"echo 1 > /tmp/1", "echo 2 > /tmp/2", "echo 3 > /tmp/3"}
		ng.OverrideBootstrapCommand = &overrideBootstrapCommand

		ng.ClusterDNS = "169.254.20.10"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=os=al2",
				"NODE_TAINTS=",
			}))

			kubeletDropInUnit := getFile(cc, "/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf")
			Expect(kubeletDropInUnit).ToNot(BeNil())
			Expect(kubeletDropInUnit.Permissions).To(Equal("0644"))
			checkAsset("10-eksclt.al2.conf", kubeletDropInUnit.Content)

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("aws-iam-authenticator")))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			script := getFile(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.al2.sh")
			Expect(script).To(BeNil())

			Expect(cc.Commands).To(HaveLen(4))
			Expect(cc.Commands[0]).To(HaveLen(3))

			for i, cmd := range ng.PreBootstrapCommands {
				c := cc.Commands[i].([]interface{})
				Expect(c[0]).To(Equal("/bin/bash"))
				Expect(c[1]).To(Equal("-c"))
				Expect(c[2]).To(Equal(cmd))
			}

			Expect(cc.Commands[3].([]interface{})[0]).To(Equal("/bin/bash"))
			Expect(cc.Commands[3].([]interface{})[1]).To(Equal("-c"))
			Expect(cc.Commands[3].([]interface{})[2]).To(Equal(overrideBootstrapCommand))
		})
	})

	Context("UserData - Ubuntu1804", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		cfg.VPC.CIDR, _ = ipnet.ParseCIDR("10.1.0.0/16")
		ng.AMIFamily = "Ubuntu1804"
		ng.InstanceType = "m5.large"

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: Ubuntu1804"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=",
				"NODE_TAINTS=",
				"MAX_PODS=29",
				"CLUSTER_DNS=172.20.0.10",
			}))

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("heptio-authenticator-aws")))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			checkScript(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.ubuntu.sh", true)
		})
	})

	Context("UserData - Ubuntu1804 (custom pre-bootstrap)", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		cfg.VPC.CIDR, _ = ipnet.ParseCIDR("10.1.0.0/16")
		ng.AMIFamily = "Ubuntu1804"
		ng.InstanceType = "m5.large"

		ng.PreBootstrapCommands = []string{
			"while true ; do echo foo > /dev/null ; done",
		}

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: Ubuntu1804"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=",
				"NODE_TAINTS=",
				"MAX_PODS=29",
				"CLUSTER_DNS=172.20.0.10",
			}))

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("heptio-authenticator-aws")))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			checkScript(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.ubuntu.sh", true)

			for i, cmd := range ng.PreBootstrapCommands {
				c := cc.Commands[i].([]interface{})
				Expect(c[0]).To(Equal("/bin/bash"))
				Expect(c[1]).To(Equal("-c"))
				Expect(c[2]).To(Equal(cmd))
			}
		})
	})

	Context("UserData - Ubuntu1804 (custom bootstrap)", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		cfg.VPC.CIDR, _ = ipnet.ParseCIDR("10.1.0.0/16")
		ng.AMIFamily = "Ubuntu1804"
		ng.InstanceType = "m5.large"

		ng.Labels = map[string]string{
			"os": "ubuntu",
		}
		ng.Taints = map[string]string{
			"key1": "value1:NoSchedule",
		}

		ng.ClusterDNS = "169.254.20.10"

		ng.OverrideBootstrapCommand = &overrideBootstrapCommand

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: Ubuntu1804"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=os=ubuntu",
				"NODE_TAINTS=key1=value1:NoSchedule",
				"MAX_PODS=29",
				"CLUSTER_DNS=169.254.20.10",
			}))

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("heptio-authenticator-aws")))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			script := getFile(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.ubuntu.sh")
			Expect(script).To(BeNil())

			Expect(cc.Commands).To(HaveLen(1))
			Expect(cc.Commands[0]).To(HaveLen(3))
			Expect(cc.Commands[0].([]interface{})[0]).To(Equal("/bin/bash"))
			Expect(cc.Commands[0].([]interface{})[1]).To(Equal("-c"))
			Expect(cc.Commands[0].([]interface{})[2]).To(Equal(overrideBootstrapCommand))
		})
	})

	Context("UserData - Ubuntu1804 (custom bootstrap and pre-bootstrap)", func() {
		cfg, ng := newClusterConfigAndNodegroup(true)

		cfg.VPC.CIDR, _ = ipnet.ParseCIDR("10.1.0.0/16")
		ng.AMIFamily = "Ubuntu1804"
		ng.InstanceType = "m5.large"

		ng.Labels = map[string]string{
			"os": "ubuntu",
		}

		ng.ClusterDNS = "169.254.20.10"

		ng.PreBootstrapCommands = []string{"echo 1 > /tmp/1", "echo 2 > /tmp/2", "echo 3 > /tmp/3"}
		ng.OverrideBootstrapCommand = &overrideBootstrapCommand

		build(cfg, "eksctl-test-123-cluster", ng)

		roundtript()

		extractCloudConfig()

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: Ubuntu1804"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(cc.Packages).Should(BeEmpty())

			kubeletEnv := getFile(cc, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"NODE_LABELS=os=ubuntu",
				"NODE_TAINTS=",
				"MAX_PODS=29",
				"CLUSTER_DNS=169.254.20.10",
			}))

			kubeconfig := getFile(cc, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(MatchYAML(kubeconfigBody("heptio-authenticator-aws")))

			ca := getFile(cc, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			script := getFile(cc, "/var/lib/cloud/scripts/per-instance/bootstrap.ubuntu.sh")
			Expect(script).To(BeNil())

			Expect(cc.Commands).To(HaveLen(4))
			Expect(cc.Commands[0]).To(HaveLen(3))

			for i, cmd := range ng.PreBootstrapCommands {
				c := cc.Commands[i].([]interface{})
				Expect(c[0]).To(Equal("/bin/bash"))
				Expect(c[1]).To(Equal("-c"))
				Expect(c[2]).To(Equal(cmd))
			}

			c3 := cc.Commands[3].([]interface{})
			Expect(c3[0]).To(Equal("/bin/bash"))
			Expect(c3[1]).To(Equal("-c"))
			Expect(c3[2]).To(Equal(overrideBootstrapCommand))
		})
	})
})

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
	Expect(o).To(HaveKey(gfn.Ref))
	Expect(o[gfn.Ref]).To(Equal(value))
}

func isFnGetAttOf(obj interface{}, value string) {
	Expect(obj).ToNot(BeEmpty())
	o, ok := obj.(map[string]interface{})
	Expect(ok).To(BeTrue())
	Expect(o).To(HaveKey(gfn.FnGetAtt))
	Expect(o[gfn.FnGetAtt]).To(Equal(value))
}

func getLaunchTemplateData(obj *Template) LaunchTemplateData {
	Expect(obj.Resources).ToNot(BeEmpty())
	Expect(obj.Resources).To(HaveKey("NodeGroupLaunchTemplate"))
	return obj.Resources["NodeGroupLaunchTemplate"].Properties.LaunchTemplateData
}
