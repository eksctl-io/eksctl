package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// EndpointConfig_ServerlessConfig AWS CloudFormation Resource (AWS::SageMaker::EndpointConfig.ServerlessConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant-serverlessconfig.html
type EndpointConfig_ServerlessConfig struct {

	// MaxConcurrency AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant-serverlessconfig.html#cfn-sagemaker-endpointconfig-productionvariant-serverlessconfig-maxconcurrency
	MaxConcurrency *types.Value `json:"MaxConcurrency"`

	// MemorySizeInMB AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-productionvariant-serverlessconfig.html#cfn-sagemaker-endpointconfig-productionvariant-serverlessconfig-memorysizeinmb
	MemorySizeInMB *types.Value `json:"MemorySizeInMB"`

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
func (r *EndpointConfig_ServerlessConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::EndpointConfig.ServerlessConfig"
}
