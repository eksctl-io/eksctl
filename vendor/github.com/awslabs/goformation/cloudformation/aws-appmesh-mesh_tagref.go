package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshMesh_TagRef AWS CloudFormation Resource (AWS::AppMesh::Mesh.TagRef)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-mesh-tagref.html
type AWSAppMeshMesh_TagRef struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-mesh-tagref.html#cfn-appmesh-mesh-tagref-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-mesh-tagref.html#cfn-appmesh-mesh-tagref-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshMesh_TagRef) AWSCloudFormationType() string {
	return "AWS::AppMesh::Mesh.TagRef"
}

func (r *AWSAppMeshMesh_TagRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
