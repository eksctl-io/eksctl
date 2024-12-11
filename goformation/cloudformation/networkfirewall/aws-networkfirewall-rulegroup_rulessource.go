package networkfirewall

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// RuleGroup_RulesSource AWS CloudFormation Resource (AWS::NetworkFirewall::RuleGroup.RulesSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-rulessource.html
type RuleGroup_RulesSource struct {

	// RulesSourceList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-rulessource.html#cfn-networkfirewall-rulegroup-rulessource-rulessourcelist
	RulesSourceList *RuleGroup_RulesSourceList `json:"RulesSourceList,omitempty"`

	// RulesString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-rulessource.html#cfn-networkfirewall-rulegroup-rulessource-rulesstring
	RulesString *types.Value `json:"RulesString,omitempty"`

	// StatefulRules AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-rulessource.html#cfn-networkfirewall-rulegroup-rulessource-statefulrules
	StatefulRules []RuleGroup_StatefulRule `json:"StatefulRules,omitempty"`

	// StatelessRulesAndCustomActions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-rulessource.html#cfn-networkfirewall-rulegroup-rulessource-statelessrulesandcustomactions
	StatelessRulesAndCustomActions *RuleGroup_StatelessRulesAndCustomActions `json:"StatelessRulesAndCustomActions,omitempty"`

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
func (r *RuleGroup_RulesSource) AWSCloudFormationType() string {
	return "AWS::NetworkFirewall::RuleGroup.RulesSource"
}
