package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// UserProfile_SharingSettings AWS CloudFormation Resource (AWS::SageMaker::UserProfile.SharingSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-sharingsettings.html
type UserProfile_SharingSettings struct {

	// NotebookOutputOption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-sharingsettings.html#cfn-sagemaker-userprofile-sharingsettings-notebookoutputoption
	NotebookOutputOption *types.Value `json:"NotebookOutputOption,omitempty"`

	// S3KmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-sharingsettings.html#cfn-sagemaker-userprofile-sharingsettings-s3kmskeyid
	S3KmsKeyId *types.Value `json:"S3KmsKeyId,omitempty"`

	// S3OutputPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-sharingsettings.html#cfn-sagemaker-userprofile-sharingsettings-s3outputpath
	S3OutputPath *types.Value `json:"S3OutputPath,omitempty"`

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
func (r *UserProfile_SharingSettings) AWSCloudFormationType() string {
	return "AWS::SageMaker::UserProfile.SharingSettings"
}
