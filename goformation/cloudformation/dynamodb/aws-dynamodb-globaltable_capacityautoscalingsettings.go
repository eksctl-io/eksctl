package dynamodb

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GlobalTable_CapacityAutoScalingSettings AWS CloudFormation Resource (AWS::DynamoDB::GlobalTable.CapacityAutoScalingSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-capacityautoscalingsettings.html
type GlobalTable_CapacityAutoScalingSettings struct {

	// MaxCapacity AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-capacityautoscalingsettings.html#cfn-dynamodb-globaltable-capacityautoscalingsettings-maxcapacity
	MaxCapacity *types.Value `json:"MaxCapacity"`

	// MinCapacity AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-capacityautoscalingsettings.html#cfn-dynamodb-globaltable-capacityautoscalingsettings-mincapacity
	MinCapacity *types.Value `json:"MinCapacity"`

	// SeedCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-capacityautoscalingsettings.html#cfn-dynamodb-globaltable-capacityautoscalingsettings-seedcapacity
	SeedCapacity *types.Value `json:"SeedCapacity,omitempty"`

	// TargetTrackingScalingPolicyConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-capacityautoscalingsettings.html#cfn-dynamodb-globaltable-capacityautoscalingsettings-targettrackingscalingpolicyconfiguration
	TargetTrackingScalingPolicyConfiguration *GlobalTable_TargetTrackingScalingPolicyConfiguration `json:"TargetTrackingScalingPolicyConfiguration,omitempty"`

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
func (r *GlobalTable_CapacityAutoScalingSettings) AWSCloudFormationType() string {
	return "AWS::DynamoDB::GlobalTable.CapacityAutoScalingSettings"
}
