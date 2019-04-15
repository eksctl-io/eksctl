package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualService_VirtualServiceProvider AWS CloudFormation Resource (AWS::AppMesh::VirtualService.VirtualServiceProvider)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualserviceprovider.html
type AWSAppMeshVirtualService_VirtualServiceProvider struct {

	// VirtualNode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualserviceprovider.html#cfn-appmesh-virtualservice-virtualserviceprovider-virtualnode
	VirtualNode *AWSAppMeshVirtualService_VirtualNodeServiceProvider `json:"VirtualNode,omitempty"`

	// VirtualRouter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualserviceprovider.html#cfn-appmesh-virtualservice-virtualserviceprovider-virtualrouter
	VirtualRouter *AWSAppMeshVirtualService_VirtualRouterServiceProvider `json:"VirtualRouter,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualService_VirtualServiceProvider) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualService.VirtualServiceProvider"
}

func (r *AWSAppMeshVirtualService_VirtualServiceProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
