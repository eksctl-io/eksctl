package appstream

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DirectoryConfig_ServiceAccountCredentials AWS CloudFormation Resource (AWS::AppStream::DirectoryConfig.ServiceAccountCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-directoryconfig-serviceaccountcredentials.html
type DirectoryConfig_ServiceAccountCredentials struct {

	// AccountName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-directoryconfig-serviceaccountcredentials.html#cfn-appstream-directoryconfig-serviceaccountcredentials-accountname
	AccountName *types.Value `json:"AccountName,omitempty"`

	// AccountPassword AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appstream-directoryconfig-serviceaccountcredentials.html#cfn-appstream-directoryconfig-serviceaccountcredentials-accountpassword
	AccountPassword *types.Value `json:"AccountPassword,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *DirectoryConfig_ServiceAccountCredentials) AWSCloudFormationType() string {
	return "AWS::AppStream::DirectoryConfig.ServiceAccountCredentials"
}
