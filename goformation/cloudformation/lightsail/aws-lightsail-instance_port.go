package lightsail

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Instance_Port AWS CloudFormation Resource (AWS::Lightsail::Instance.Port)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html
type Instance_Port struct {

	// AccessDirection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-accessdirection
	AccessDirection *types.Value `json:"AccessDirection,omitempty"`

	// AccessFrom AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-accessfrom
	AccessFrom *types.Value `json:"AccessFrom,omitempty"`

	// AccessType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-accesstype
	AccessType *types.Value `json:"AccessType,omitempty"`

	// CidrListAliases AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-cidrlistaliases
	CidrListAliases *types.Value `json:"CidrListAliases,omitempty"`

	// Cidrs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-cidrs
	Cidrs *types.Value `json:"Cidrs,omitempty"`

	// CommonName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-commonname
	CommonName *types.Value `json:"CommonName,omitempty"`

	// FromPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-fromport
	FromPort *types.Value `json:"FromPort,omitempty"`

	// Ipv6Cidrs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-ipv6cidrs
	Ipv6Cidrs *types.Value `json:"Ipv6Cidrs,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-protocol
	Protocol *types.Value `json:"Protocol,omitempty"`

	// ToPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-port.html#cfn-lightsail-instance-port-toport
	ToPort *types.Value `json:"ToPort,omitempty"`

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
func (r *Instance_Port) AWSCloudFormationType() string {
	return "AWS::Lightsail::Instance.Port"
}
