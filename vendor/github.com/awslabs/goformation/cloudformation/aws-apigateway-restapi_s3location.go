package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayRestApi_S3Location AWS CloudFormation Resource (AWS::ApiGateway::RestApi.S3Location)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-s3location.html
type AWSApiGatewayRestApi_S3Location struct {

	// Bucket AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-s3location.html#cfn-apigateway-restapi-s3location-bucket
	Bucket *Value `json:"Bucket,omitempty"`

	// ETag AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-s3location.html#cfn-apigateway-restapi-s3location-etag
	ETag *Value `json:"ETag,omitempty"`

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-s3location.html#cfn-apigateway-restapi-s3location-key
	Key *Value `json:"Key,omitempty"`

	// Version AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-s3location.html#cfn-apigateway-restapi-s3location-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayRestApi_S3Location) AWSCloudFormationType() string {
	return "AWS::ApiGateway::RestApi.S3Location"
}

func (r *AWSApiGatewayRestApi_S3Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
