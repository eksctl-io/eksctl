package cloudformation

// AWSKinesisFirehoseDeliveryStream_SplunkRetryOptions AWS CloudFormation Resource (AWS::KinesisFirehose::DeliveryStream.SplunkRetryOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-splunkretryoptions.html
type AWSKinesisFirehoseDeliveryStream_SplunkRetryOptions struct {

	// DurationInSeconds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-splunkretryoptions.html#cfn-kinesisfirehose-deliverystream-splunkretryoptions-durationinseconds
	DurationInSeconds int `json:"DurationInSeconds,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisFirehoseDeliveryStream_SplunkRetryOptions) AWSCloudFormationType() string {
	return "AWS::KinesisFirehose::DeliveryStream.SplunkRetryOptions"
}
