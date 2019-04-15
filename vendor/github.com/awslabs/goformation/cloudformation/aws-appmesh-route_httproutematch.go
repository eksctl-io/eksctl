package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_HttpRouteMatch AWS CloudFormation Resource (AWS::AppMesh::Route.HttpRouteMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-httproutematch.html
type AWSAppMeshRoute_HttpRouteMatch struct {

	// Prefix AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-httproutematch.html#cfn-appmesh-route-httproutematch-prefix
	Prefix *Value `json:"Prefix,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_HttpRouteMatch) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.HttpRouteMatch"
}

func (r *AWSAppMeshRoute_HttpRouteMatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
