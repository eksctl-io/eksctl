package gamelift

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GameServerGroup_AutoScalingPolicy AWS CloudFormation Resource (AWS::GameLift::GameServerGroup.AutoScalingPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-autoscalingpolicy.html
type GameServerGroup_AutoScalingPolicy struct {

	// EstimatedInstanceWarmup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-autoscalingpolicy.html#cfn-gamelift-gameservergroup-autoscalingpolicy-estimatedinstancewarmup
	EstimatedInstanceWarmup *types.Value `json:"EstimatedInstanceWarmup,omitempty"`

	// TargetTrackingConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-autoscalingpolicy.html#cfn-gamelift-gameservergroup-autoscalingpolicy-targettrackingconfiguration
	TargetTrackingConfiguration *GameServerGroup_TargetTrackingConfiguration `json:"TargetTrackingConfiguration,omitempty"`

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
func (r *GameServerGroup_AutoScalingPolicy) AWSCloudFormationType() string {
	return "AWS::GameLift::GameServerGroup.AutoScalingPolicy"
}
