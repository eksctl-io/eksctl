package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SpotFleet_NetworkInterfaceCountRequest AWS CloudFormation Resource (AWS::EC2::SpotFleet.NetworkInterfaceCountRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-networkinterfacecountrequest.html
type SpotFleet_NetworkInterfaceCountRequest struct {

	// Max AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-networkinterfacecountrequest.html#cfn-ec2-spotfleet-networkinterfacecountrequest-max
	Max *types.Value `json:"Max,omitempty"`

	// Min AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-networkinterfacecountrequest.html#cfn-ec2-spotfleet-networkinterfacecountrequest-min
	Min *types.Value `json:"Min,omitempty"`

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
func (r *SpotFleet_NetworkInterfaceCountRequest) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.NetworkInterfaceCountRequest"
}
