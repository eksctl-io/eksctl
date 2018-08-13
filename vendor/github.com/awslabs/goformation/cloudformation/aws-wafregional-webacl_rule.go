package cloudformation

import (
	"encoding/json"
)

// AWSWAFRegionalWebACL_Rule AWS CloudFormation Resource (AWS::WAFRegional::WebACL.Rule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-webacl-rule.html
type AWSWAFRegionalWebACL_Rule struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-webacl-rule.html#cfn-wafregional-webacl-rule-action
	Action *AWSWAFRegionalWebACL_Action `json:"Action,omitempty"`

	// Priority AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-webacl-rule.html#cfn-wafregional-webacl-rule-priority
	Priority *Value `json:"Priority,omitempty"`

	// RuleId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-webacl-rule.html#cfn-wafregional-webacl-rule-ruleid
	RuleId *Value `json:"RuleId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFRegionalWebACL_Rule) AWSCloudFormationType() string {
	return "AWS::WAFRegional::WebACL.Rule"
}

func (r *AWSWAFRegionalWebACL_Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
