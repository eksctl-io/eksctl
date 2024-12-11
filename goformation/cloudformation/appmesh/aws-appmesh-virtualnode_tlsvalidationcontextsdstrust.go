package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// VirtualNode_TlsValidationContextSdsTrust AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.TlsValidationContextSdsTrust)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tlsvalidationcontextsdstrust.html
type VirtualNode_TlsValidationContextSdsTrust struct {

	// SecretName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tlsvalidationcontextsdstrust.html#cfn-appmesh-virtualnode-tlsvalidationcontextsdstrust-secretname
	SecretName *types.Value `json:"SecretName,omitempty"`

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
func (r *VirtualNode_TlsValidationContextSdsTrust) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.TlsValidationContextSdsTrust"
}
