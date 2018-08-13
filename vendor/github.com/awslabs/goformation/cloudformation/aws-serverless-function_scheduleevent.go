package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_ScheduleEvent AWS CloudFormation Resource (AWS::Serverless::Function.ScheduleEvent)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#schedule
type AWSServerlessFunction_ScheduleEvent struct {

	// Input AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#schedule
	Input *Value `json:"Input,omitempty"`

	// Schedule AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#schedule
	Schedule *Value `json:"Schedule,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_ScheduleEvent) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.ScheduleEvent"
}

func (r *AWSServerlessFunction_ScheduleEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
