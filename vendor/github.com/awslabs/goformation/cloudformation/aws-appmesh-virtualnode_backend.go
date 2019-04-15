package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_Backend AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.Backend)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-backend.html
type AWSAppMeshVirtualNode_Backend struct {

	// VirtualService AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-backend.html#cfn-appmesh-virtualnode-backend-virtualservice
	VirtualService *AWSAppMeshVirtualNode_VirtualServiceBackend `json:"VirtualService,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_Backend) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.Backend"
}

func (r *AWSAppMeshVirtualNode_Backend) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
