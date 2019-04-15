package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_TagRef AWS CloudFormation Resource (AWS::AppMesh::Route.TagRef)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-tagref.html
type AWSAppMeshRoute_TagRef struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-tagref.html#cfn-appmesh-route-tagref-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-tagref.html#cfn-appmesh-route-tagref-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_TagRef) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.TagRef"
}

func (r *AWSAppMeshRoute_TagRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
