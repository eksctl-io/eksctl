package cloudformation

import (
	"encoding/json"
)

// AWSCloudWatchAlarm_MetricDataQuery AWS CloudFormation Resource (AWS::CloudWatch::Alarm.MetricDataQuery)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricdataquery.html
type AWSCloudWatchAlarm_MetricDataQuery struct {

	// Expression AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricdataquery.html#cfn-cloudwatch-alarm-metricdataquery-expression
	Expression *Value `json:"Expression,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricdataquery.html#cfn-cloudwatch-alarm-metricdataquery-id
	Id *Value `json:"Id,omitempty"`

	// Label AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricdataquery.html#cfn-cloudwatch-alarm-metricdataquery-label
	Label *Value `json:"Label,omitempty"`

	// MetricStat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricdataquery.html#cfn-cloudwatch-alarm-metricdataquery-metricstat
	MetricStat *AWSCloudWatchAlarm_MetricStat `json:"MetricStat,omitempty"`

	// ReturnData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-alarm-metricdataquery.html#cfn-cloudwatch-alarm-metricdataquery-returndata
	ReturnData *Value `json:"ReturnData,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudWatchAlarm_MetricDataQuery) AWSCloudFormationType() string {
	return "AWS::CloudWatch::Alarm.MetricDataQuery"
}

func (r *AWSCloudWatchAlarm_MetricDataQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
