package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_FlinkApplicationConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.FlinkApplicationConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-flinkapplicationconfiguration.html
type AWSKinesisAnalyticsV2Application_FlinkApplicationConfiguration struct {

	// CheckpointConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-flinkapplicationconfiguration.html#cfn-kinesisanalyticsv2-application-flinkapplicationconfiguration-checkpointconfiguration
	CheckpointConfiguration *AWSKinesisAnalyticsV2Application_CheckpointConfiguration `json:"CheckpointConfiguration,omitempty"`

	// MonitoringConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-flinkapplicationconfiguration.html#cfn-kinesisanalyticsv2-application-flinkapplicationconfiguration-monitoringconfiguration
	MonitoringConfiguration *AWSKinesisAnalyticsV2Application_MonitoringConfiguration `json:"MonitoringConfiguration,omitempty"`

	// ParallelismConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-flinkapplicationconfiguration.html#cfn-kinesisanalyticsv2-application-flinkapplicationconfiguration-parallelismconfiguration
	ParallelismConfiguration *AWSKinesisAnalyticsV2Application_ParallelismConfiguration `json:"ParallelismConfiguration,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_FlinkApplicationConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.FlinkApplicationConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_FlinkApplicationConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
