package dms

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Endpoint_RedisSettings AWS CloudFormation Resource (AWS::DMS::Endpoint.RedisSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html
type Endpoint_RedisSettings struct {

	// AuthPassword AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html#cfn-dms-endpoint-redissettings-authpassword
	AuthPassword *types.Value `json:"AuthPassword,omitempty"`

	// AuthType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html#cfn-dms-endpoint-redissettings-authtype
	AuthType *types.Value `json:"AuthType,omitempty"`

	// AuthUserName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html#cfn-dms-endpoint-redissettings-authusername
	AuthUserName *types.Value `json:"AuthUserName,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html#cfn-dms-endpoint-redissettings-port
	Port *types.Value `json:"Port,omitempty"`

	// ServerName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html#cfn-dms-endpoint-redissettings-servername
	ServerName *types.Value `json:"ServerName,omitempty"`

	// SslCaCertificateArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html#cfn-dms-endpoint-redissettings-sslcacertificatearn
	SslCaCertificateArn *types.Value `json:"SslCaCertificateArn,omitempty"`

	// SslSecurityProtocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-redissettings.html#cfn-dms-endpoint-redissettings-sslsecurityprotocol
	SslSecurityProtocol *types.Value `json:"SslSecurityProtocol,omitempty"`

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
func (r *Endpoint_RedisSettings) AWSCloudFormationType() string {
	return "AWS::DMS::Endpoint.RedisSettings"
}
