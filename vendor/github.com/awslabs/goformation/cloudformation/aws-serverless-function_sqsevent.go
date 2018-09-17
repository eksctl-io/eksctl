package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_SQSEvent AWS CloudFormation Resource (AWS::Serverless::Function.SQSEvent)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#sqs
type AWSServerlessFunction_SQSEvent struct {

	// BatchSize AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#sqs
	BatchSize *Value `json:"BatchSize,omitempty"`

	// Queue AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#sqs
	Queue *Value `json:"Queue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_SQSEvent) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.SQSEvent"
}

func (r *AWSServerlessFunction_SQSEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
