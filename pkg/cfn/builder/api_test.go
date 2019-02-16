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
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	. "github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	totalNodeResources = 10
	clusterName        = "ferocious-mushroom-1532594698"
	endpoint           = "https://DE37D8AFB23F7275D2361AD6B2599143.yl4.us-west-2.eks.amazonaws.com"
	caCert             = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNE1EWXdOekExTlRBMU5Wb1hEVEk0TURZd05EQTFOVEExTlZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTWJoCnpvZElYR0drckNSZE1jUmVEN0YvMnB1NFZweTdvd3FEVDgrdk9zeGs2bXFMNWxQd3ZicFhmYkE3R0xzMDVHa0wKaDdqL0ZjcU91cnMwUFZSK3N5REtuQXltdDFORWxGNllGQktSV1dUQ1hNd2lwN1pweW9XMXdoYTlJYUlPUGxCTQpPTEVlckRabFVrVDFVV0dWeVdsMmxPeFgxa2JhV2gvakptWWdkeW5jMXhZZ3kxa2JybmVMSkkwLzVUVTRCajJxClB1emtrYW5Xd3lKbGdXQzhBSXlpWW82WFh2UVZmRzYrM3RISE5XM1F1b3ZoRng2MTFOYnl6RUI3QTdtZGNiNmgKR0ZpWjdOeThHZnFzdjJJSmI2Nk9FVzBSdW9oY1k3UDZPdnZmYnlKREhaU2hqTStRWFkxQXN5b3g4Ri9UelhHSgpQUWpoWUZWWEVhZU1wQmJqNmNFQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFCa2hKRVd4MHk1LzlMSklWdXJ1c1hZbjN6Z2EKRkZ6V0JsQU44WTlqUHB3S2t0Vy9JNFYyUGg3bWY2Z3ZwZ3Jhc2t1Slk1aHZPcDdBQmcxSTFhaHUxNUFpMUI0ZApuMllRaDlOaHdXM2pKMmhuRXk0VElpb0gza2JFdHRnUVB2bWhUQzNEYUJreEpkbmZJSEJCV1RFTTU1czRwRmxUClpzQVJ3aDc1Q3hYbjdScVU0akpKcWNPaTRjeU5qeFVpRDBqR1FaTmNiZWEyMkRCeTJXaEEzUWZnbGNScGtDVGUKRDVPS3NOWlF4MW9MZFAwci9TSmtPT1NPeUdnbVJURTIrODQxN21PRW02Z3RPMCszdWJkbXQ0aENsWEtFTTZYdwpuQWNlK0JxVUNYblVIN2ZNS3p2TDE5UExvMm5KbFU1TnlCbU1nL1pNVHVlUy80eFZmKy94WnpsQ0Q1WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	arn                = "arn:aws:eks:us-west-2:122333:cluster/" + clusterName

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

