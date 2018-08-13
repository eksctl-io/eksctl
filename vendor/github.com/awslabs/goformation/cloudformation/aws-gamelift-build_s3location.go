package cloudformation

import (
	"encoding/json"
)

// AWSGameLiftBuild_S3Location AWS CloudFormation Resource (AWS::GameLift::Build.S3Location)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-build-storagelocation.html
type AWSGameLiftBuild_S3Location struct {

	// Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-build-storagelocation.html#cfn-gamelift-build-storage-bucket
	Bucket *Value `json:"Bucket,omitempty"`

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-build-storagelocation.html#cfn-gamelift-build-storage-key
	Key *Value `json:"Key,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-build-storagelocation.html#cfn-gamelift-build-storage-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGameLiftBuild_S3Location) AWSCloudFormationType() string {
	return "AWS::GameLift::Build.S3Location"
}

func (r *AWSGameLiftBuild_S3Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
