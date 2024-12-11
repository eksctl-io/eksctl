package sagemaker

import (
	"goformation/v4/cloudformation/policies"
)

// EndpointConfig_AsyncInferenceConfig AWS CloudFormation Resource (AWS::SageMaker::EndpointConfig.AsyncInferenceConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-asyncinferenceconfig.html
type EndpointConfig_AsyncInferenceConfig struct {

	// ClientConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-asyncinferenceconfig.html#cfn-sagemaker-endpointconfig-asyncinferenceconfig-clientconfig
	ClientConfig *EndpointConfig_AsyncInferenceClientConfig `json:"ClientConfig,omitempty"`

	// OutputConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-asyncinferenceconfig.html#cfn-sagemaker-endpointconfig-asyncinferenceconfig-outputconfig
	OutputConfig *EndpointConfig_AsyncInferenceOutputConfig `json:"OutputConfig,omitempty"`

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
func (r *EndpointConfig_AsyncInferenceConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::EndpointConfig.AsyncInferenceConfig"
}
