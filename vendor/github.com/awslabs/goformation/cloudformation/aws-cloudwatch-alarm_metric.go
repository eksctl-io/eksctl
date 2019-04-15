package cloudformation

import (
	"encoding/json"
)

// AWSCloudWatchAlarm_Metric AWS CloudFormation Resource (AWS::CloudWatch::Alarm.Metric)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metric.html
type AWSCloudWatchAlarm_Metric struct {

	// Dimensions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metric.html#cfn-cloudwatch-alarm-metric-dimensions
	Dimensions []AWSCloudWatchAlarm_Dimension `json:"Dimensions,omitempty"`

	// MetricName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metric.html#cfn-cloudwatch-alarm-metric-metricname
	MetricName *Value `json:"MetricName,omitempty"`

	// Namespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metric.html#cfn-cloudwatch-alarm-metric-namespace
	Namespace *Value `json:"Namespace,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudWatchAlarm_Metric) AWSCloudFormationType() string {
	return "AWS::CloudWatch::Alarm.Metric"
}

func (r *AWSCloudWatchAlarm_Metric) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
