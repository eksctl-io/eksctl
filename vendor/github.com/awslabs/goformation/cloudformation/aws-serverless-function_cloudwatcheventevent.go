package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_CloudWatchEventEvent AWS CloudFormation Resource (AWS::Serverless::Function.CloudWatchEventEvent)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cloudwatchevent
type AWSServerlessFunction_CloudWatchEventEvent struct {

	// Input AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cloudwatchevent
	Input *Value `json:"Input,omitempty"`

	// InputPath AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cloudwatchevent
	InputPath *Value `json:"InputPath,omitempty"`

	// Pattern AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AmazonCloudWatch/latest/events/CloudWatchEventsandEventPatterns.html
	Pattern interface{} `json:"Pattern,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_CloudWatchEventEvent) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.CloudWatchEventEvent"
}

func (r *AWSServerlessFunction_CloudWatchEventEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
