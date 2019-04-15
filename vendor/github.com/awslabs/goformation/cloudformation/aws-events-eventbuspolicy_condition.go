package cloudformation

import (
	"encoding/json"
)

// AWSEventsEventBusPolicy_Condition AWS CloudFormation Resource (AWS::Events::EventBusPolicy.Condition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-eventbuspolicy-condition.html
type AWSEventsEventBusPolicy_Condition struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-eventbuspolicy-condition.html#cfn-events-eventbuspolicy-condition-key
	Key *Value `json:"Key,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-eventbuspolicy-condition.html#cfn-events-eventbuspolicy-condition-type
	Type *Value `json:"Type,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-eventbuspolicy-condition.html#cfn-events-eventbuspolicy-condition-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEventsEventBusPolicy_Condition) AWSCloudFormationType() string {
	return "AWS::Events::EventBusPolicy.Condition"
}

func (r *AWSEventsEventBusPolicy_Condition) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
