package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_SalesforceStandardObjectConfigurationList AWS CloudFormation Resource (AWS::Kendra::DataSource.SalesforceStandardObjectConfigurationList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcestandardobjectconfigurationlist.html
type DataSource_SalesforceStandardObjectConfigurationList struct {

	// SalesforceStandardObjectConfigurationList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcestandardobjectconfigurationlist.html#cfn-kendra-datasource-salesforcestandardobjectconfigurationlist-salesforcestandardobjectconfigurationlist
	SalesforceStandardObjectConfigurationList []DataSource_SalesforceStandardObjectConfiguration `json:"SalesforceStandardObjectConfigurationList,omitempty"`

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
func (r *DataSource_SalesforceStandardObjectConfigurationList) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.SalesforceStandardObjectConfigurationList"
}