type Template struct {
	Description string
	Resources   map[string]struct {
		Properties struct {
			Tags           []Tag
			UserData       string
			PolicyDocument struct {
				Statement []struct {
					Action   []string
					Effect   string
					Resource interface{}
				}
			}
			VPCZoneIdentifier        interface{}
			AssociatePublicIpAddress bool
			CidrIp                   string
			CidrIpv6                 string
			IpProtocol               string
			FromPort                 int
			ToPort                   int
		}
	}
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
	var caCertData []byte

	Describe("should decode CA data", func() {
		var err error
		caCertData, err = base64.StdEncoding.DecodeString(caCert)
		It("should not error", func() { Expect(err).ShouldNot(HaveOccurred()) })
	})

	testAZs := []string{"us-west-2b", "us-west-2a", "us-west-2c"}

	newClusterConfigAndNodegroup := func() (*api.ClusterConfig, *api.NodeGroup) {
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
		ng.VolumeSize = 2
		ng.VolumeType = api.NodeVolumeTypeIO1

		*cfg.VPC.CIDR = api.DefaultCIDR()

		return cfg, ng
	}

	newClusterConfig := func() *api.ClusterConfig {
		cfg, _ := newClusterConfigAndNodegroup()
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
				Version: "1.11",
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
					SecurityGroups: &api.NodeGroupSGs{
						WithLocal:  api.NewBoolTrue(),
						WithShared: api.NewBoolTrue(),
						AttachIDs:  []string{},
					},
					DesiredCapacity: nil,
					VolumeSize:      2,
					VolumeType:      api.NodeVolumeTypeIO1,
					IAM: &api.NodeGroupIAM{
						WithAddonPolicies: api.NodeGroupIAMAddonPolicies{
							ImageBuilder: api.NewBoolFalse(),
							AutoScaler:   api.NewBoolFalse(),
							ExternalDNS:  api.NewBoolFalse(),
						},
					},
				},
			},
		}

		cfg := newClusterConfig()

		It("should not error when calling SetSubnets", func() {
			err := vpc.SetSubnets(cfg)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have public and private subnets", func() {
			Expect(len(cfg.VPC.Subnets.Private)).To(Equal(3))
			Expect(len(cfg.VPC.Subnets.Public)).To(Equal(3))
		})

		rs := NewClusterResourceSet(p, cfg)
		It("should add all resources without error", func() {
			err := rs.AddAllResources()
			Expect(err).ShouldNot(HaveOccurred())
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

		sampleStack := newStackWithOutputs(sampleOutputs)

		It("should not error", func() {
			err := rs.GetAllOutputs(sampleStack)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be equal", func() {
			cfgData, err := json.Marshal(cfg)
			Expect(err).ShouldNot(HaveOccurred())
			expectedData, err := json.Marshal(expected)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfgData).To(MatchJSON(expectedData))
		})
	})

	Describe("AutoNameTag", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-123-cluster", ng)

		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
			t := rs.Template()
			Expect(t.Resources).ToNot(BeEmpty())
			Expect(len(t.Resources)).To(Equal(totalNodeResources))
			templateBody, err := t.JSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(templateBody).ShouldNot(BeEmpty())
		})

		templateBody, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(templateBody).ShouldNot(BeEmpty())
		})
		obj := Template{}
		It("should parse JSON without errors", func() {
			err := json.Unmarshal(templateBody, &obj)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("SG should have correct tags", func() {
			Expect(obj.Resources).ToNot(BeNil())
			Expect(len(obj.Resources)).To(Equal(totalNodeResources))
			Expect(len(obj.Resources["SG"].Properties.Tags)).To(Equal(2))
			Expect(obj.Resources["SG"].Properties.Tags[0].Key).To(Equal("kubernetes.io/cluster/" + clusterName))
			Expect(obj.Resources["SG"].Properties.Tags[0].Value).To(Equal("owned"))
			Expect(obj.Resources["SG"].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(obj.Resources["SG"].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SG",
			}))
		})
	})

	Describe("NodeGroupTags", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		ng.InstanceType = "t2.medium"
		ng.Name = "ng-abcd1234"

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-123-cluster", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors", func() {
			err := json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have correct tags", func() {
			Expect(len(obj.Resources)).ToNot(Equal(0))
			Expect(len(obj.Resources["NodeGroup"].Properties.Tags)).To(Equal(2))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[0].Key).To(Equal("Name"))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[0].Value).To(Equal(clusterName + "-ng-abcd1234-Node"))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[0].PropagateAtLaunch).To(Equal("true"))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[1].Key).To(Equal("kubernetes.io/cluster/" + clusterName))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[1].Value).To(Equal("owned"))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[1].PropagateAtLaunch).To(Equal("true"))
		})
	})

	Describe("NodeGroupAutoScaling", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		ng.IAM.WithAddonPolicies.AutoScaler = api.NewBoolTrue()

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-123-cluster", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors", func() {
			err := json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have correct policies", func() {
			Expect(len(obj.Resources)).ToNot(Equal(0))
			Expect(obj.Resources["PolicyAutoScaling"]).ToNot(BeNil())
			Expect(len(obj.Resources["PolicyAutoScaling"].Properties.PolicyDocument.Statement)).To(Equal(1))
			Expect(obj.Resources["PolicyAutoScaling"].Properties.PolicyDocument.Statement[0].Effect).To(Equal("Allow"))
			Expect(obj.Resources["PolicyAutoScaling"].Properties.PolicyDocument.Statement[0].Resource).To(Equal("*"))
			Expect(obj.Resources["PolicyAutoScaling"].Properties.PolicyDocument.Statement[0].Action).To(Equal([]string{
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

			Expect(obj.Resources).To(HaveKey("NodeGroup"))
			ng := obj.Resources["NodeGroup"]
			Expect(ng).ToNot(BeNil())
			Expect(ng.Properties).ToNot(BeNil())
			Expect(ng.Properties.Tags).ToNot(BeNil())
			Expect(ng.Properties.Tags).To(Equal(expectedTags))
		})
	})

	Describe("NodeGroup{PrivateNetworking=true AllowSSH=true}", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		ng.AllowSSH = true
		ng.InstanceType = "t2.medium"
		ng.PrivateNetworking = true
		ng.AMIFamily = "AmazonLinux2"

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-private-ng", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors", func() {
			err := json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(obj.Description).To(ContainSubstring("SSH access: true"))
			Expect(obj.Description).To(ContainSubstring("private networking: true"))
		})

		It("should have correct resources and attributes", func() {
			Expect(len(obj.Resources)).ToNot(Equal(0))

			Expect(obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier).To(Not(BeNil()))
			x, ok := obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(len(x)).To(Equal(1))
			refSubnets := map[string]interface{}{
				"Fn::Split": []interface{}{
					",",
					map[string]interface{}{
						"Fn::ImportValue": "eksctl-test-private-ng::SubnetsPrivate",
					},
				},
			}
			Expect(x).To(Equal(refSubnets))

			Expect(obj.Resources["NodeLaunchConfig"].Properties.AssociatePublicIpAddress).To(BeFalse())

			Expect(obj.Resources["SSHIPv4"].Properties.CidrIp).To(Equal("192.168.0.0/16"))
			Expect(obj.Resources["SSHIPv4"].Properties.FromPort).To(Equal(22))
			Expect(obj.Resources["SSHIPv4"].Properties.ToPort).To(Equal(22))

			Expect(obj.Resources).To(Not(HaveKey("SSHIPv6")))
		})
	})

	Describe("NodeGroup{PrivateNetworking=false AllowSSH=true}", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		ng.AllowSSH = true
		ng.InstanceType = "t2.medium"
		ng.PrivateNetworking = false
		ng.AMIFamily = "AmazonLinux2"

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-public-ng", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		obj := Template{}
		It("should parse JSON without errors", func() {
			err := json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(obj.Description).To(ContainSubstring("SSH access: true"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have correct resources and attributes", func() {
			Expect(len(obj.Resources)).ToNot(Equal(0))

			Expect(obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier).To(Not(BeNil()))
			x, ok := obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(len(x)).To(Equal(1))
			refSubnets := map[string]interface{}{
				"Fn::Split": []interface{}{
					",",
					map[string]interface{}{
						"Fn::ImportValue": "eksctl-test-public-ng::SubnetsPublic",
					},
				},
			}
			Expect(x).To(Equal(refSubnets))

			Expect(obj.Resources["NodeLaunchConfig"].Properties.AssociatePublicIpAddress).To(BeTrue())

			Expect(obj.Resources["SSHIPv4"].Properties.CidrIp).To(Equal("0.0.0.0/0"))
			Expect(obj.Resources["SSHIPv4"].Properties.FromPort).To(Equal(22))
			Expect(obj.Resources["SSHIPv4"].Properties.ToPort).To(Equal(22))

			Expect(obj.Resources["SSHIPv6"].Properties.CidrIpv6).To(Equal("::/0"))
			Expect(obj.Resources["SSHIPv6"].Properties.FromPort).To(Equal(22))
			Expect(obj.Resources["SSHIPv6"].Properties.ToPort).To(Equal(22))
		})
	})

	Describe("NodeGroup{PrivateNetworking=false AllowSSH=false}", func() {
		cfg, ng := newClusterConfigAndNodegroup()

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
		ng.AllowSSH = false
		ng.InstanceType = "t2.medium"
		ng.PrivateNetworking = false
		ng.AMIFamily = "AmazonLinux2"

		It("should have 1 AZ for the nodegroup", func() {
			Expect(ng.AvailabilityZones).To(Equal([]string{"us-west-2a"}))
		})

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-public-ng", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors", func() {
			err := json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: AmazonLinux2"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have correct resources and attributes", func() {
			Expect(len(obj.Resources)).ToNot(Equal(0))

			Expect(obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier).To(Not(BeNil()))
			x, ok := obj.Resources["NodeGroup"].Properties.VPCZoneIdentifier.([]interface{})
			Expect(ok).To(BeTrue())
			refSubnets := []interface{}{
				cfg.VPC.Subnets.Public["us-west-2a"].ID,
			}
			Expect(x).To((Equal(refSubnets)))

			Expect(obj.Resources["NodeLaunchConfig"].Properties.AssociatePublicIpAddress).To(BeTrue())

			Expect(obj.Resources).To(Not(HaveKey("SSHIPv4")))

			Expect(obj.Resources).To(Not(HaveKey("SSHIPv6")))

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

	Describe("UserData - AmazonLinux2", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		var c *cloudconfig.CloudConfig

		cfg.NodeGroups[0].InstanceType = "m5.large"

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-123-cluster", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors and extract valid cloud-config using our implementation", func() {
			err = json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(obj.Resources)).ToNot(Equal(0))

			userData := obj.Resources["NodeLaunchConfig"].Properties.UserData
			Expect(userData).ToNot(BeEmpty())

			c, err = cloudconfig.DecodeCloudConfig(userData)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(c).ToNot(BeNil())

			Expect(c.Packages).Should(BeEmpty())

			kubeletEnv := getFile(c, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"MAX_PODS=29",
				"CLUSTER_DNS=10.100.0.10",
				"NODE_LABELS=",
			}))

			kubeletDropInUnit := getFile(c, "/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf")
			Expect(kubeletDropInUnit).ToNot(BeNil())
			Expect(kubeletDropInUnit.Permissions).To(Equal("0644"))
			checkAsset("10-eksclt.al2.conf", kubeletDropInUnit.Content)

			kubeconfig := getFile(c, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(Equal(kubeconfigBody("aws-iam-authenticator")))

			ca := getFile(c, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			checkScript(c, "/var/lib/cloud/scripts/per-instance/bootstrap.al2.sh", true)
		})
	})

	Describe("UserData - AmazonLinux2 (custom)", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		var c *cloudconfig.CloudConfig

		cfg.NodeGroups[0].InstanceType = "m5.large"
		cfg.NodeGroups[0].Labels = map[string]string{
			"os": "al2",
		}

		cfg.NodeGroups[0].OverrideBootstrapCommand = &overrideBootstrapCommand

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-123-cluster", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors and extract valid cloud-config using our implementation", func() {
			err = json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(obj.Resources)).ToNot(Equal(0))

			userData := obj.Resources["NodeLaunchConfig"].Properties.UserData
			Expect(userData).ToNot(BeEmpty())

			c, err = cloudconfig.DecodeCloudConfig(userData)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(c).ToNot(BeNil())

			Expect(c.Packages).Should(BeEmpty())

			kubeletEnv := getFile(c, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"MAX_PODS=29",
				"CLUSTER_DNS=10.100.0.10",
				"NODE_LABELS=os=al2",
			}))

			kubeletDropInUnit := getFile(c, "/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf")
			Expect(kubeletDropInUnit).ToNot(BeNil())
			Expect(kubeletDropInUnit.Permissions).To(Equal("0644"))
			checkAsset("10-eksclt.al2.conf", kubeletDropInUnit.Content)

			kubeconfig := getFile(c, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(Equal(kubeconfigBody("aws-iam-authenticator")))

			ca := getFile(c, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			script := getFile(c, "/var/lib/cloud/scripts/per-instance/bootstrap.al2.sh")
			Expect(script).To(BeNil())

			Expect(c.Commands).To(HaveLen(1))
			Expect(c.Commands[0]).To(HaveLen(3))
			Expect(c.Commands[0].([]interface{})[0]).To(Equal("/bin/bash"))
			Expect(c.Commands[0].([]interface{})[1]).To(Equal("-c"))
			Expect(c.Commands[0].([]interface{})[2]).To(Equal(overrideBootstrapCommand))
		})
	})

	Describe("UserData - Ubuntu1804", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		var c *cloudconfig.CloudConfig

		cfg.VPC.CIDR, _ = ipnet.ParseCIDR("10.1.0.0/16")
		cfg.NodeGroups[0].AMIFamily = "Ubuntu1804"
		cfg.NodeGroups[0].InstanceType = "m5.large"

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-123-cluster", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors and extract valid cloud-config using our implementation", func() {
			err = json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(obj.Resources)).ToNot(Equal(0))

			userData := obj.Resources["NodeLaunchConfig"].Properties.UserData
			Expect(userData).ToNot(BeEmpty())

			c, err = cloudconfig.DecodeCloudConfig(userData)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: Ubuntu1804"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(c).ToNot(BeNil())

			Expect(c.Packages).Should(BeEmpty())

			kubeletEnv := getFile(c, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"MAX_PODS=29",
				"CLUSTER_DNS=172.20.0.10",
				"NODE_LABELS=",
			}))

			kubeconfig := getFile(c, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(Equal(kubeconfigBody("heptio-authenticator-aws")))

			ca := getFile(c, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			checkScript(c, "/var/lib/cloud/scripts/per-instance/bootstrap.ubuntu.sh", true)
		})
	})

	Describe("UserData - Ubuntu1804 (custom)", func() {
		cfg, ng := newClusterConfigAndNodegroup()

		var c *cloudconfig.CloudConfig

		cfg.VPC.CIDR, _ = ipnet.ParseCIDR("10.1.0.0/16")
		cfg.NodeGroups[0].AMIFamily = "Ubuntu1804"
		cfg.NodeGroups[0].InstanceType = "m5.large"

		cfg.NodeGroups[0].Labels = map[string]string{
			"os": "ubuntu",
		}

		cfg.NodeGroups[0].OverrideBootstrapCommand = &overrideBootstrapCommand

		rs := NewNodeGroupResourceSet(p, cfg, "eksctl-test-123-cluster", ng)
		err := rs.AddAllResources()
		It("should add all resources without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON without errors and extract valid cloud-config using our implementation", func() {
			err = json.Unmarshal(template, &obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(obj.Resources)).ToNot(Equal(0))

			userData := obj.Resources["NodeLaunchConfig"].Properties.UserData
			Expect(userData).ToNot(BeEmpty())

			c, err = cloudconfig.DecodeCloudConfig(userData)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have correct description", func() {
			Expect(obj.Description).To(ContainSubstring("AMI family: Ubuntu1804"))
			Expect(obj.Description).To(ContainSubstring("SSH access: false"))
			Expect(obj.Description).To(ContainSubstring("private networking: false"))
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(c).ToNot(BeNil())

			Expect(c.Packages).Should(BeEmpty())

			kubeletEnv := getFile(c, "/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"MAX_PODS=29",
				"CLUSTER_DNS=172.20.0.10",
				"NODE_LABELS=os=ubuntu",
			}))

			kubeconfig := getFile(c, "/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			Expect(kubeconfig.Content).To(Equal(kubeconfigBody("heptio-authenticator-aws")))

			ca := getFile(c, "/etc/eksctl/ca.crt")
			Expect(ca).ToNot(BeNil())
			Expect(ca.Permissions).To(Equal("0644"))
			Expect(ca.Content).To(Equal(string(caCertData)))

			script := getFile(c, "/var/lib/cloud/scripts/per-instance/bootstrap.ubuntu.sh")
			Expect(script).To(BeNil())

			Expect(c.Commands).To(HaveLen(1))
			Expect(c.Commands[0]).To(HaveLen(3))
			Expect(c.Commands[0].([]interface{})[0]).To(Equal("/bin/bash"))
			Expect(c.Commands[0].([]interface{})[1]).To(Equal("-c"))
			Expect(c.Commands[0].([]interface{})[2]).To(Equal(overrideBootstrapCommand))
		})
	})
})
