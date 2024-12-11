package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Flow_S3SourceProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.S3SourceProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3sourceproperties.html
type Flow_S3SourceProperties struct {

	// BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3sourceproperties.html#cfn-appflow-flow-s3sourceproperties-bucketname
	BucketName *types.Value `json:"BucketName,omitempty"`

	// BucketPrefix AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3sourceproperties.html#cfn-appflow-flow-s3sourceproperties-bucketprefix
	BucketPrefix *types.Value `json:"BucketPrefix,omitempty"`

	// S3InputFormatConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3sourceproperties.html#cfn-appflow-flow-s3sourceproperties-s3inputformatconfig
	S3InputFormatConfig *Flow_S3InputFormatConfig `json:"S3InputFormatConfig,omitempty"`

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
func (r *Flow_S3SourceProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.S3SourceProperties"
}
