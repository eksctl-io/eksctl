package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_VirtualServiceBackend AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.VirtualServiceBackend)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualservicebackend.html
type AWSAppMeshVirtualNode_VirtualServiceBackend struct {

	// VirtualServiceName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-virtualservicebackend.html#cfn-appmesh-virtualnode-virtualservicebackend-virtualservicename
	VirtualServiceName *Value `json:"VirtualServiceName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_VirtualServiceBackend) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.VirtualServiceBackend"
}

func (r *AWSAppMeshVirtualNode_VirtualServiceBackend) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
