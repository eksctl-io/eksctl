package kinesisfirehose

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DeliveryStream_AmazonopensearchserviceRetryOptions AWS CloudFormation Resource (AWS::KinesisFirehose::DeliveryStream.AmazonopensearchserviceRetryOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-amazonopensearchserviceretryoptions.html
type DeliveryStream_AmazonopensearchserviceRetryOptions struct {

	// DurationInSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-amazonopensearchserviceretryoptions.html#cfn-kinesisfirehose-deliverystream-amazonopensearchserviceretryoptions-durationinseconds
	DurationInSeconds *types.Value `json:"DurationInSeconds,omitempty"`

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
func (r *DeliveryStream_AmazonopensearchserviceRetryOptions) AWSCloudFormationType() string {
	return "AWS::KinesisFirehose::DeliveryStream.AmazonopensearchserviceRetryOptions"
}
