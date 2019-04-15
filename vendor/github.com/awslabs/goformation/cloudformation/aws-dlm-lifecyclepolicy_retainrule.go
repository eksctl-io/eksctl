package cloudformation

import (
	"encoding/json"
)

// AWSDLMLifecyclePolicy_RetainRule AWS CloudFormation Resource (AWS::DLM::LifecyclePolicy.RetainRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-retainrule.html
type AWSDLMLifecyclePolicy_RetainRule struct {

	// Count AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-retainrule.html#cfn-dlm-lifecyclepolicy-retainrule-count
	Count *Value `json:"Count,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDLMLifecyclePolicy_RetainRule) AWSCloudFormationType() string {
	return "AWS::DLM::LifecyclePolicy.RetainRule"
}

func (r *AWSDLMLifecyclePolicy_RetainRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
