package cloudformation

import (
	"encoding/json"
)

// AWSStepFunctionsStateMachine_TagsEntry AWS CloudFormation Resource (AWS::StepFunctions::StateMachine.TagsEntry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stepfunctions-statemachine-tagsentry.html
type AWSStepFunctionsStateMachine_TagsEntry struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stepfunctions-statemachine-tagsentry.html#cfn-stepfunctions-statemachine-tagsentry-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stepfunctions-statemachine-tagsentry.html#cfn-stepfunctions-statemachine-tagsentry-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSStepFunctionsStateMachine_TagsEntry) AWSCloudFormationType() string {
	return "AWS::StepFunctions::StateMachine.TagsEntry"
}

func (r *AWSStepFunctionsStateMachine_TagsEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
