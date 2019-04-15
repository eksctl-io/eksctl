package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualRouter_VirtualRouterListener AWS CloudFormation Resource (AWS::AppMesh::VirtualRouter.VirtualRouterListener)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-virtualrouterlistener.html
type AWSAppMeshVirtualRouter_VirtualRouterListener struct {

	// PortMapping AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-virtualrouterlistener.html#cfn-appmesh-virtualrouter-virtualrouterlistener-portmapping
	PortMapping *AWSAppMeshVirtualRouter_PortMapping `json:"PortMapping,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualRouter_VirtualRouterListener) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualRouter.VirtualRouterListener"
}

func (r *AWSAppMeshVirtualRouter_VirtualRouterListener) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
