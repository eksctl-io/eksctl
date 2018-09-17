package cloudformation

import (
	"encoding/json"
)

// AWSWAFIPSet_IPSetDescriptor AWS CloudFormation Resource (AWS::WAF::IPSet.IPSetDescriptor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-ipset-ipsetdescriptors.html
type AWSWAFIPSet_IPSetDescriptor struct {

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-ipset-ipsetdescriptors.html#cfn-waf-ipset-ipsetdescriptors-type
	Type *Value `json:"Type,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-ipset-ipsetdescriptors.html#cfn-waf-ipset-ipsetdescriptors-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFIPSet_IPSetDescriptor) AWSCloudFormationType() string {
	return "AWS::WAF::IPSet.IPSetDescriptor"
}

func (r *AWSWAFIPSet_IPSetDescriptor) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
