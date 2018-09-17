package cloudformation

import (
	"encoding/json"
)

// AWSLambdaFunction_Code AWS CloudFormation Resource (AWS::Lambda::Function.Code)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-code.html
type AWSLambdaFunction_Code struct {

	// S3Bucket AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-code.html#cfn-lambda-function-code-s3bucket
	S3Bucket *Value `json:"S3Bucket,omitempty"`

	// S3Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-code.html#cfn-lambda-function-code-s3key
	S3Key *Value `json:"S3Key,omitempty"`

	// S3ObjectVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-code.html#cfn-lambda-function-code-s3objectversion
	S3ObjectVersion *Value `json:"S3ObjectVersion,omitempty"`

	// ZipFile AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-code.html#cfn-lambda-function-code-zipfile
	ZipFile *Value `json:"ZipFile,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSLambdaFunction_Code) AWSCloudFormationType() string {
	return "AWS::Lambda::Function.Code"
}

func (r *AWSLambdaFunction_Code) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
