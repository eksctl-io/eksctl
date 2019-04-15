package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_Input AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.Input)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-input.html
type AWSKinesisAnalyticsV2Application_Input struct {

	// InputParallelism AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-input.html#cfn-kinesisanalyticsv2-application-input-inputparallelism
	InputParallelism *AWSKinesisAnalyticsV2Application_InputParallelism `json:"InputParallelism,omitempty"`

	// InputProcessingConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-input.html#cfn-kinesisanalyticsv2-application-input-inputprocessingconfiguration
	InputProcessingConfiguration *AWSKinesisAnalyticsV2Application_InputProcessingConfiguration `json:"InputProcessingConfiguration,omitempty"`

	// InputSchema AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-input.html#cfn-kinesisanalyticsv2-application-input-inputschema
	InputSchema *AWSKinesisAnalyticsV2Application_InputSchema `json:"InputSchema,omitempty"`

	// KinesisFirehoseInput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-input.html#cfn-kinesisanalyticsv2-application-input-kinesisfirehoseinput
	KinesisFirehoseInput *AWSKinesisAnalyticsV2Application_KinesisFirehoseInput `json:"KinesisFirehoseInput,omitempty"`

	// KinesisStreamsInput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-input.html#cfn-kinesisanalyticsv2-application-input-kinesisstreamsinput
	KinesisStreamsInput *AWSKinesisAnalyticsV2Application_KinesisStreamsInput `json:"KinesisStreamsInput,omitempty"`

	// NamePrefix AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-input.html#cfn-kinesisanalyticsv2-application-input-nameprefix
	NamePrefix *Value `json:"NamePrefix,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_Input) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.Input"
}

func (r *AWSKinesisAnalyticsV2Application_Input) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
