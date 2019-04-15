package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_Logging AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.Logging)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-logging.html
type AWSAppMeshVirtualNode_Logging struct {

	// AccessLog AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-logging.html#cfn-appmesh-virtualnode-logging-accesslog
	AccessLog *AWSAppMeshVirtualNode_AccessLog `json:"AccessLog,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_Logging) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.Logging"
}

func (r *AWSAppMeshVirtualNode_Logging) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
