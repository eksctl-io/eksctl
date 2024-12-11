package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_SalesforceChatterFeedConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.SalesforceChatterFeedConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcechatterfeedconfiguration.html
type DataSource_SalesforceChatterFeedConfiguration struct {

	// DocumentDataFieldName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcechatterfeedconfiguration.html#cfn-kendra-datasource-salesforcechatterfeedconfiguration-documentdatafieldname
	DocumentDataFieldName *types.Value `json:"DocumentDataFieldName,omitempty"`

	// DocumentTitleFieldName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcechatterfeedconfiguration.html#cfn-kendra-datasource-salesforcechatterfeedconfiguration-documenttitlefieldname
	DocumentTitleFieldName *types.Value `json:"DocumentTitleFieldName,omitempty"`

	// FieldMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcechatterfeedconfiguration.html#cfn-kendra-datasource-salesforcechatterfeedconfiguration-fieldmappings
	FieldMappings []DataSource_DataSourceToIndexFieldMapping `json:"FieldMappings,omitempty"`

	// IncludeFilterTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcechatterfeedconfiguration.html#cfn-kendra-datasource-salesforcechatterfeedconfiguration-includefiltertypes
	IncludeFilterTypes *types.Value `json:"IncludeFilterTypes,omitempty"`

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
func (r *DataSource_SalesforceChatterFeedConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.SalesforceChatterFeedConfiguration"
}
