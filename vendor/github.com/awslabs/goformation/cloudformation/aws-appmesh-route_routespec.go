package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_RouteSpec AWS CloudFormation Resource (AWS::AppMesh::Route.RouteSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-routespec.html
type AWSAppMeshRoute_RouteSpec struct {

	// HttpRoute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-routespec.html#cfn-appmesh-route-routespec-httproute
	HttpRoute *AWSAppMeshRoute_HttpRoute `json:"HttpRoute,omitempty"`

	// TcpRoute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-routespec.html#cfn-appmesh-route-routespec-tcproute
	TcpRoute *AWSAppMeshRoute_TcpRoute `json:"TcpRoute,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_RouteSpec) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.RouteSpec"
}

func (r *AWSAppMeshRoute_RouteSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
