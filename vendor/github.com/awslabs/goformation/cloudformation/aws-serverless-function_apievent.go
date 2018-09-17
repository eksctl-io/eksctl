package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_ApiEvent AWS CloudFormation Resource (AWS::Serverless::Function.ApiEvent)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
type AWSServerlessFunction_ApiEvent struct {

	// Method AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
	Method *Value `json:"Method,omitempty"`

	// Path AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
	Path *Value `json:"Path,omitempty"`

	// RestApiId AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
	RestApiId *Value `json:"RestApiId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_ApiEvent) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.ApiEvent"
}

func (r *AWSServerlessFunction_ApiEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
