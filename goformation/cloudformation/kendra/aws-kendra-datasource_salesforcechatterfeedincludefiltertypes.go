package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_SalesforceChatterFeedIncludeFilterTypes AWS CloudFormation Resource (AWS::Kendra::DataSource.SalesforceChatterFeedIncludeFilterTypes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcechatterfeedincludefiltertypes.html
type DataSource_SalesforceChatterFeedIncludeFilterTypes struct {

	// SalesforceChatterFeedIncludeFilterTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcechatterfeedincludefiltertypes.html#cfn-kendra-datasource-salesforcechatterfeedincludefiltertypes-salesforcechatterfeedincludefiltertypes
	SalesforceChatterFeedIncludeFilterTypes *types.Value `json:"SalesforceChatterFeedIncludeFilterTypes,omitempty"`

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
func (r *DataSource_SalesforceChatterFeedIncludeFilterTypes) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.SalesforceChatterFeedIncludeFilterTypes"
}
