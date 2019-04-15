package cloudformation

import (
	"encoding/json"
)

// AWSStepFunctionsActivity_TagsEntry AWS CloudFormation Resource (AWS::StepFunctions::Activity.TagsEntry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stepfunctions-activity-tagsentry.html
type AWSStepFunctionsActivity_TagsEntry struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stepfunctions-activity-tagsentry.html#cfn-stepfunctions-activity-tagsentry-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stepfunctions-activity-tagsentry.html#cfn-stepfunctions-activity-tagsentry-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSStepFunctionsActivity_TagsEntry) AWSCloudFormationType() string {
	return "AWS::StepFunctions::Activity.TagsEntry"
}

func (r *AWSStepFunctionsActivity_TagsEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
