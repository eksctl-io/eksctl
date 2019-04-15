package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualService_VirtualServiceSpec AWS CloudFormation Resource (AWS::AppMesh::VirtualService.VirtualServiceSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualservicespec.html
type AWSAppMeshVirtualService_VirtualServiceSpec struct {

	// Provider AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualservicespec.html#cfn-appmesh-virtualservice-virtualservicespec-provider
	Provider *AWSAppMeshVirtualService_VirtualServiceProvider `json:"Provider,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualService_VirtualServiceSpec) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualService.VirtualServiceSpec"
}

func (r *AWSAppMeshVirtualService_VirtualServiceSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
