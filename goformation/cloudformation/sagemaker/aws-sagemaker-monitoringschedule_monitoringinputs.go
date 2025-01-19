package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// MonitoringSchedule_MonitoringInputs AWS CloudFormation Resource (AWS::SageMaker::MonitoringSchedule.MonitoringInputs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringinputs.html
type MonitoringSchedule_MonitoringInputs struct {

	// MonitoringInputs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-monitoringinputs.html#cfn-sagemaker-monitoringschedule-monitoringinputs-monitoringinputs
	MonitoringInputs []MonitoringSchedule_MonitoringInput `json:"MonitoringInputs,omitempty"`

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
func (r *MonitoringSchedule_MonitoringInputs) AWSCloudFormationType() string {
	return "AWS::SageMaker::MonitoringSchedule.MonitoringInputs"
}
