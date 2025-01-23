package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// MonitoringSchedule_MonitoringOutput AWS CloudFormation Resource (AWS::SageMaker::MonitoringSchedule.MonitoringOutput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringoutput.html
type MonitoringSchedule_MonitoringOutput struct {

	// S3Output AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringoutput.html#cfn-sagemaker-monitoringschedule-monitoringoutput-s3output
	S3Output *MonitoringSchedule_S3Output `json:"S3Output,omitempty"`

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
func (r *MonitoringSchedule_MonitoringOutput) AWSCloudFormationType() string {
	return "AWS::SageMaker::MonitoringSchedule.MonitoringOutput"
}
