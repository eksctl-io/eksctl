package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_S3Location AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.S3Location)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-s3location.html
type AWSCodeDeployDeploymentGroup_S3Location struct {

	// Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-s3location.html#cfn-properties-codedeploy-deploymentgroup-deployment-revision-s3location-bucket
	Bucket *Value `json:"Bucket,omitempty"`

	// BundleType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-s3location.html#cfn-properties-codedeploy-deploymentgroup-deployment-revision-s3location-bundletype
	BundleType *Value `json:"BundleType,omitempty"`

	// ETag AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-s3location.html#cfn-properties-codedeploy-deploymentgroup-deployment-revision-s3location-etag
	ETag *Value `json:"ETag,omitempty"`

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-s3location.html#cfn-properties-codedeploy-deploymentgroup-deployment-revision-s3location-key
	Key *Value `json:"Key,omitempty"`

	// Version AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-s3location.html#cfn-properties-codedeploy-deploymentgroup-deployment-revision-s3location-value
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_S3Location) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.S3Location"
}

func (r *AWSCodeDeployDeploymentGroup_S3Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
