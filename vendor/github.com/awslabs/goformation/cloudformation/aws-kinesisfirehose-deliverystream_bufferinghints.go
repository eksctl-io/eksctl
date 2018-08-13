package cloudformation

import (
	"encoding/json"
)

// AWSKinesisFirehoseDeliveryStream_BufferingHints AWS CloudFormation Resource (AWS::KinesisFirehose::DeliveryStream.BufferingHints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-bufferinghints.html
type AWSKinesisFirehoseDeliveryStream_BufferingHints struct {

	// IntervalInSeconds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-bufferinghints.html#cfn-kinesisfirehose-deliverystream-bufferinghints-intervalinseconds
	IntervalInSeconds *Value `json:"IntervalInSeconds,omitempty"`

	// SizeInMBs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-bufferinghints.html#cfn-kinesisfirehose-deliverystream-bufferinghints-sizeinmbs
	SizeInMBs *Value `json:"SizeInMBs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisFirehoseDeliveryStream_BufferingHints) AWSCloudFormationType() string {
	return "AWS::KinesisFirehose::DeliveryStream.BufferingHints"
}

func (r *AWSKinesisFirehoseDeliveryStream_BufferingHints) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
