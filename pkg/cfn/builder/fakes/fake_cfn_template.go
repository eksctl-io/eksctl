package fakes

type FakeTemplate struct {
	Description string
	Resources   map[string]struct {
		Type         string
		Properties   Properties
		DependsOn    []string
		UpdatePolicy map[string]map[string]interface{}
	}
	Mappings map[string]interface{}
	Outputs  interface{}
}

type Tag struct {
	Key   interface{}
	Value interface{}

	PropagateAtLaunch string
}

type Properties struct {
	AcceptRoleSessionName                bool
	EnableDNSHostnames, EnableDNSSupport bool
	GroupDescription                     string
	Description                          string
	Tags                                 []Tag
	SecurityGroupIngress                 []SGIngress
	BootstrapSelfManagedAddons           bool
	GroupID                              interface{}
	SourceSecurityGroupID                interface{}
	DestinationSecurityGroupID           interface{}

	Type                     string
	Path, RoleName           string
	Roles, ManagedPolicyArns []interface{}
	PermissionsBoundary      interface{}
	AssumeRolePolicyDocument interface{}

	PolicyDocument struct {
		Version   string
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
	MaxInstanceLifetime               int

	CidrIP, CidrIPv6, IPProtocol string
	FromPort, ToPort             int

	VpcID, SubnetID                                                            interface{}
	EgressOnlyInternetGatewayID, RouteTableID, AllocationID                    interface{}
	GatewayID, InternetGatewayID, NatGatewayID, VpnGatewayId, TransitGatewayId interface{}
	DestinationCidrBlock, DestinationIpv6CidrBlock                             interface{}
	MapPublicIPOnLaunch                                                        bool
	AssignIpv6AddressOnCreation                                                *bool

	Ipv6CidrBlock           interface{}
	Ipv6Pool                string
	CidrBlock               interface{}
	KubernetesNetworkConfig KubernetesNetworkConfig

	AmazonProvidedIpv6CidrBlock bool
	AvailabilityZone, Domain    string

	Name               interface{}
	Version            string
	RoleArn            interface{}
	ResourcesVpcConfig struct {
		SecurityGroupIDs      []interface{}
		SubnetIDs             []interface{}
		EndpointPublicAccess  bool
		EndpointPrivateAccess bool
		PublicAccessCidrs     []string
	}
	EncryptionConfig []struct {
		Provider struct {
			KeyARN interface{}
		}
		Resources []string
	}
	AccessConfig struct {
		AuthenticationMode                      string
		BootstrapClusterCreatorAdminPermissions bool
	}
	LaunchTemplate struct {
		LaunchTemplateName map[string]interface{}
		Version            map[string]interface{}
		Overrides          []struct {
			InstanceType string
		}
	}
	Logging              ClusterLogging
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

type KubernetesNetworkConfig struct {
	ServiceIPv4CIDR string
	ServiceIPv6CIDR interface{}
	IPFamily        string
}

type ClusterLogging struct {
	ClusterLogging struct {
		EnabledTypes []ClusterLoggingType
	}
}

type ClusterLoggingType struct {
	Type string
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
	MetadataOptions                  MetadataOptions
	TagSpecifications                []TagSpecification
	Placement                        Placement
	KeyName                          string
	CapacityReservationSpecification *CapacityReservationSpecification
}

type CapacityReservationSpecification struct {
	CapacityReservationPreference *string
	CapacityReservationTarget     *CapacityReservationTarget
}

type CapacityReservationTarget struct {
	CapacityReservationID               *string
	CapacityReservationResourceGroupARN *string
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
