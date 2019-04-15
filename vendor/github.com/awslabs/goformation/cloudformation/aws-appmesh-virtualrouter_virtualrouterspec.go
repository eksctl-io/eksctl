package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualRouter_VirtualRouterSpec AWS CloudFormation Resource (AWS::AppMesh::VirtualRouter.VirtualRouterSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-virtualrouterspec.html
type AWSAppMeshVirtualRouter_VirtualRouterSpec struct {

	// Listeners AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-virtualrouterspec.html#cfn-appmesh-virtualrouter-virtualrouterspec-listeners
	Listeners []AWSAppMeshVirtualRouter_VirtualRouterListener `json:"Listeners,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualRouter_VirtualRouterSpec) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualRouter.VirtualRouterSpec"
}

func (r *AWSAppMeshVirtualRouter_VirtualRouterSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
