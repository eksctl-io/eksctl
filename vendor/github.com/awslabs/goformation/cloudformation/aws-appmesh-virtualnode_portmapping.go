package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_PortMapping AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.PortMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-portmapping.html
type AWSAppMeshVirtualNode_PortMapping struct {

	// Port AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-portmapping.html#cfn-appmesh-virtualnode-portmapping-port
	Port *Value `json:"Port,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-portmapping.html#cfn-appmesh-virtualnode-portmapping-protocol
	Protocol *Value `json:"Protocol,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_PortMapping) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.PortMapping"
}

func (r *AWSAppMeshVirtualNode_PortMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
