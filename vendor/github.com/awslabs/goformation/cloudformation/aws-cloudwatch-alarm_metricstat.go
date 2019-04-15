package cloudformation

import (
	"encoding/json"
)

// AWSCloudWatchAlarm_MetricStat AWS CloudFormation Resource (AWS::CloudWatch::Alarm.MetricStat)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricstat.html
type AWSCloudWatchAlarm_MetricStat struct {

	// Metric AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricstat.html#cfn-cloudwatch-alarm-metricstat-metric
	Metric *AWSCloudWatchAlarm_Metric `json:"Metric,omitempty"`

	// Period AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricstat.html#cfn-cloudwatch-alarm-metricstat-period
	Period *Value `json:"Period,omitempty"`

	// Stat AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricstat.html#cfn-cloudwatch-alarm-metricstat-stat
	Stat *Value `json:"Stat,omitempty"`

	// Unit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricstat.html#cfn-cloudwatch-alarm-metricstat-unit
	Unit *Value `json:"Unit,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudWatchAlarm_MetricStat) AWSCloudFormationType() string {
	return "AWS::CloudWatch::Alarm.MetricStat"
}

func (r *AWSCloudWatchAlarm_MetricStat) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
