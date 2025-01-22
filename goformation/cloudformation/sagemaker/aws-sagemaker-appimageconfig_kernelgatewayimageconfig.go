package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// AppImageConfig_KernelGatewayImageConfig AWS CloudFormation Resource (AWS::SageMaker::AppImageConfig.KernelGatewayImageConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-appimageconfig-kernelgatewayimageconfig.html
type AppImageConfig_KernelGatewayImageConfig struct {

	// FileSystemConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-appimageconfig-kernelgatewayimageconfig.html#cfn-sagemaker-appimageconfig-kernelgatewayimageconfig-filesystemconfig
	FileSystemConfig *AppImageConfig_FileSystemConfig `json:"FileSystemConfig,omitempty"`

	// KernelSpecs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-appimageconfig-kernelgatewayimageconfig.html#cfn-sagemaker-appimageconfig-kernelgatewayimageconfig-kernelspecs
	KernelSpecs []AppImageConfig_KernelSpec `json:"KernelSpecs,omitempty"`

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
func (r *AppImageConfig_KernelGatewayImageConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::AppImageConfig.KernelGatewayImageConfig"
}
