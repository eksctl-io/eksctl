package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2ApplicationOutput_Output AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationOutput.Output)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-output.html
type AWSKinesisAnalyticsV2ApplicationOutput_Output struct {

	// DestinationSchema AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-output.html#cfn-kinesisanalyticsv2-applicationoutput-output-destinationschema
	DestinationSchema *AWSKinesisAnalyticsV2ApplicationOutput_DestinationSchema `json:"DestinationSchema,omitempty"`

	// KinesisFirehoseOutput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-output.html#cfn-kinesisanalyticsv2-applicationoutput-output-kinesisfirehoseoutput
	KinesisFirehoseOutput *AWSKinesisAnalyticsV2ApplicationOutput_KinesisFirehoseOutput `json:"KinesisFirehoseOutput,omitempty"`

	// KinesisStreamsOutput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-output.html#cfn-kinesisanalyticsv2-applicationoutput-output-kinesisstreamsoutput
	KinesisStreamsOutput *AWSKinesisAnalyticsV2ApplicationOutput_KinesisStreamsOutput `json:"KinesisStreamsOutput,omitempty"`

	// LambdaOutput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-output.html#cfn-kinesisanalyticsv2-applicationoutput-output-lambdaoutput
	LambdaOutput *AWSKinesisAnalyticsV2ApplicationOutput_LambdaOutput `json:"LambdaOutput,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-applicationoutput-output.html#cfn-kinesisanalyticsv2-applicationoutput-output-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationOutput_Output) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationOutput.Output"
}

func (r *AWSKinesisAnalyticsV2ApplicationOutput_Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
