package cloudformation

import (
	"encoding/json"
)

// AWSWAFRule_Predicate AWS CloudFormation Resource (AWS::WAF::Rule.Predicate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-rule-predicates.html
type AWSWAFRule_Predicate struct {

	// DataId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-rule-predicates.html#cfn-waf-rule-predicates-dataid
	DataId *Value `json:"DataId,omitempty"`

	// Negated AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-rule-predicates.html#cfn-waf-rule-predicates-negated
	Negated *Value `json:"Negated,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-rule-predicates.html#cfn-waf-rule-predicates-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFRule_Predicate) AWSCloudFormationType() string {
	return "AWS::WAF::Rule.Predicate"
}

func (r *AWSWAFRule_Predicate) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
