package cloudformation

import (
	"encoding/json"
)

// AWSWAFSqlInjectionMatchSet_FieldToMatch AWS CloudFormation Resource (AWS::WAF::SqlInjectionMatchSet.FieldToMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-bytematchset-bytematchtuples-fieldtomatch.html
type AWSWAFSqlInjectionMatchSet_FieldToMatch struct {

	// Data AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-bytematchset-bytematchtuples-fieldtomatch.html#cfn-waf-sizeconstraintset-sizeconstraint-fieldtomatch-data
	Data *Value `json:"Data,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-bytematchset-bytematchtuples-fieldtomatch.html#cfn-waf-sizeconstraintset-sizeconstraint-fieldtomatch-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFSqlInjectionMatchSet_FieldToMatch) AWSCloudFormationType() string {
	return "AWS::WAF::SqlInjectionMatchSet.FieldToMatch"
}

func (r *AWSWAFSqlInjectionMatchSet_FieldToMatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
