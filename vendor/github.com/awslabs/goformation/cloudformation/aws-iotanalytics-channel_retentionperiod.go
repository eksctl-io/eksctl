package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsChannel_RetentionPeriod AWS CloudFormation Resource (AWS::IoTAnalytics::Channel.RetentionPeriod)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-channel-retentionperiod.html
type AWSIoTAnalyticsChannel_RetentionPeriod struct {

	// NumberOfDays AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-channel-retentionperiod.html#cfn-iotanalytics-channel-retentionperiod-numberofdays
	NumberOfDays *Value `json:"NumberOfDays,omitempty"`

	// Unlimited AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-channel-retentionperiod.html#cfn-iotanalytics-channel-retentionperiod-unlimited
	Unlimited *Value `json:"Unlimited,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsChannel_RetentionPeriod) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Channel.RetentionPeriod"
}

func (r *AWSIoTAnalyticsChannel_RetentionPeriod) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
