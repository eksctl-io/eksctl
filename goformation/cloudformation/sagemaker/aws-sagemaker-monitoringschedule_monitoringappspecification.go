package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// MonitoringSchedule_MonitoringAppSpecification AWS CloudFormation Resource (AWS::SageMaker::MonitoringSchedule.MonitoringAppSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringappspecification.html
type MonitoringSchedule_MonitoringAppSpecification struct {

	// ContainerArguments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringappspecification.html#cfn-sagemaker-monitoringschedule-monitoringappspecification-containerarguments
	ContainerArguments *types.Value `json:"ContainerArguments,omitempty"`

	// ContainerEntrypoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringappspecification.html#cfn-sagemaker-monitoringschedule-monitoringappspecification-containerentrypoint
	ContainerEntrypoint *types.Value `json:"ContainerEntrypoint,omitempty"`

	// ImageUri AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringappspecification.html#cfn-sagemaker-monitoringschedule-monitoringappspecification-imageuri
	ImageUri *types.Value `json:"ImageUri,omitempty"`

	// PostAnalyticsProcessorSourceUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringappspecification.html#cfn-sagemaker-monitoringschedule-monitoringappspecification-postanalyticsprocessorsourceuri
	PostAnalyticsProcessorSourceUri *types.Value `json:"PostAnalyticsProcessorSourceUri,omitempty"`

	// RecordPreprocessorSourceUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringappspecification.html#cfn-sagemaker-monitoringschedule-monitoringappspecification-recordpreprocessorsourceuri
	RecordPreprocessorSourceUri *types.Value `json:"RecordPreprocessorSourceUri,omitempty"`

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
func (r *MonitoringSchedule_MonitoringAppSpecification) AWSCloudFormationType() string {
	return "AWS::SageMaker::MonitoringSchedule.MonitoringAppSpecification"
}
