package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// NetworkInsightsAnalysis_Explanation AWS CloudFormation Resource (AWS::EC2::NetworkInsightsAnalysis.Explanation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html
type NetworkInsightsAnalysis_Explanation struct {

	// Acl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-acl
	Acl *NetworkInsightsAnalysis_AnalysisComponent `json:"Acl,omitempty"`

	// AclRule AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-aclrule
	AclRule *NetworkInsightsAnalysis_AnalysisAclRule `json:"AclRule,omitempty"`

	// Address AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-address
	Address *types.Value `json:"Address,omitempty"`

	// Addresses AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-addresses
	Addresses *types.Value `json:"Addresses,omitempty"`

	// AttachedTo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-attachedto
	AttachedTo *NetworkInsightsAnalysis_AnalysisComponent `json:"AttachedTo,omitempty"`

	// AvailabilityZones AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-availabilityzones
	AvailabilityZones *types.Value `json:"AvailabilityZones,omitempty"`

	// Cidrs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-cidrs
	Cidrs *types.Value `json:"Cidrs,omitempty"`

	// ClassicLoadBalancerListener AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-classicloadbalancerlistener
	ClassicLoadBalancerListener *NetworkInsightsAnalysis_AnalysisLoadBalancerListener `json:"ClassicLoadBalancerListener,omitempty"`

	// Component AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-component
	Component *NetworkInsightsAnalysis_AnalysisComponent `json:"Component,omitempty"`

	// CustomerGateway AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-customergateway
	CustomerGateway *NetworkInsightsAnalysis_AnalysisComponent `json:"CustomerGateway,omitempty"`

	// Destination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-destination
	Destination *NetworkInsightsAnalysis_AnalysisComponent `json:"Destination,omitempty"`

	// DestinationVpc AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-destinationvpc
	DestinationVpc *NetworkInsightsAnalysis_AnalysisComponent `json:"DestinationVpc,omitempty"`

	// Direction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-direction
	Direction *types.Value `json:"Direction,omitempty"`

	// ElasticLoadBalancerListener AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-elasticloadbalancerlistener
	ElasticLoadBalancerListener *NetworkInsightsAnalysis_AnalysisComponent `json:"ElasticLoadBalancerListener,omitempty"`

	// ExplanationCode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-explanationcode
	ExplanationCode *types.Value `json:"ExplanationCode,omitempty"`

	// IngressRouteTable AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-ingressroutetable
	IngressRouteTable *NetworkInsightsAnalysis_AnalysisComponent `json:"IngressRouteTable,omitempty"`

	// InternetGateway AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-internetgateway
	InternetGateway *NetworkInsightsAnalysis_AnalysisComponent `json:"InternetGateway,omitempty"`

	// LoadBalancerArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-loadbalancerarn
	LoadBalancerArn *types.Value `json:"LoadBalancerArn,omitempty"`

	// LoadBalancerListenerPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-loadbalancerlistenerport
	LoadBalancerListenerPort *types.Value `json:"LoadBalancerListenerPort,omitempty"`

	// LoadBalancerTarget AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-loadbalancertarget
	LoadBalancerTarget *NetworkInsightsAnalysis_AnalysisLoadBalancerTarget `json:"LoadBalancerTarget,omitempty"`

	// LoadBalancerTargetGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-loadbalancertargetgroup
	LoadBalancerTargetGroup *NetworkInsightsAnalysis_AnalysisComponent `json:"LoadBalancerTargetGroup,omitempty"`

	// LoadBalancerTargetGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-loadbalancertargetgroups
	LoadBalancerTargetGroups []NetworkInsightsAnalysis_AnalysisComponent `json:"LoadBalancerTargetGroups,omitempty"`

	// LoadBalancerTargetPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-loadbalancertargetport
	LoadBalancerTargetPort *types.Value `json:"LoadBalancerTargetPort,omitempty"`

	// MissingComponent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-missingcomponent
	MissingComponent *types.Value `json:"MissingComponent,omitempty"`

	// NatGateway AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-natgateway
	NatGateway *NetworkInsightsAnalysis_AnalysisComponent `json:"NatGateway,omitempty"`

	// NetworkInterface AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-networkinterface
	NetworkInterface *NetworkInsightsAnalysis_AnalysisComponent `json:"NetworkInterface,omitempty"`

	// PacketField AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-packetfield
	PacketField *types.Value `json:"PacketField,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-port
	Port *types.Value `json:"Port,omitempty"`

	// PortRanges AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-portranges
	PortRanges []NetworkInsightsAnalysis_PortRange `json:"PortRanges,omitempty"`

	// PrefixList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-prefixlist
	PrefixList *NetworkInsightsAnalysis_AnalysisComponent `json:"PrefixList,omitempty"`

	// Protocols AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-protocols
	Protocols *types.Value `json:"Protocols,omitempty"`

	// RouteTable AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-routetable
	RouteTable *NetworkInsightsAnalysis_AnalysisComponent `json:"RouteTable,omitempty"`

	// RouteTableRoute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-routetableroute
	RouteTableRoute *NetworkInsightsAnalysis_AnalysisRouteTableRoute `json:"RouteTableRoute,omitempty"`

	// SecurityGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-securitygroup
	SecurityGroup *NetworkInsightsAnalysis_AnalysisComponent `json:"SecurityGroup,omitempty"`

	// SecurityGroupRule AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-securitygrouprule
	SecurityGroupRule *NetworkInsightsAnalysis_AnalysisSecurityGroupRule `json:"SecurityGroupRule,omitempty"`

	// SecurityGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-securitygroups
	SecurityGroups []NetworkInsightsAnalysis_AnalysisComponent `json:"SecurityGroups,omitempty"`

	// SourceVpc AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-sourcevpc
	SourceVpc *NetworkInsightsAnalysis_AnalysisComponent `json:"SourceVpc,omitempty"`

	// State AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-state
	State *types.Value `json:"State,omitempty"`

	// Subnet AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-subnet
	Subnet *NetworkInsightsAnalysis_AnalysisComponent `json:"Subnet,omitempty"`

	// SubnetRouteTable AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-subnetroutetable
	SubnetRouteTable *NetworkInsightsAnalysis_AnalysisComponent `json:"SubnetRouteTable,omitempty"`

	// Vpc AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-vpc
	Vpc *NetworkInsightsAnalysis_AnalysisComponent `json:"Vpc,omitempty"`

	// VpcPeeringConnection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-vpcpeeringconnection
	VpcPeeringConnection *NetworkInsightsAnalysis_AnalysisComponent `json:"VpcPeeringConnection,omitempty"`

	// VpnConnection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-vpnconnection
	VpnConnection *NetworkInsightsAnalysis_AnalysisComponent `json:"VpnConnection,omitempty"`

	// VpnGateway AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-vpngateway
	VpnGateway *NetworkInsightsAnalysis_AnalysisComponent `json:"VpnGateway,omitempty"`

	// vpcEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-explanation.html#cfn-ec2-networkinsightsanalysis-explanation-vpcendpoint
	vpcEndpoint *NetworkInsightsAnalysis_AnalysisComponent `json:"vpcEndpoint,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *NetworkInsightsAnalysis_Explanation) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkInsightsAnalysis.Explanation"
}
