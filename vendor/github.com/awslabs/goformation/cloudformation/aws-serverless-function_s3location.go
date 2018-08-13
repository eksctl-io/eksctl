package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_S3Location AWS CloudFormation Resource (AWS::Serverless::Function.S3Location)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#s3-location-object
type AWSServerlessFunction_S3Location struct {

	// Bucket AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Bucket *Value `json:"Bucket,omitempty"`

	// Key AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Key *Value `json:"Key,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_S3Location) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.S3Location"
}

func (r *AWSServerlessFunction_S3Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
