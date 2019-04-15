package cloudformation

import (
	"encoding/json"
)

// AWSDLMLifecyclePolicy_CreateRule AWS CloudFormation Resource (AWS::DLM::LifecyclePolicy.CreateRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-createrule.html
type AWSDLMLifecyclePolicy_CreateRule struct {

	// Interval AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-createrule.html#cfn-dlm-lifecyclepolicy-createrule-interval
	Interval *Value `json:"Interval,omitempty"`

	// IntervalUnit AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-createrule.html#cfn-dlm-lifecyclepolicy-createrule-intervalunit
	IntervalUnit *Value `json:"IntervalUnit,omitempty"`

	// Times AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-createrule.html#cfn-dlm-lifecyclepolicy-createrule-times
	Times []*Value `json:"Times,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDLMLifecyclePolicy_CreateRule) AWSCloudFormationType() string {
	return "AWS::DLM::LifecyclePolicy.CreateRule"
}

func (r *AWSDLMLifecyclePolicy_CreateRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
