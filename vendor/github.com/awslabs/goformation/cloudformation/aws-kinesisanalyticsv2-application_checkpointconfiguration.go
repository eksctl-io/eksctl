package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_CheckpointConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.CheckpointConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-checkpointconfiguration.html
type AWSKinesisAnalyticsV2Application_CheckpointConfiguration struct {

	// CheckpointInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-checkpointconfiguration.html#cfn-kinesisanalyticsv2-application-checkpointconfiguration-checkpointinterval
	CheckpointInterval *Value `json:"CheckpointInterval,omitempty"`

	// CheckpointingEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-checkpointconfiguration.html#cfn-kinesisanalyticsv2-application-checkpointconfiguration-checkpointingenabled
	CheckpointingEnabled *Value `json:"CheckpointingEnabled,omitempty"`

	// ConfigurationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-checkpointconfiguration.html#cfn-kinesisanalyticsv2-application-checkpointconfiguration-configurationtype
	ConfigurationType *Value `json:"ConfigurationType,omitempty"`

	// MinPauseBetweenCheckpoints AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-checkpointconfiguration.html#cfn-kinesisanalyticsv2-application-checkpointconfiguration-minpausebetweencheckpoints
	MinPauseBetweenCheckpoints *Value `json:"MinPauseBetweenCheckpoints,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_CheckpointConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.CheckpointConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_CheckpointConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
