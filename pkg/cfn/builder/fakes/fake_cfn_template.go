package fakes

import (
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

type FakeTemplate struct {
	Description string
	Resources   map[string]struct {
		Type         string
		Properties   Properties
		DependsOn    []string
		UpdatePolicy map[string]map[string]interface{}
	}
	Mappings map[string]interface{}
	Outputs  map[string]cfn.Output
}

type Tag struct {
	Key   interface{}
	Value interface{}

	PropagateAtLaunch string
}

type Properties struct {
	GroupDescription           string
	Description                string
	Tags                       []Tag
	SecurityGroupIngress       []SGIngress
	GroupID                    interface{}
	SourceSecurityGroupID      interface{}
	DestinationSecurityGroupID interface{}

	Path, RoleName           string
	Roles, ManagedPolicyArns []interface{}
	PermissionsBoundary      interface{}
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
	LaunchTemplateName interface{}
	Strategy           string

	CapacityRebalance bool

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
	MapPublicIPOnLaunch                        bool

	Ipv6CidrBlock map[string][]interface{}

	AmazonProvidedIpv6CidrBlock         bool
	AvailabilityZone, Domain, CidrBlock string

	Name, Version      string
	RoleArn            interface{}
	ResourcesVpcConfig struct {
		SecurityGroupIds []interface{}
		SubnetIds        []interface{}
	}
	EncryptionConfig []struct {
		Provider struct {
			KeyARN interface{}
		}
		Resources []string
	}
	LaunchTemplate struct {
		LaunchTemplateName map[string]interface{}
		Version            map[string]interface{}
		Overrides          []struct {
			InstanceType string
		}
	}
	MixedInstancesPolicy *struct {
		LaunchTemplate struct {
			LaunchTemplateSpecification struct {
				LaunchTemplateName map[string]interface{}
				Version            map[string]interface{}
			}
			Overrides []struct {
				InstanceType string
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

type SGIngress struct {
	SourceSecurityGroupID interface{}
	FromPort              float64
	ToPort                float64
	Description           string
	IPProtocol            string
}

type LaunchTemplateData struct {
	IamInstanceProfile              struct{ Arn interface{} }
	UserData, InstanceType, ImageID string
	BlockDeviceMappings             []BlockDeviceMappings
	EbsOptimized                    *bool
	Monitoring                      *Monitoring
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
	MetadataOptions   MetadataOptions
	TagSpecifications []TagSpecification
	Placement         Placement
	KeyName           string
}

type Placement struct {
	GroupName interface{}
}

type BlockDeviceMappings struct {
	DeviceName string
	Ebs        map[string]interface{}
}

type MetadataOptions struct {
	HTTPPutResponseHopLimit float64
	HTTPTokens              string
}

type TagSpecification struct {
	ResourceType *string
	Tags         []Tag
}

type NetworkInterface struct {
	DeviceIndex              int
	AssociatePublicIPAddress bool
	NetworkCardIndex         int
	InterfaceType            string
}

type Monitoring struct {
	Enabled bool
}
