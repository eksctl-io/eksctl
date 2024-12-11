package ecs

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CapacityProvider_AutoScalingGroupProvider AWS CloudFormation Resource (AWS::ECS::CapacityProvider.AutoScalingGroupProvider)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-autoscalinggroupprovider.html
type CapacityProvider_AutoScalingGroupProvider struct {

	// AutoScalingGroupArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-autoscalinggroupprovider.html#cfn-ecs-capacityprovider-autoscalinggroupprovider-autoscalinggrouparn
	AutoScalingGroupArn *types.Value `json:"AutoScalingGroupArn,omitempty"`

	// ManagedScaling AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-autoscalinggroupprovider.html#cfn-ecs-capacityprovider-autoscalinggroupprovider-managedscaling
	ManagedScaling *CapacityProvider_ManagedScaling `json:"ManagedScaling,omitempty"`

	// ManagedTerminationProtection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-autoscalinggroupprovider.html#cfn-ecs-capacityprovider-autoscalinggroupprovider-managedterminationprotection
	ManagedTerminationProtection *types.Value `json:"ManagedTerminationProtection,omitempty"`

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
func (r *CapacityProvider_AutoScalingGroupProvider) AWSCloudFormationType() string {
	return "AWS::ECS::CapacityProvider.AutoScalingGroupProvider"
}
