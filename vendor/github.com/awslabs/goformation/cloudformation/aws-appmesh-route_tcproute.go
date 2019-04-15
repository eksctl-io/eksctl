package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_TcpRoute AWS CloudFormation Resource (AWS::AppMesh::Route.TcpRoute)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-tcproute.html
type AWSAppMeshRoute_TcpRoute struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-tcproute.html#cfn-appmesh-route-tcproute-action
	Action *AWSAppMeshRoute_TcpRouteAction `json:"Action,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_TcpRoute) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.TcpRoute"
}

func (r *AWSAppMeshRoute_TcpRoute) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
