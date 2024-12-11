package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// UserProfile_ResourceSpec AWS CloudFormation Resource (AWS::SageMaker::UserProfile.ResourceSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-resourcespec.html
type UserProfile_ResourceSpec struct {

	// InstanceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-resourcespec.html#cfn-sagemaker-userprofile-resourcespec-instancetype
	InstanceType *types.Value `json:"InstanceType,omitempty"`

	// SageMakerImageArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-resourcespec.html#cfn-sagemaker-userprofile-resourcespec-sagemakerimagearn
	SageMakerImageArn *types.Value `json:"SageMakerImageArn,omitempty"`

	// SageMakerImageVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-resourcespec.html#cfn-sagemaker-userprofile-resourcespec-sagemakerimageversionarn
	SageMakerImageVersionArn *types.Value `json:"SageMakerImageVersionArn,omitempty"`

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
func (r *UserProfile_ResourceSpec) AWSCloudFormationType() string {
	return "AWS::SageMaker::UserProfile.ResourceSpec"
}
