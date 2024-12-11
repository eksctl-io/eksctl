package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Domain_CustomImage AWS CloudFormation Resource (AWS::SageMaker::Domain.CustomImage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-customimage.html
type Domain_CustomImage struct {

	// AppImageConfigName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-customimage.html#cfn-sagemaker-domain-customimage-appimageconfigname
	AppImageConfigName *types.Value `json:"AppImageConfigName,omitempty"`

	// ImageName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-customimage.html#cfn-sagemaker-domain-customimage-imagename
	ImageName *types.Value `json:"ImageName,omitempty"`

	// ImageVersionNumber AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-customimage.html#cfn-sagemaker-domain-customimage-imageversionnumber
	ImageVersionNumber *types.Value `json:"ImageVersionNumber,omitempty"`

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
func (r *Domain_CustomImage) AWSCloudFormationType() string {
	return "AWS::SageMaker::Domain.CustomImage"
}
