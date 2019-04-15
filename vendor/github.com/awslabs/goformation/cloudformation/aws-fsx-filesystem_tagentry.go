package cloudformation

import (
	"encoding/json"
)

// AWSFSxFileSystem_TagEntry AWS CloudFormation Resource (AWS::FSx::FileSystem.TagEntry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-tagentry.html
type AWSFSxFileSystem_TagEntry struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-tagentry.html#cfn-fsx-filesystem-tagentry-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fsx-filesystem-tagentry.html#cfn-fsx-filesystem-tagentry-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSFSxFileSystem_TagEntry) AWSCloudFormationType() string {
	return "AWS::FSx::FileSystem.TagEntry"
}

func (r *AWSFSxFileSystem_TagEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
