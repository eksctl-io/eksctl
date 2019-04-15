package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshMesh_EgressFilter AWS CloudFormation Resource (AWS::AppMesh::Mesh.EgressFilter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-mesh-egressfilter.html
type AWSAppMeshMesh_EgressFilter struct {

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-mesh-egressfilter.html#cfn-appmesh-mesh-egressfilter-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshMesh_EgressFilter) AWSCloudFormationType() string {
	return "AWS::AppMesh::Mesh.EgressFilter"
}

func (r *AWSAppMeshMesh_EgressFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
