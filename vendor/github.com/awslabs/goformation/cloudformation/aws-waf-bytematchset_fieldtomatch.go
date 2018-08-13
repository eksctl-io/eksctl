package cloudformation

import (
	"encoding/json"
)

// AWSWAFByteMatchSet_FieldToMatch AWS CloudFormation Resource (AWS::WAF::ByteMatchSet.FieldToMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-bytematchset-bytematchtuples-fieldtomatch.html
type AWSWAFByteMatchSet_FieldToMatch struct {

	// Data AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-bytematchset-bytematchtuples-fieldtomatch.html#cfn-waf-bytematchset-bytematchtuples-fieldtomatch-data
	Data *Value `json:"Data,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-bytematchset-bytematchtuples-fieldtomatch.html#cfn-waf-bytematchset-bytematchtuples-fieldtomatch-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFByteMatchSet_FieldToMatch) AWSCloudFormationType() string {
	return "AWS::WAF::ByteMatchSet.FieldToMatch"
}

func (r *AWSWAFByteMatchSet_FieldToMatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
