package appmesh

import (
	"github.com/awslabs/goformation/v4/cloudformation/types"

	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// VirtualNode_VirtualNodeTcpConnectionPool AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.VirtualNodeTcpConnectionPool)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualnodetcpconnectionpool.html
type VirtualNode_VirtualNodeTcpConnectionPool struct {

	// MaxConnections AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualnodetcpconnectionpool.html#cfn-appmesh-virtualnode-virtualnodetcpconnectionpool-maxconnections
	MaxConnections *types.Value `json:"MaxConnections"`

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
func (r *VirtualNode_VirtualNodeTcpConnectionPool) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.VirtualNodeTcpConnectionPool"
}
