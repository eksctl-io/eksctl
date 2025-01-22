package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ModelQualityJobDefinition_MonitoringResources AWS CloudFormation Resource (AWS::SageMaker::ModelQualityJobDefinition.MonitoringResources)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-monitoringresources.html
type ModelQualityJobDefinition_MonitoringResources struct {

	// ClusterConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-monitoringresources.html#cfn-sagemaker-modelqualityjobdefinition-monitoringresources-clusterconfig
	ClusterConfig *ModelQualityJobDefinition_ClusterConfig `json:"ClusterConfig,omitempty"`

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
func (r *ModelQualityJobDefinition_MonitoringResources) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelQualityJobDefinition.MonitoringResources"
}
