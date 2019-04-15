package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_MonitoringConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.MonitoringConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-monitoringconfiguration.html
type AWSKinesisAnalyticsV2Application_MonitoringConfiguration struct {

	// ConfigurationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-monitoringconfiguration.html#cfn-kinesisanalyticsv2-application-monitoringconfiguration-configurationtype
	ConfigurationType *Value `json:"ConfigurationType,omitempty"`

	// LogLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-monitoringconfiguration.html#cfn-kinesisanalyticsv2-application-monitoringconfiguration-loglevel
	LogLevel *Value `json:"LogLevel,omitempty"`

	// MetricsLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-monitoringconfiguration.html#cfn-kinesisanalyticsv2-application-monitoringconfiguration-metricslevel
	MetricsLevel *Value `json:"MetricsLevel,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_MonitoringConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.MonitoringConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_MonitoringConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
