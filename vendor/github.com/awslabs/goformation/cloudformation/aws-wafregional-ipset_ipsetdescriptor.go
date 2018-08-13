package cloudformation

import (
	"encoding/json"
)

// AWSWAFRegionalIPSet_IPSetDescriptor AWS CloudFormation Resource (AWS::WAFRegional::IPSet.IPSetDescriptor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-ipset-ipsetdescriptor.html
type AWSWAFRegionalIPSet_IPSetDescriptor struct {

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-ipset-ipsetdescriptor.html#cfn-wafregional-ipset-ipsetdescriptor-type
	Type *Value `json:"Type,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-ipset-ipsetdescriptor.html#cfn-wafregional-ipset-ipsetdescriptor-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFRegionalIPSet_IPSetDescriptor) AWSCloudFormationType() string {
	return "AWS::WAFRegional::IPSet.IPSetDescriptor"
}

func (r *AWSWAFRegionalIPSet_IPSetDescriptor) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
