package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_StepFunctionsAction AWS CloudFormation Resource (AWS::IoT::TopicRule.StepFunctionsAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-stepfunctionsaction.html
type AWSIoTTopicRule_StepFunctionsAction struct {

	// ExecutionNamePrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-stepfunctionsaction.html#cfn-iot-topicrule-stepfunctionsaction-executionnameprefix
	ExecutionNamePrefix *Value `json:"ExecutionNamePrefix,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-stepfunctionsaction.html#cfn-iot-topicrule-stepfunctionsaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// StateMachineName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-stepfunctionsaction.html#cfn-iot-topicrule-stepfunctionsaction-statemachinename
	StateMachineName *Value `json:"StateMachineName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_StepFunctionsAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.StepFunctionsAction"
}

func (r *AWSIoTTopicRule_StepFunctionsAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
