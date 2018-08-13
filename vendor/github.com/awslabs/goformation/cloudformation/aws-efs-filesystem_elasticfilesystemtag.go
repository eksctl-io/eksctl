package cloudformation

import (
	"encoding/json"
)

// AWSEFSFileSystem_ElasticFileSystemTag AWS CloudFormation Resource (AWS::EFS::FileSystem.ElasticFileSystemTag)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-efs-filesystem-filesystemtags.html
type AWSEFSFileSystem_ElasticFileSystemTag struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-efs-filesystem-filesystemtags.html#cfn-efs-filesystem-filesystemtags-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-efs-filesystem-filesystemtags.html#cfn-efs-filesystem-filesystemtags-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEFSFileSystem_ElasticFileSystemTag) AWSCloudFormationType() string {
	return "AWS::EFS::FileSystem.ElasticFileSystemTag"
}

func (r *AWSEFSFileSystem_ElasticFileSystemTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
