package kinesisanalyticsv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Application_CustomArtifactConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.CustomArtifactConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-customartifactconfiguration.html
type Application_CustomArtifactConfiguration struct {

	// ArtifactType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-customartifactconfiguration.html#cfn-kinesisanalyticsv2-application-customartifactconfiguration-artifacttype
	ArtifactType *types.Value `json:"ArtifactType,omitempty"`

	// MavenReference AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-customartifactconfiguration.html#cfn-kinesisanalyticsv2-application-customartifactconfiguration-mavenreference
	MavenReference *Application_MavenReference `json:"MavenReference,omitempty"`

	// S3ContentLocation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-customartifactconfiguration.html#cfn-kinesisanalyticsv2-application-customartifactconfiguration-s3contentlocation
	S3ContentLocation *Application_S3ContentLocation `json:"S3ContentLocation,omitempty"`

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
func (r *Application_CustomArtifactConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.CustomArtifactConfiguration"
}
