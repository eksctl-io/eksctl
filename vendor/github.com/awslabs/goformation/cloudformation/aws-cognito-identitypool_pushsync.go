package cloudformation

import (
	"encoding/json"
)

// AWSCognitoIdentityPool_PushSync AWS CloudFormation Resource (AWS::Cognito::IdentityPool.PushSync)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypool-pushsync.html
type AWSCognitoIdentityPool_PushSync struct {

	// ApplicationArns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypool-pushsync.html#cfn-cognito-identitypool-pushsync-applicationarns
	ApplicationArns []*Value `json:"ApplicationArns,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypool-pushsync.html#cfn-cognito-identitypool-pushsync-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCognitoIdentityPool_PushSync) AWSCloudFormationType() string {
	return "AWS::Cognito::IdentityPool.PushSync"
}

func (r *AWSCognitoIdentityPool_PushSync) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
