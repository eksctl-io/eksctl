package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshMesh_MeshSpec AWS CloudFormation Resource (AWS::AppMesh::Mesh.MeshSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-mesh-meshspec.html
type AWSAppMeshMesh_MeshSpec struct {

	// EgressFilter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-mesh-meshspec.html#cfn-appmesh-mesh-meshspec-egressfilter
	EgressFilter *AWSAppMeshMesh_EgressFilter `json:"EgressFilter,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshMesh_MeshSpec) AWSCloudFormationType() string {
	return "AWS::AppMesh::Mesh.MeshSpec"
}

func (r *AWSAppMeshMesh_MeshSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
