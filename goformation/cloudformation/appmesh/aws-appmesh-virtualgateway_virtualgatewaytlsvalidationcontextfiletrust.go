package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// VirtualGateway_VirtualGatewayTlsValidationContextFileTrust AWS CloudFormation Resource (AWS::AppMesh::VirtualGateway.VirtualGatewayTlsValidationContextFileTrust)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewaytlsvalidationcontextfiletrust.html
type VirtualGateway_VirtualGatewayTlsValidationContextFileTrust struct {

	// CertificateChain AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewaytlsvalidationcontextfiletrust.html#cfn-appmesh-virtualgateway-virtualgatewaytlsvalidationcontextfiletrust-certificatechain
	CertificateChain *types.Value `json:"CertificateChain,omitempty"`

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
func (r *VirtualGateway_VirtualGatewayTlsValidationContextFileTrust) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualGateway.VirtualGatewayTlsValidationContextFileTrust"
}
