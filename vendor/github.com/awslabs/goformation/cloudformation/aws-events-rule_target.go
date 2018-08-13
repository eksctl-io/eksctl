package cloudformation

import (
	"encoding/json"
)

// AWSEventsRule_Target AWS CloudFormation Resource (AWS::Events::Rule.Target)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html
type AWSEventsRule_Target struct {

	// Arn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-arn
	Arn *Value `json:"Arn,omitempty"`

	// EcsParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-ecsparameters
	EcsParameters *AWSEventsRule_EcsParameters `json:"EcsParameters,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-id
	Id *Value `json:"Id,omitempty"`

	// Input AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-input
	Input *Value `json:"Input,omitempty"`

	// InputPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-inputpath
	InputPath *Value `json:"InputPath,omitempty"`

	// InputTransformer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-inputtransformer
	InputTransformer *AWSEventsRule_InputTransformer `json:"InputTransformer,omitempty"`

	// KinesisParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-kinesisparameters
	KinesisParameters *AWSEventsRule_KinesisParameters `json:"KinesisParameters,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// RunCommandParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-target.html#cfn-events-rule-target-runcommandparameters
	RunCommandParameters *AWSEventsRule_RunCommandParameters `json:"RunCommandParameters,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEventsRule_Target) AWSCloudFormationType() string {
	return "AWS::Events::Rule.Target"
}

func (r *AWSEventsRule_Target) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
