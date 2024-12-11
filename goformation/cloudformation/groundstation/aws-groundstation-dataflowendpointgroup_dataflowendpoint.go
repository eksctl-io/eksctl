package groundstation

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataflowEndpointGroup_DataflowEndpoint AWS CloudFormation Resource (AWS::GroundStation::DataflowEndpointGroup.DataflowEndpoint)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-dataflowendpointgroup-dataflowendpoint.html
type DataflowEndpointGroup_DataflowEndpoint struct {

	// Address AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-dataflowendpointgroup-dataflowendpoint.html#cfn-groundstation-dataflowendpointgroup-dataflowendpoint-address
	Address *DataflowEndpointGroup_SocketAddress `json:"Address,omitempty"`

	// Mtu AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-dataflowendpointgroup-dataflowendpoint.html#cfn-groundstation-dataflowendpointgroup-dataflowendpoint-mtu
	Mtu *types.Value `json:"Mtu,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-dataflowendpointgroup-dataflowendpoint.html#cfn-groundstation-dataflowendpointgroup-dataflowendpoint-name
	Name *types.Value `json:"Name,omitempty"`

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
func (r *DataflowEndpointGroup_DataflowEndpoint) AWSCloudFormationType() string {
	return "AWS::GroundStation::DataflowEndpointGroup.DataflowEndpoint"
}
