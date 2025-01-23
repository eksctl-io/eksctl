package appmesh

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// VirtualGateway_VirtualGatewayTlsValidationContext AWS CloudFormation Resource (AWS::AppMesh::VirtualGateway.VirtualGatewayTlsValidationContext)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewaytlsvalidationcontext.html
type VirtualGateway_VirtualGatewayTlsValidationContext struct {

	// SubjectAlternativeNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewaytlsvalidationcontext.html#cfn-appmesh-virtualgateway-virtualgatewaytlsvalidationcontext-subjectalternativenames
	SubjectAlternativeNames *VirtualGateway_SubjectAlternativeNames `json:"SubjectAlternativeNames,omitempty"`

	// Trust AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewaytlsvalidationcontext.html#cfn-appmesh-virtualgateway-virtualgatewaytlsvalidationcontext-trust
	Trust *VirtualGateway_VirtualGatewayTlsValidationContextTrust `json:"Trust,omitempty"`

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
func (r *VirtualGateway_VirtualGatewayTlsValidationContext) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualGateway.VirtualGatewayTlsValidationContext"
}
