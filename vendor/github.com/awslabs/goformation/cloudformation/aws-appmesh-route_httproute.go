package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_HttpRoute AWS CloudFormation Resource (AWS::AppMesh::Route.HttpRoute)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-httproute.html
type AWSAppMeshRoute_HttpRoute struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-httproute.html#cfn-appmesh-route-httproute-action
	Action *AWSAppMeshRoute_HttpRouteAction `json:"Action,omitempty"`

	// Match AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-httproute.html#cfn-appmesh-route-httproute-match
	Match *AWSAppMeshRoute_HttpRouteMatch `json:"Match,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_HttpRoute) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.HttpRoute"
}

func (r *AWSAppMeshRoute_HttpRoute) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
