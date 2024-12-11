package cognito

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// UserPoolClient_TokenValidityUnits AWS CloudFormation Resource (AWS::Cognito::UserPoolClient.TokenValidityUnits)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpoolclient-tokenvalidityunits.html
type UserPoolClient_TokenValidityUnits struct {

	// AccessToken AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpoolclient-tokenvalidityunits.html#cfn-cognito-userpoolclient-tokenvalidityunits-accesstoken
	AccessToken *types.Value `json:"AccessToken,omitempty"`

	// IdToken AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpoolclient-tokenvalidityunits.html#cfn-cognito-userpoolclient-tokenvalidityunits-idtoken
	IdToken *types.Value `json:"IdToken,omitempty"`

	// RefreshToken AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpoolclient-tokenvalidityunits.html#cfn-cognito-userpoolclient-tokenvalidityunits-refreshtoken
	RefreshToken *types.Value `json:"RefreshToken,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *UserPoolClient_TokenValidityUnits) AWSCloudFormationType() string {
	return "AWS::Cognito::UserPoolClient.TokenValidityUnits"
}
