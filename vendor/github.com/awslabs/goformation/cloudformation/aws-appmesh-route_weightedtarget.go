package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_WeightedTarget AWS CloudFormation Resource (AWS::AppMesh::Route.WeightedTarget)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-weightedtarget.html
type AWSAppMeshRoute_WeightedTarget struct {

	// VirtualNode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-weightedtarget.html#cfn-appmesh-route-weightedtarget-virtualnode
	VirtualNode *Value `json:"VirtualNode,omitempty"`

	// Weight AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-weightedtarget.html#cfn-appmesh-route-weightedtarget-weight
	Weight *Value `json:"Weight,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_WeightedTarget) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.WeightedTarget"
}

func (r *AWSAppMeshRoute_WeightedTarget) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
