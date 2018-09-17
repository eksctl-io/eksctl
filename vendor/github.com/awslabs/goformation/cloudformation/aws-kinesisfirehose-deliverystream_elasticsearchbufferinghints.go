package cloudformation

import (
	"encoding/json"
)

// AWSKinesisFirehoseDeliveryStream_ElasticsearchBufferingHints AWS CloudFormation Resource (AWS::KinesisFirehose::DeliveryStream.ElasticsearchBufferingHints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-elasticsearchbufferinghints.html
type AWSKinesisFirehoseDeliveryStream_ElasticsearchBufferingHints struct {

	// IntervalInSeconds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-elasticsearchbufferinghints.html#cfn-kinesisfirehose-deliverystream-elasticsearchbufferinghints-intervalinseconds
	IntervalInSeconds *Value `json:"IntervalInSeconds,omitempty"`

	// SizeInMBs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-elasticsearchbufferinghints.html#cfn-kinesisfirehose-deliverystream-elasticsearchbufferinghints-sizeinmbs
	SizeInMBs *Value `json:"SizeInMBs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisFirehoseDeliveryStream_ElasticsearchBufferingHints) AWSCloudFormationType() string {
	return "AWS::KinesisFirehose::DeliveryStream.ElasticsearchBufferingHints"
}

func (r *AWSKinesisFirehoseDeliveryStream_ElasticsearchBufferingHints) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
