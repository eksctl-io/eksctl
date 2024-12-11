package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// EndpointConfig_AsyncInferenceNotificationConfig AWS CloudFormation Resource (AWS::SageMaker::EndpointConfig.AsyncInferenceNotificationConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-asyncinferencenotificationconfig.html
type EndpointConfig_AsyncInferenceNotificationConfig struct {

	// ErrorTopic AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-asyncinferencenotificationconfig.html#cfn-sagemaker-endpointconfig-asyncinferencenotificationconfig-errortopic
	ErrorTopic *types.Value `json:"ErrorTopic,omitempty"`

	// SuccessTopic AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-asyncinferencenotificationconfig.html#cfn-sagemaker-endpointconfig-asyncinferencenotificationconfig-successtopic
	SuccessTopic *types.Value `json:"SuccessTopic,omitempty"`

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
func (r *EndpointConfig_AsyncInferenceNotificationConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::EndpointConfig.AsyncInferenceNotificationConfig"
}
