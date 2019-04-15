package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_ApplicationCodeConfiguration AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.ApplicationCodeConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationcodeconfiguration.html
type AWSKinesisAnalyticsV2Application_ApplicationCodeConfiguration struct {

	// CodeContent AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationcodeconfiguration.html#cfn-kinesisanalyticsv2-application-applicationcodeconfiguration-codecontent
	CodeContent *AWSKinesisAnalyticsV2Application_CodeContent `json:"CodeContent,omitempty"`

	// CodeContentType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-applicationcodeconfiguration.html#cfn-kinesisanalyticsv2-application-applicationcodeconfiguration-codecontenttype
	CodeContentType *Value `json:"CodeContentType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_ApplicationCodeConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.ApplicationCodeConfiguration"
}

func (r *AWSKinesisAnalyticsV2Application_ApplicationCodeConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
