package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDatastore_RetentionPeriod AWS CloudFormation Resource (AWS::IoTAnalytics::Datastore.RetentionPeriod)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-retentionperiod.html
type AWSIoTAnalyticsDatastore_RetentionPeriod struct {

	// NumberOfDays AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-retentionperiod.html#cfn-iotanalytics-datastore-retentionperiod-numberofdays
	NumberOfDays *Value `json:"NumberOfDays,omitempty"`

	// Unlimited AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-retentionperiod.html#cfn-iotanalytics-datastore-retentionperiod-unlimited
	Unlimited *Value `json:"Unlimited,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDatastore_RetentionPeriod) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Datastore.RetentionPeriod"
}

func (r *AWSIoTAnalyticsDatastore_RetentionPeriod) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
