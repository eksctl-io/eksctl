package sso

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// InstanceAccessControlAttributeConfiguration_AccessControlAttributeValue AWS CloudFormation Resource (AWS::SSO::InstanceAccessControlAttributeConfiguration.AccessControlAttributeValue)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sso-instanceaccesscontrolattributeconfiguration-accesscontrolattributevalue.html
type InstanceAccessControlAttributeConfiguration_AccessControlAttributeValue struct {

	// Source AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sso-instanceaccesscontrolattributeconfiguration-accesscontrolattributevalue.html#cfn-sso-instanceaccesscontrolattributeconfiguration-accesscontrolattributevalue-source
	Source *types.Value `json:"Source,omitempty"`

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
func (r *InstanceAccessControlAttributeConfiguration_AccessControlAttributeValue) AWSCloudFormationType() string {
	return "AWS::SSO::InstanceAccessControlAttributeConfiguration.AccessControlAttributeValue"
}
