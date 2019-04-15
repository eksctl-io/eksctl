package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_VirtualNodeSpec AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.VirtualNodeSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualnodespec.html
type AWSAppMeshVirtualNode_VirtualNodeSpec struct {

	// Backends AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualnodespec.html#cfn-appmesh-virtualnode-virtualnodespec-backends
	Backends []AWSAppMeshVirtualNode_Backend `json:"Backends,omitempty"`

	// Listeners AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualnodespec.html#cfn-appmesh-virtualnode-virtualnodespec-listeners
	Listeners []AWSAppMeshVirtualNode_Listener `json:"Listeners,omitempty"`

	// Logging AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualnodespec.html#cfn-appmesh-virtualnode-virtualnodespec-logging
	Logging *AWSAppMeshVirtualNode_Logging `json:"Logging,omitempty"`

	// ServiceDiscovery AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualnodespec.html#cfn-appmesh-virtualnode-virtualnodespec-servicediscovery
	ServiceDiscovery *AWSAppMeshVirtualNode_ServiceDiscovery `json:"ServiceDiscovery,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_VirtualNodeSpec) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.VirtualNodeSpec"
}

func (r *AWSAppMeshVirtualNode_VirtualNodeSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
