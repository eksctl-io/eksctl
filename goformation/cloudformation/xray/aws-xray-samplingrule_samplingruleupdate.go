package xray

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SamplingRule_SamplingRuleUpdate AWS CloudFormation Resource (AWS::XRay::SamplingRule.SamplingRuleUpdate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html
type SamplingRule_SamplingRuleUpdate struct {

	// Attributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-attributes
	Attributes map[string]*types.Value `json:"Attributes,omitempty"`

	// FixedRate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-fixedrate
	FixedRate *types.Value `json:"FixedRate,omitempty"`

	// HTTPMethod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-httpmethod
	HTTPMethod *types.Value `json:"HTTPMethod,omitempty"`

	// Host AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-host
	Host *types.Value `json:"Host,omitempty"`

	// Priority AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-priority
	Priority *types.Value `json:"Priority,omitempty"`

	// ReservoirSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-reservoirsize
	ReservoirSize *types.Value `json:"ReservoirSize,omitempty"`

	// ResourceARN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-resourcearn
	ResourceARN *types.Value `json:"ResourceARN,omitempty"`

	// RuleARN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-rulearn
	RuleARN *types.Value `json:"RuleARN,omitempty"`

	// RuleName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-rulename
	RuleName *types.Value `json:"RuleName,omitempty"`

	// ServiceName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-servicename
	ServiceName *types.Value `json:"ServiceName,omitempty"`

	// ServiceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-servicetype
	ServiceType *types.Value `json:"ServiceType,omitempty"`

	// URLPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-xray-samplingrule-samplingruleupdate.html#cfn-xray-samplingrule-samplingruleupdate-urlpath
	URLPath *types.Value `json:"URLPath,omitempty"`

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
func (r *SamplingRule_SamplingRuleUpdate) AWSCloudFormationType() string {
	return "AWS::XRay::SamplingRule.SamplingRuleUpdate"
}
