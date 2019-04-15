package cloudformation

import (
	"encoding/json"
)

// AWSLambdaLayerVersion_Content AWS CloudFormation Resource (AWS::Lambda::LayerVersion.Content)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-layerversion-content.html
type AWSLambdaLayerVersion_Content struct {

	// S3Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-layerversion-content.html#cfn-lambda-layerversion-content-s3bucket
	S3Bucket *Value `json:"S3Bucket,omitempty"`

	// S3Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-layerversion-content.html#cfn-lambda-layerversion-content-s3key
	S3Key *Value `json:"S3Key,omitempty"`

	// S3ObjectVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-layerversion-content.html#cfn-lambda-layerversion-content-s3objectversion
	S3ObjectVersion *Value `json:"S3ObjectVersion,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSLambdaLayerVersion_Content) AWSCloudFormationType() string {
	return "AWS::Lambda::LayerVersion.Content"
}

func (r *AWSLambdaLayerVersion_Content) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
