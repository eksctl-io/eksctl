package cloudformation

import (
	"encoding/json"
)

// AWSAppStreamDirectoryConfig_ServiceAccountCredentials AWS CloudFormation Resource (AWS::AppStream::DirectoryConfig.ServiceAccountCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-directoryconfig-serviceaccountcredentials.html
type AWSAppStreamDirectoryConfig_ServiceAccountCredentials struct {

	// AccountName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-directoryconfig-serviceaccountcredentials.html#cfn-appstream-directoryconfig-serviceaccountcredentials-accountname
	AccountName *Value `json:"AccountName,omitempty"`

	// AccountPassword AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-directoryconfig-serviceaccountcredentials.html#cfn-appstream-directoryconfig-serviceaccountcredentials-accountpassword
	AccountPassword *Value `json:"AccountPassword,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamDirectoryConfig_ServiceAccountCredentials) AWSCloudFormationType() string {
	return "AWS::AppStream::DirectoryConfig.ServiceAccountCredentials"
}

func (r *AWSAppStreamDirectoryConfig_ServiceAccountCredentials) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
