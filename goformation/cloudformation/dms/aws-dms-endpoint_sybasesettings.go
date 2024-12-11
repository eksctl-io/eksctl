package dms

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Endpoint_SybaseSettings AWS CloudFormation Resource (AWS::DMS::Endpoint.SybaseSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-sybasesettings.html
type Endpoint_SybaseSettings struct {

	// SecretsManagerAccessRoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-sybasesettings.html#cfn-dms-endpoint-sybasesettings-secretsmanageraccessrolearn
	SecretsManagerAccessRoleArn *types.Value `json:"SecretsManagerAccessRoleArn,omitempty"`

	// SecretsManagerSecretId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-sybasesettings.html#cfn-dms-endpoint-sybasesettings-secretsmanagersecretid
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
func (r *Endpoint_SybaseSettings) AWSCloudFormationType() string {
	return "AWS::DMS::Endpoint.SybaseSettings"
}
