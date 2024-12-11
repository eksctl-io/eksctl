package wafv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// RuleGroup_VisibilityConfig AWS CloudFormation Resource (AWS::WAFv2::RuleGroup.VisibilityConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-visibilityconfig.html
type RuleGroup_VisibilityConfig struct {

	// CloudWatchMetricsEnabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-visibilityconfig.html#cfn-wafv2-rulegroup-visibilityconfig-cloudwatchmetricsenabled
	CloudWatchMetricsEnabled *types.Value `json:"CloudWatchMetricsEnabled"`

	// MetricName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-visibilityconfig.html#cfn-wafv2-rulegroup-visibilityconfig-metricname
	MetricName *types.Value `json:"MetricName,omitempty"`

	// SampledRequestsEnabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-visibilityconfig.html#cfn-wafv2-rulegroup-visibilityconfig-sampledrequestsenabled
	SampledRequestsEnabled *types.Value `json:"SampledRequestsEnabled"`

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
func (r *RuleGroup_VisibilityConfig) AWSCloudFormationType() string {
	return "AWS::WAFv2::RuleGroup.VisibilityConfig"
}
