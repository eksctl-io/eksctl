package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Flow_S3DestinationProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.S3DestinationProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3destinationproperties.html
type Flow_S3DestinationProperties struct {

	// BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3destinationproperties.html#cfn-appflow-flow-s3destinationproperties-bucketname
	BucketName *types.Value `json:"BucketName,omitempty"`

	// BucketPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3destinationproperties.html#cfn-appflow-flow-s3destinationproperties-bucketprefix
	BucketPrefix *types.Value `json:"BucketPrefix,omitempty"`

	// S3OutputFormatConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-s3destinationproperties.html#cfn-appflow-flow-s3destinationproperties-s3outputformatconfig
	S3OutputFormatConfig *Flow_S3OutputFormatConfig `json:"S3OutputFormatConfig,omitempty"`

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
func (r *Flow_S3DestinationProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.S3DestinationProperties"
}
