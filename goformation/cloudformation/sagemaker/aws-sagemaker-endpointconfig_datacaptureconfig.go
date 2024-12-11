package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// EndpointConfig_DataCaptureConfig AWS CloudFormation Resource (AWS::SageMaker::EndpointConfig.DataCaptureConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-datacaptureconfig.html
type EndpointConfig_DataCaptureConfig struct {

	// CaptureContentTypeHeader AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-datacaptureconfig.html#cfn-sagemaker-endpointconfig-datacaptureconfig-capturecontenttypeheader
	CaptureContentTypeHeader *EndpointConfig_CaptureContentTypeHeader `json:"CaptureContentTypeHeader,omitempty"`

	// CaptureOptions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-datacaptureconfig.html#cfn-sagemaker-endpointconfig-datacaptureconfig-captureoptions
	CaptureOptions []EndpointConfig_CaptureOption `json:"CaptureOptions,omitempty"`

	// DestinationS3Uri AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-datacaptureconfig.html#cfn-sagemaker-endpointconfig-datacaptureconfig-destinations3uri
	DestinationS3Uri *types.Value `json:"DestinationS3Uri,omitempty"`

	// EnableCapture AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-datacaptureconfig.html#cfn-sagemaker-endpointconfig-datacaptureconfig-enablecapture
	EnableCapture *types.Value `json:"EnableCapture,omitempty"`

	// InitialSamplingPercentage AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-datacaptureconfig.html#cfn-sagemaker-endpointconfig-datacaptureconfig-initialsamplingpercentage
	InitialSamplingPercentage *types.Value `json:"InitialSamplingPercentage"`

	// KmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpointconfig-datacaptureconfig.html#cfn-sagemaker-endpointconfig-datacaptureconfig-kmskeyid
	KmsKeyId *types.Value `json:"KmsKeyId,omitempty"`

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
func (r *EndpointConfig_DataCaptureConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::EndpointConfig.DataCaptureConfig"
}
