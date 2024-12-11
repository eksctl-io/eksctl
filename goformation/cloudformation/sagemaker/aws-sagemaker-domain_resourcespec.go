package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Domain_ResourceSpec AWS CloudFormation Resource (AWS::SageMaker::Domain.ResourceSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-resourcespec.html
type Domain_ResourceSpec struct {

	// InstanceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-resourcespec.html#cfn-sagemaker-domain-resourcespec-instancetype
	InstanceType *types.Value `json:"InstanceType,omitempty"`

	// SageMakerImageArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-resourcespec.html#cfn-sagemaker-domain-resourcespec-sagemakerimagearn
	SageMakerImageArn *types.Value `json:"SageMakerImageArn,omitempty"`

	// SageMakerImageVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-resourcespec.html#cfn-sagemaker-domain-resourcespec-sagemakerimageversionarn
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
func (r *Domain_ResourceSpec) AWSCloudFormationType() string {
	return "AWS::SageMaker::Domain.ResourceSpec"
}
