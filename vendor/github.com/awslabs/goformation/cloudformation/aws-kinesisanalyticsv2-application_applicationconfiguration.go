package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_ApplicationConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.ApplicationConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationconfiguration.html
type AWSKinesisAnalyticsV2Application_ApplicationConfiguration struct {

	// ApplicationCodeConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationconfiguration.html#cfn-kinesisanalyticsv2-application-applicationconfiguration-applicationcodeconfiguration
	ApplicationCodeConfiguration *AWSKinesisAnalyticsV2Application_ApplicationCodeConfiguration `json:"ApplicationCodeConfiguration,omitempty"`

	// ApplicationSnapshotConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationconfiguration.html#cfn-kinesisanalyticsv2-application-applicationconfiguration-applicationsnapshotconfiguration
	ApplicationSnapshotConfiguration *AWSKinesisAnalyticsV2Application_ApplicationSnapshotConfiguration `json:"ApplicationSnapshotConfiguration,omitempty"`

	// EnvironmentProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationconfiguration.html#cfn-kinesisanalyticsv2-application-applicationconfiguration-environmentproperties
	EnvironmentProperties *AWSKinesisAnalyticsV2Application_EnvironmentProperties `json:"EnvironmentProperties,omitempty"`

	// FlinkApplicationConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationconfiguration.html#cfn-kinesisanalyticsv2-application-applicationconfiguration-flinkapplicationconfiguration
	FlinkApplicationConfiguration *AWSKinesisAnalyticsV2Application_FlinkApplicationConfiguration `json:"FlinkApplicationConfiguration,omitempty"`

	// SqlApplicationConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationconfiguration.html#cfn-kinesisanalyticsv2-application-applicationconfiguration-sqlapplicationconfiguration
	SqlApplicationConfiguration *AWSKinesisAnalyticsV2Application_SqlApplicationConfiguration `json:"SqlApplicationConfiguration,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_ApplicationConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.ApplicationConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_ApplicationConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
