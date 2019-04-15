package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualRouter_PortMapping AWS CloudFormation Resource (AWS::AppMesh::VirtualRouter.PortMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-portmapping.html
type AWSAppMeshVirtualRouter_PortMapping struct {

	// Port AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-portmapping.html#cfn-appmesh-virtualrouter-portmapping-port
	Port *Value `json:"Port,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-portmapping.html#cfn-appmesh-virtualrouter-portmapping-protocol
	Protocol *Value `json:"Protocol,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualRouter_PortMapping) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualRouter.PortMapping"
}

func (r *AWSAppMeshVirtualRouter_PortMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
