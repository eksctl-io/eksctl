package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_AccessLog AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.AccessLog)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-accesslog.html
type AWSAppMeshVirtualNode_AccessLog struct {

	// File AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-accesslog.html#cfn-appmesh-virtualnode-accesslog-file
	File *AWSAppMeshVirtualNode_FileAccessLog `json:"File,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_AccessLog) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.AccessLog"
}

func (r *AWSAppMeshVirtualNode_AccessLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
