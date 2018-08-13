package cloudformation

import (
	"encoding/json"
)

// AWSEventsRule_EcsParameters AWS CloudFormation Resource (AWS::Events::Rule.EcsParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-ecsparameters.html
type AWSEventsRule_EcsParameters struct {

	// TaskCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-ecsparameters.html#cfn-events-rule-ecsparameters-taskcount
	TaskCount *Value `json:"TaskCount,omitempty"`

	// TaskDefinitionArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-ecsparameters.html#cfn-events-rule-ecsparameters-taskdefinitionarn
	TaskDefinitionArn *Value `json:"TaskDefinitionArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEventsRule_EcsParameters) AWSCloudFormationType() string {
	return "AWS::Events::Rule.EcsParameters"
}

func (r *AWSEventsRule_EcsParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
