package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_ParallelismConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.ParallelismConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-parallelismconfiguration.html
type AWSKinesisAnalyticsV2Application_ParallelismConfiguration struct {

	// AutoScalingEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-parallelismconfiguration.html#cfn-kinesisanalyticsv2-application-parallelismconfiguration-autoscalingenabled
	AutoScalingEnabled *Value `json:"AutoScalingEnabled,omitempty"`

	// ConfigurationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-parallelismconfiguration.html#cfn-kinesisanalyticsv2-application-parallelismconfiguration-configurationtype
	ConfigurationType *Value `json:"ConfigurationType,omitempty"`

	// Parallelism AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-parallelismconfiguration.html#cfn-kinesisanalyticsv2-application-parallelismconfiguration-parallelism
	Parallelism *Value `json:"Parallelism,omitempty"`

	// ParallelismPerKPU AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-parallelismconfiguration.html#cfn-kinesisanalyticsv2-application-parallelismconfiguration-parallelismperkpu
	ParallelismPerKPU *Value `json:"ParallelismPerKPU,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_ParallelismConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.ParallelismConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_ParallelismConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
