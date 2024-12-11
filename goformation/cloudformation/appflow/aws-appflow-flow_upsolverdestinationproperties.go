package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Flow_UpsolverDestinationProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.UpsolverDestinationProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-upsolverdestinationproperties.html
type Flow_UpsolverDestinationProperties struct {

	// BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-upsolverdestinationproperties.html#cfn-appflow-flow-upsolverdestinationproperties-bucketname
	BucketName *types.Value `json:"BucketName,omitempty"`

	// BucketPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-upsolverdestinationproperties.html#cfn-appflow-flow-upsolverdestinationproperties-bucketprefix
	BucketPrefix *types.Value `json:"BucketPrefix,omitempty"`

	// S3OutputFormatConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-upsolverdestinationproperties.html#cfn-appflow-flow-upsolverdestinationproperties-s3outputformatconfig
	S3OutputFormatConfig *Flow_UpsolverS3OutputFormatConfig `json:"S3OutputFormatConfig,omitempty"`

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
func (r *Flow_UpsolverDestinationProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.UpsolverDestinationProperties"
}
