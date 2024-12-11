package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CodeRepository_GitConfig AWS CloudFormation Resource (AWS::SageMaker::CodeRepository.GitConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-coderepository-gitconfig.html
type CodeRepository_GitConfig struct {

	// Branch AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-coderepository-gitconfig.html#cfn-sagemaker-coderepository-gitconfig-branch
	Branch *types.Value `json:"Branch,omitempty"`

	// RepositoryUrl AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-coderepository-gitconfig.html#cfn-sagemaker-coderepository-gitconfig-repositoryurl
	RepositoryUrl *types.Value `json:"RepositoryUrl,omitempty"`

	// SecretArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-coderepository-gitconfig.html#cfn-sagemaker-coderepository-gitconfig-secretarn
	SecretArn *types.Value `json:"SecretArn,omitempty"`

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
func (r *CodeRepository_GitConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::CodeRepository.GitConfig"
}
