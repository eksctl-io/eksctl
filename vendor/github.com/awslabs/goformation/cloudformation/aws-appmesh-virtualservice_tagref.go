package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualService_TagRef AWS CloudFormation Resource (AWS::AppMesh::VirtualService.TagRef)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-tagref.html
type AWSAppMeshVirtualService_TagRef struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-tagref.html#cfn-appmesh-virtualservice-tagref-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-tagref.html#cfn-appmesh-virtualservice-tagref-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualService_TagRef) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualService.TagRef"
}

func (r *AWSAppMeshVirtualService_TagRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
