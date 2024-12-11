package xray

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SamplingRule_SamplingRuleRecord AWS CloudFormation Resource (AWS::XRay::SamplingRule.SamplingRuleRecord)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingrulerecord.html
type SamplingRule_SamplingRuleRecord struct {

	// CreatedAt AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingrulerecord.html#cfn-xray-samplingrule-samplingrulerecord-createdat
	CreatedAt *types.Value `json:"CreatedAt,omitempty"`

	// ModifiedAt AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingrulerecord.html#cfn-xray-samplingrule-samplingrulerecord-modifiedat
	ModifiedAt *types.Value `json:"ModifiedAt,omitempty"`

	// SamplingRule AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingrulerecord.html#cfn-xray-samplingrule-samplingrulerecord-samplingrule
	SamplingRule *SamplingRule_SamplingRule `json:"SamplingRule,omitempty"`

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
func (r *SamplingRule_SamplingRuleRecord) AWSCloudFormationType() string {
	return "AWS::XRay::SamplingRule.SamplingRuleRecord"
}
