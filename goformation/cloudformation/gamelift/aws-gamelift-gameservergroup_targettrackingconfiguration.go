package gamelift

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GameServerGroup_TargetTrackingConfiguration AWS CloudFormation Resource (AWS::GameLift::GameServerGroup.TargetTrackingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-targettrackingconfiguration.html
type GameServerGroup_TargetTrackingConfiguration struct {

	// TargetValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-targettrackingconfiguration.html#cfn-gamelift-gameservergroup-targettrackingconfiguration-targetvalue
	TargetValue *types.Value `json:"TargetValue"`

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
func (r *GameServerGroup_TargetTrackingConfiguration) AWSCloudFormationType() string {
	return "AWS::GameLift::GameServerGroup.TargetTrackingConfiguration"
}
