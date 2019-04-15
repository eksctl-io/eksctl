package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_FileAccessLog AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.FileAccessLog)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-fileaccesslog.html
type AWSAppMeshVirtualNode_FileAccessLog struct {

	// Path AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-fileaccesslog.html#cfn-appmesh-virtualnode-fileaccesslog-path
	Path *Value `json:"Path,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_FileAccessLog) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.FileAccessLog"
}

func (r *AWSAppMeshVirtualNode_FileAccessLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
