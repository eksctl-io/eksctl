package cloudformation

import (
	"encoding/json"
)

// AWSIAMUser_LoginProfile AWS CloudFormation Resource (AWS::IAM::User.LoginProfile)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-user-loginprofile.html
type AWSIAMUser_LoginProfile struct {

	// Password AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-user-loginprofile.html#cfn-iam-user-loginprofile-password
	Password *Value `json:"Password,omitempty"`

	// PasswordResetRequired AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-user-loginprofile.html#cfn-iam-user-loginprofile-passwordresetrequired
	PasswordResetRequired *Value `json:"PasswordResetRequired,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIAMUser_LoginProfile) AWSCloudFormationType() string {
	return "AWS::IAM::User.LoginProfile"
}

func (r *AWSIAMUser_LoginProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
