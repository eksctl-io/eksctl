package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_CloudwatchMetricAction AWS CloudFormation Resource (AWS::IoT::TopicRule.CloudwatchMetricAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-cloudwatchmetricaction.html
type AWSIoTTopicRule_CloudwatchMetricAction struct {

	// MetricName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-cloudwatchmetricaction.html#cfn-iot-topicrule-cloudwatchmetricaction-metricname
	MetricName *Value `json:"MetricName,omitempty"`

	// MetricNamespace AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-cloudwatchmetricaction.html#cfn-iot-topicrule-cloudwatchmetricaction-metricnamespace
	MetricNamespace *Value `json:"MetricNamespace,omitempty"`

	// MetricTimestamp AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-cloudwatchmetricaction.html#cfn-iot-topicrule-cloudwatchmetricaction-metrictimestamp
	MetricTimestamp *Value `json:"MetricTimestamp,omitempty"`

	// MetricUnit AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-cloudwatchmetricaction.html#cfn-iot-topicrule-cloudwatchmetricaction-metricunit
	MetricUnit *Value `json:"MetricUnit,omitempty"`

	// MetricValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-cloudwatchmetricaction.html#cfn-iot-topicrule-cloudwatchmetricaction-metricvalue
	MetricValue *Value `json:"MetricValue,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-cloudwatchmetricaction.html#cfn-iot-topicrule-cloudwatchmetricaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_CloudwatchMetricAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.CloudwatchMetricAction"
}

func (r *AWSIoTTopicRule_CloudwatchMetricAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
