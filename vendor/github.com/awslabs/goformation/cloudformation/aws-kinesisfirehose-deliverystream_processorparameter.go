package cloudformation

// AWSKinesisFirehoseDeliveryStream_ProcessorParameter AWS CloudFormation Resource (AWS::KinesisFirehose::DeliveryStream.ProcessorParameter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-processorparameter.html
type AWSKinesisFirehoseDeliveryStream_ProcessorParameter struct {

	// ParameterName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-processorparameter.html#cfn-kinesisfirehose-deliverystream-processorparameter-parametername
	ParameterName *StringIntrinsic `json:"ParameterName,omitempty"`

	// ParameterValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisfirehose-deliverystream-processorparameter.html#cfn-kinesisfirehose-deliverystream-processorparameter-parametervalue
	ParameterValue *StringIntrinsic `json:"ParameterValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisFirehoseDeliveryStream_ProcessorParameter) AWSCloudFormationType() string {
	return "AWS::KinesisFirehose::DeliveryStream.ProcessorParameter"
}
