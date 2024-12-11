package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// VirtualNode_ListenerTlsFileCertificate AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.ListenerTlsFileCertificate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-listenertlsfilecertificate.html
type VirtualNode_ListenerTlsFileCertificate struct {

	// CertificateChain AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-listenertlsfilecertificate.html#cfn-appmesh-virtualnode-listenertlsfilecertificate-certificatechain
	CertificateChain *types.Value `json:"CertificateChain,omitempty"`

	// PrivateKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-listenertlsfilecertificate.html#cfn-appmesh-virtualnode-listenertlsfilecertificate-privatekey
	PrivateKey *types.Value `json:"PrivateKey,omitempty"`

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
func (r *VirtualNode_ListenerTlsFileCertificate) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.ListenerTlsFileCertificate"
}
