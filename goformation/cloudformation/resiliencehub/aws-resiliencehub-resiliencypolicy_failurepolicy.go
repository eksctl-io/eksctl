package resiliencehub

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResiliencyPolicy_FailurePolicy AWS CloudFormation Resource (AWS::ResilienceHub::ResiliencyPolicy.FailurePolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-resiliencehub-resiliencypolicy-failurepolicy.html
type ResiliencyPolicy_FailurePolicy struct {

	// RpoInSecs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-resiliencehub-resiliencypolicy-failurepolicy.html#cfn-resiliencehub-resiliencypolicy-failurepolicy-rpoinsecs
	RpoInSecs *types.Value `json:"RpoInSecs"`

	// RtoInSecs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-resiliencehub-resiliencypolicy-failurepolicy.html#cfn-resiliencehub-resiliencypolicy-failurepolicy-rtoinsecs
	RtoInSecs *types.Value `json:"RtoInSecs"`

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
func (r *ResiliencyPolicy_FailurePolicy) AWSCloudFormationType() string {
	return "AWS::ResilienceHub::ResiliencyPolicy.FailurePolicy"
}
