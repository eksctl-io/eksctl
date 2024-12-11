package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ModelQualityJobDefinition_NetworkConfig AWS CloudFormation Resource (AWS::SageMaker::ModelQualityJobDefinition.NetworkConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-networkconfig.html
type ModelQualityJobDefinition_NetworkConfig struct {

	// EnableInterContainerTrafficEncryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-networkconfig.html#cfn-sagemaker-modelqualityjobdefinition-networkconfig-enableintercontainertrafficencryption
	EnableInterContainerTrafficEncryption *types.Value `json:"EnableInterContainerTrafficEncryption,omitempty"`

	// EnableNetworkIsolation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-networkconfig.html#cfn-sagemaker-modelqualityjobdefinition-networkconfig-enablenetworkisolation
	EnableNetworkIsolation *types.Value `json:"EnableNetworkIsolation,omitempty"`

	// VpcConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-networkconfig.html#cfn-sagemaker-modelqualityjobdefinition-networkconfig-vpcconfig
	VpcConfig *ModelQualityJobDefinition_VpcConfig `json:"VpcConfig,omitempty"`

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
func (r *ModelQualityJobDefinition_NetworkConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelQualityJobDefinition.NetworkConfig"
}
