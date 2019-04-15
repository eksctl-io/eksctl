package cloudformation

import (
	"encoding/json"
)

// AWSDLMLifecyclePolicy_Schedule AWS CloudFormation Resource (AWS::DLM::LifecyclePolicy.Schedule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-schedule.html
type AWSDLMLifecyclePolicy_Schedule struct {

	// CopyTags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-schedule.html#cfn-dlm-lifecyclepolicy-schedule-copytags
	CopyTags *Value `json:"CopyTags,omitempty"`

	// CreateRule AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-schedule.html#cfn-dlm-lifecyclepolicy-schedule-createrule
	CreateRule *AWSDLMLifecyclePolicy_CreateRule `json:"CreateRule,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-schedule.html#cfn-dlm-lifecyclepolicy-schedule-name
	Name *Value `json:"Name,omitempty"`

	// RetainRule AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-schedule.html#cfn-dlm-lifecyclepolicy-schedule-retainrule
	RetainRule *AWSDLMLifecyclePolicy_RetainRule `json:"RetainRule,omitempty"`

	// TagsToAdd AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-schedule.html#cfn-dlm-lifecyclepolicy-schedule-tagstoadd
	TagsToAdd []Tag `json:"TagsToAdd,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDLMLifecyclePolicy_Schedule) AWSCloudFormationType() string {
	return "AWS::DLM::LifecyclePolicy.Schedule"
}

func (r *AWSDLMLifecyclePolicy_Schedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
