package wafv2

import (
	"github.com/awslabs/goformation/v4/cloudformation/types"

	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// RuleGroup_RateBasedStatementOne AWS CloudFormation Resource (AWS::WAFv2::RuleGroup.RateBasedStatementOne)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-ratebasedstatementone.html
type RuleGroup_RateBasedStatementOne struct {

	// AggregateKeyType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-ratebasedstatementone.html#cfn-wafv2-rulegroup-ratebasedstatementone-aggregatekeytype
	AggregateKeyType *types.Value `json:"AggregateKeyType,omitempty"`

	// ForwardedIPConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-ratebasedstatementone.html#cfn-wafv2-rulegroup-ratebasedstatementone-forwardedipconfig
	ForwardedIPConfig *RuleGroup_ForwardedIPConfiguration `json:"ForwardedIPConfig,omitempty"`

	// Limit AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-ratebasedstatementone.html#cfn-wafv2-rulegroup-ratebasedstatementone-limit
	Limit *types.Value `json:"Limit"`

	// ScopeDownStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-rulegroup-ratebasedstatementone.html#cfn-wafv2-rulegroup-ratebasedstatementone-scopedownstatement
	ScopeDownStatement *RuleGroup_StatementTwo `json:"ScopeDownStatement,omitempty"`

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
func (r *RuleGroup_RateBasedStatementOne) AWSCloudFormationType() string {
	return "AWS::WAFv2::RuleGroup.RateBasedStatementOne"
}
