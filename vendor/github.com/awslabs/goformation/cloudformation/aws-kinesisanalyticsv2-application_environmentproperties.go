package cloudformation

import (
	"encoding/json"
)

// AWSKinesisAnalyticsV2Application_EnvironmentProperties AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application.EnvironmentProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-environmentproperties.html
type AWSKinesisAnalyticsV2Application_EnvironmentProperties struct {

	// PropertyGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalyticsv2-application-environmentproperties.html#cfn-kinesisanalyticsv2-application-environmentproperties-propertygroups
	PropertyGroups []AWSKinesisAnalyticsV2Application_PropertyGroup `json:"PropertyGroups,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application_EnvironmentProperties) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application.EnvironmentProperties"
}

func (r *AWSKinesisAnalyticsV2Application_EnvironmentProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
