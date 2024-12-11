package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataQualityJobDefinition_DataQualityBaselineConfig AWS CloudFormation Resource (AWS::SageMaker::DataQualityJobDefinition.DataQualityBaselineConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-dataqualityjobdefinition-dataqualitybaselineconfig.html
type DataQualityJobDefinition_DataQualityBaselineConfig struct {

	// BaseliningJobName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-dataqualityjobdefinition-dataqualitybaselineconfig.html#cfn-sagemaker-dataqualityjobdefinition-dataqualitybaselineconfig-baseliningjobname
	BaseliningJobName *types.Value `json:"BaseliningJobName,omitempty"`

	// ConstraintsResource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-dataqualityjobdefinition-dataqualitybaselineconfig.html#cfn-sagemaker-dataqualityjobdefinition-dataqualitybaselineconfig-constraintsresource
	ConstraintsResource *DataQualityJobDefinition_ConstraintsResource `json:"ConstraintsResource,omitempty"`

	// StatisticsResource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-dataqualityjobdefinition-dataqualitybaselineconfig.html#cfn-sagemaker-dataqualityjobdefinition-dataqualitybaselineconfig-statisticsresource
	StatisticsResource *DataQualityJobDefinition_StatisticsResource `json:"StatisticsResource,omitempty"`

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
func (r *DataQualityJobDefinition_DataQualityBaselineConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::DataQualityJobDefinition.DataQualityBaselineConfig"
}
