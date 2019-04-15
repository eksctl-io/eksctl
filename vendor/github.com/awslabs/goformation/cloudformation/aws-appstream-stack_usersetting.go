package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamStack_UserSetting AWS CloudFormation Resource (AWS::AppStream::Stack.UserSetting)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-usersetting.html
type AWSAppStreamStack_UserSetting struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-usersetting.html#cfn-appstream-stack-usersetting-action
	Action *Value `json:"Action,omitempty"`

	// Permission AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-stack-usersetting.html#cfn-appstream-stack-usersetting-permission
	Permission *Value `json:"Permission,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamStack_UserSetting) AWSCloudFormationType() string {
	return "AWS::AppStream::Stack.UserSetting"
}

func (r *AWSAppStreamStack_UserSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
