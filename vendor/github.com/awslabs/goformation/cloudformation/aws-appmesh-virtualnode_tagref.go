package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_TagRef AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.TagRef)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tagref.html
type AWSAppMeshVirtualNode_TagRef struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tagref.html#cfn-appmesh-virtualnode-tagref-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-tagref.html#cfn-appmesh-virtualnode-tagref-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_TagRef) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.TagRef"
}

func (r *AWSAppMeshVirtualNode_TagRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
