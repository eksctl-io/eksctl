package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamStack_ApplicationSettings AWS CloudFormation Resource (AWS::AppStream::Stack.ApplicationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-applicationsettings.html
type AWSAppStreamStack_ApplicationSettings struct {

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-applicationsettings.html#cfn-appstream-stack-applicationsettings-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// SettingsGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-applicationsettings.html#cfn-appstream-stack-applicationsettings-settingsgroup
	SettingsGroup *Value `json:"SettingsGroup,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamStack_ApplicationSettings) AWSCloudFormationType() string {
	return "AWS::AppStream::Stack.ApplicationSettings"
}

func (r *AWSAppStreamStack_ApplicationSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
