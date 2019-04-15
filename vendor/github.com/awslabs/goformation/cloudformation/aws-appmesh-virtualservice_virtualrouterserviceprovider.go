package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualService_VirtualRouterServiceProvider AWS CloudFormation Resource (AWS::AppMesh::VirtualService.VirtualRouterServiceProvider)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualrouterserviceprovider.html
type AWSAppMeshVirtualService_VirtualRouterServiceProvider struct {

	// VirtualRouterName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualservice-virtualrouterserviceprovider.html#cfn-appmesh-virtualservice-virtualrouterserviceprovider-virtualroutername
	VirtualRouterName *Value `json:"VirtualRouterName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualService_VirtualRouterServiceProvider) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualService.VirtualRouterServiceProvider"
}

func (r *AWSAppMeshVirtualService_VirtualRouterServiceProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
