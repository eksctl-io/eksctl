package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// NetworkInterface_ConnectionTrackingSpecification AWS CloudFormation Resource (AWS::EC2::NetworkInterface.ConnectionTrackingSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinterface-connectiontrackingspecification.html
type NetworkInterface_ConnectionTrackingSpecification struct {

	// TcpEstablishedTimeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinterface-connectiontrackingspecification.html#cfn-ec2-networkinterface-connectiontrackingspecification-tcpestablishedtimeout
	TcpEstablishedTimeout *types.Value `json:"TcpEstablishedTimeout,omitempty"`

	// UdpStreamTimeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinterface-connectiontrackingspecification.html#cfn-ec2-networkinterface-connectiontrackingspecification-udpstreamtimeout
	UdpStreamTimeout *types.Value `json:"UdpStreamTimeout,omitempty"`

	// UdpTimeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinterface-connectiontrackingspecification.html#cfn-ec2-networkinterface-connectiontrackingspecification-udptimeout
	UdpTimeout *types.Value `json:"UdpTimeout,omitempty"`

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
func (r *NetworkInterface_ConnectionTrackingSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkInterface.ConnectionTrackingSpecification"
}
