package gamelift

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GameServerGroup_VpcSubnets AWS CloudFormation Resource (AWS::GameLift::GameServerGroup.VpcSubnets)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-vpcsubnets.html
type GameServerGroup_VpcSubnets struct {

	// VpcSubnets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-vpcsubnets.html#cfn-gamelift-gameservergroup-vpcsubnets-vpcsubnets
	VpcSubnets *types.Value `json:"VpcSubnets,omitempty"`

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
func (r *GameServerGroup_VpcSubnets) AWSCloudFormationType() string {
	return "AWS::GameLift::GameServerGroup.VpcSubnets"
}
