package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_FunctionEnvironment AWS CloudFormation Resource (AWS::Serverless::Function.FunctionEnvironment)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#environment-object
type AWSServerlessFunction_FunctionEnvironment struct {

	// Variables AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#environment-object
	Variables map[string]*Value `json:"Variables,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_FunctionEnvironment) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.FunctionEnvironment"
}

func (r *AWSServerlessFunction_FunctionEnvironment) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
