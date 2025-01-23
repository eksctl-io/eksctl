package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// UserProfile_KernelGatewayAppSettings AWS CloudFormation Resource (AWS::SageMaker::UserProfile.KernelGatewayAppSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-kernelgatewayappsettings.html
type UserProfile_KernelGatewayAppSettings struct {

	// CustomImages AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-kernelgatewayappsettings.html#cfn-sagemaker-userprofile-kernelgatewayappsettings-customimages
	CustomImages []UserProfile_CustomImage `json:"CustomImages,omitempty"`

	// DefaultResourceSpec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-userprofile-kernelgatewayappsettings.html#cfn-sagemaker-userprofile-kernelgatewayappsettings-defaultresourcespec
	DefaultResourceSpec *UserProfile_ResourceSpec `json:"DefaultResourceSpec,omitempty"`

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
func (r *UserProfile_KernelGatewayAppSettings) AWSCloudFormationType() string {
	return "AWS::SageMaker::UserProfile.KernelGatewayAppSettings"
}
