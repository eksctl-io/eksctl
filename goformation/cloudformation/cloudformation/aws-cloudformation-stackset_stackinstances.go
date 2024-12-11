package cloudformation

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StackSet_StackInstances AWS CloudFormation Resource (AWS::CloudFormation::StackSet.StackInstances)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-stackinstances.html
type StackSet_StackInstances struct {

	// DeploymentTargets AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-stackinstances.html#cfn-cloudformation-stackset-stackinstances-deploymenttargets
	DeploymentTargets *StackSet_DeploymentTargets `json:"DeploymentTargets,omitempty"`

	// ParameterOverrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-stackinstances.html#cfn-cloudformation-stackset-stackinstances-parameteroverrides
	ParameterOverrides []StackSet_Parameter `json:"ParameterOverrides,omitempty"`

	// Regions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-stackinstances.html#cfn-cloudformation-stackset-stackinstances-regions
	Regions *types.Value `json:"Regions,omitempty"`

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
func (r *StackSet_StackInstances) AWSCloudFormationType() string {
	return "AWS::CloudFormation::StackSet.StackInstances"
}
