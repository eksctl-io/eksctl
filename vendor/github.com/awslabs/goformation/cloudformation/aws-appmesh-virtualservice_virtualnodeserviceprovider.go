package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualService_VirtualNodeServiceProvider AWS CloudFormation Resource (AWS::AppMesh::VirtualService.VirtualNodeServiceProvider)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualnodeserviceprovider.html
type AWSAppMeshVirtualService_VirtualNodeServiceProvider struct {

	// VirtualNodeName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualnodeserviceprovider.html#cfn-appmesh-virtualservice-virtualnodeserviceprovider-virtualnodename
	VirtualNodeName *Value `json:"VirtualNodeName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualService_VirtualNodeServiceProvider) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualService.VirtualNodeServiceProvider"
}

func (r *AWSAppMeshVirtualService_VirtualNodeServiceProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
