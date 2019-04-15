package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualRouter_TagRef AWS CloudFormation Resource (AWS::AppMesh::VirtualRouter.TagRef)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-tagref.html
type AWSAppMeshVirtualRouter_TagRef struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-tagref.html#cfn-appmesh-virtualrouter-tagref-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualrouter-tagref.html#cfn-appmesh-virtualrouter-tagref-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualRouter_TagRef) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualRouter.TagRef"
}

func (r *AWSAppMeshVirtualRouter_TagRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
