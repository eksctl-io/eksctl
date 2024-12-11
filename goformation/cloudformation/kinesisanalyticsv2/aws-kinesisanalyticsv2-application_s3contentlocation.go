package kinesisanalyticsv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Application_S3ContentLocation AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.S3ContentLocation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html
type Application_S3ContentLocation struct {

	// BucketARN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html#cfn-kinesisanalyticsv2-application-s3contentlocation-bucketarn
	BucketARN *types.Value `json:"BucketARN,omitempty"`

	// FileKey AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html#cfn-kinesisanalyticsv2-application-s3contentlocation-filekey
	FileKey *types.Value `json:"FileKey,omitempty"`

	// ObjectVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-s3contentlocation.html#cfn-kinesisanalyticsv2-application-s3contentlocation-objectversion
	ObjectVersion *types.Value `json:"ObjectVersion,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Application_S3ContentLocation) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.S3ContentLocation"
}
