package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_ApplicationSnapshotConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.ApplicationSnapshotConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationsnapshotconfiguration.html
type AWSKinesisAnalyticsV2Application_ApplicationSnapshotConfiguration struct {

	// SnapshotsEnabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationsnapshotconfiguration.html#cfn-kinesisanalyticsv2-application-applicationsnapshotconfiguration-snapshotsenabled
	SnapshotsEnabled *Value `json:"SnapshotsEnabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_ApplicationSnapshotConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.ApplicationSnapshotConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_ApplicationSnapshotConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
