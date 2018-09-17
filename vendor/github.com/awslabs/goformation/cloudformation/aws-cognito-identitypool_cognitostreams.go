package cloudformation

import (
	"encoding/json"
)

// AWSCognitoIdentityPool_CognitoStreams AWS CloudFormation Resource (AWS::Cognito::IdentityPool.CognitoStreams)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypool-cognitostreams.html
type AWSCognitoIdentityPool_CognitoStreams struct {

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypool-cognitostreams.html#cfn-cognito-identitypool-cognitostreams-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// StreamName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypool-cognitostreams.html#cfn-cognito-identitypool-cognitostreams-streamname
	StreamName *Value `json:"StreamName,omitempty"`

	// StreamingStatus AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypool-cognitostreams.html#cfn-cognito-identitypool-cognitostreams-streamingstatus
	StreamingStatus *Value `json:"StreamingStatus,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCognitoIdentityPool_CognitoStreams) AWSCloudFormationType() string {
	return "AWS::Cognito::IdentityPool.CognitoStreams"
}

func (r *AWSCognitoIdentityPool_CognitoStreams) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
