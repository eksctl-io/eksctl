package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// MonitoringSchedule_S3Output AWS CloudFormation Resource (AWS::SageMaker::MonitoringSchedule.S3Output)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-s3output.html
type MonitoringSchedule_S3Output struct {

	// LocalPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-s3output.html#cfn-sagemaker-monitoringschedule-s3output-localpath
	LocalPath *types.Value `json:"LocalPath,omitempty"`

	// S3UploadMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-s3output.html#cfn-sagemaker-monitoringschedule-s3output-s3uploadmode
	S3UploadMode *types.Value `json:"S3UploadMode,omitempty"`

	// S3Uri AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-monitoringschedule-s3output.html#cfn-sagemaker-monitoringschedule-s3output-s3uri
	S3Uri *types.Value `json:"S3Uri,omitempty"`

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
func (r *MonitoringSchedule_S3Output) AWSCloudFormationType() string {
	return "AWS::SageMaker::MonitoringSchedule.S3Output"
}
