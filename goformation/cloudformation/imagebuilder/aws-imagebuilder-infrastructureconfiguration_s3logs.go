package imagebuilder

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// InfrastructureConfiguration_S3Logs AWS CloudFormation Resource (AWS::ImageBuilder::InfrastructureConfiguration.S3Logs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-infrastructureconfiguration-s3logs.html
type InfrastructureConfiguration_S3Logs struct {

	// S3BucketName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-infrastructureconfiguration-s3logs.html#cfn-imagebuilder-infrastructureconfiguration-s3logs-s3bucketname
	S3BucketName *types.Value `json:"S3BucketName,omitempty"`

	// S3KeyPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-infrastructureconfiguration-s3logs.html#cfn-imagebuilder-infrastructureconfiguration-s3logs-s3keyprefix
	S3KeyPrefix *types.Value `json:"S3KeyPrefix,omitempty"`

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
func (r *InfrastructureConfiguration_S3Logs) AWSCloudFormationType() string {
	return "AWS::ImageBuilder::InfrastructureConfiguration.S3Logs"
}
