package cloudformation

import (
	"encoding/json"
)

// AWSCognitoIdentityPoolRoleAttachment_RoleMapping AWS CloudFormation Resource (AWS::Cognito::IdentityPoolRoleAttachment.RoleMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypoolroleattachment-rolemapping.html
type AWSCognitoIdentityPoolRoleAttachment_RoleMapping struct {

	// AmbiguousRoleResolution AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypoolroleattachment-rolemapping.html#cfn-cognito-identitypoolroleattachment-rolemapping-ambiguousroleresolution
	AmbiguousRoleResolution *Value `json:"AmbiguousRoleResolution,omitempty"`

	// RulesConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypoolroleattachment-rolemapping.html#cfn-cognito-identitypoolroleattachment-rolemapping-rulesconfiguration
	RulesConfiguration *AWSCognitoIdentityPoolRoleAttachment_RulesConfigurationType `json:"RulesConfiguration,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-identitypoolroleattachment-rolemapping.html#cfn-cognito-identitypoolroleattachment-rolemapping-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCognitoIdentityPoolRoleAttachment_RoleMapping) AWSCloudFormationType() string {
	return "AWS::Cognito::IdentityPoolRoleAttachment.RoleMapping"
}

func (r *AWSCognitoIdentityPoolRoleAttachment_RoleMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
