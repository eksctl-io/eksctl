package kinesisfirehose

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DeliveryStream_AmazonopensearchserviceBufferingHints AWS CloudFormation Resource (AWS::KinesisFirehose::DeliveryStream.AmazonopensearchserviceBufferingHints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-amazonopensearchservicebufferinghints.html
type DeliveryStream_AmazonopensearchserviceBufferingHints struct {

	// IntervalInSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-amazonopensearchservicebufferinghints.html#cfn-kinesisfirehose-deliverystream-amazonopensearchservicebufferinghints-intervalinseconds
	IntervalInSeconds *types.Value `json:"IntervalInSeconds,omitempty"`

	// SizeInMBs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-amazonopensearchservicebufferinghints.html#cfn-kinesisfirehose-deliverystream-amazonopensearchservicebufferinghints-sizeinmbs
	SizeInMBs *types.Value `json:"SizeInMBs,omitempty"`

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
func (r *DeliveryStream_AmazonopensearchserviceBufferingHints) AWSCloudFormationType() string {
	return "AWS::KinesisFirehose::DeliveryStream.AmazonopensearchserviceBufferingHints"
}
