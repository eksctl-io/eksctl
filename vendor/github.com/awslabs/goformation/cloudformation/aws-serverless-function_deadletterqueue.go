package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_DeadLetterQueue AWS CloudFormation Resource (AWS::Serverless::Function.DeadLetterQueue)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#deadletterqueue-object
type AWSServerlessFunction_DeadLetterQueue struct {

	// TargetArn AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	TargetArn *Value `json:"TargetArn,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_DeadLetterQueue) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.DeadLetterQueue"
}

func (r *AWSServerlessFunction_DeadLetterQueue) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
