package dms

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Endpoint_DocDbSettings AWS CloudFormation Resource (AWS::DMS::Endpoint.DocDbSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-docdbsettings.html
type Endpoint_DocDbSettings struct {

	// SecretsManagerAccessRoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-docdbsettings.html#cfn-dms-endpoint-docdbsettings-secretsmanageraccessrolearn
	SecretsManagerAccessRoleArn *types.Value `json:"SecretsManagerAccessRoleArn,omitempty"`

	// SecretsManagerSecretId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-docdbsettings.html#cfn-dms-endpoint-docdbsettings-secretsmanagersecretid
	SecretsManagerSecretId *types.Value `json:"SecretsManagerSecretId,omitempty"`

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
func (r *Endpoint_DocDbSettings) AWSCloudFormationType() string {
	return "AWS::DMS::Endpoint.DocDbSettings"
}
