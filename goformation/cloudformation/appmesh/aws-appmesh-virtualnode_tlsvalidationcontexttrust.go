package appmesh

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// VirtualNode_TlsValidationContextTrust AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.TlsValidationContextTrust)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tlsvalidationcontexttrust.html
type VirtualNode_TlsValidationContextTrust struct {

	// ACM AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tlsvalidationcontexttrust.html#cfn-appmesh-virtualnode-tlsvalidationcontexttrust-acm
	ACM *VirtualNode_TlsValidationContextAcmTrust `json:"ACM,omitempty"`

	// File AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tlsvalidationcontexttrust.html#cfn-appmesh-virtualnode-tlsvalidationcontexttrust-file
	File *VirtualNode_TlsValidationContextFileTrust `json:"File,omitempty"`

	// SDS AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tlsvalidationcontexttrust.html#cfn-appmesh-virtualnode-tlsvalidationcontexttrust-sds
	SDS *VirtualNode_TlsValidationContextSdsTrust `json:"SDS,omitempty"`

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
func (r *VirtualNode_TlsValidationContextTrust) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.TlsValidationContextTrust"
}
