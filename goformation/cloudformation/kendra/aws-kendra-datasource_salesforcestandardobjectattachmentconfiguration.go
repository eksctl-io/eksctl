package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_SalesforceStandardObjectAttachmentConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.SalesforceStandardObjectAttachmentConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcestandardobjectattachmentconfiguration.html
type DataSource_SalesforceStandardObjectAttachmentConfiguration struct {

	// DocumentTitleFieldName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcestandardobjectattachmentconfiguration.html#cfn-kendra-datasource-salesforcestandardobjectattachmentconfiguration-documenttitlefieldname
	DocumentTitleFieldName *types.Value `json:"DocumentTitleFieldName,omitempty"`

	// FieldMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcestandardobjectattachmentconfiguration.html#cfn-kendra-datasource-salesforcestandardobjectattachmentconfiguration-fieldmappings
	FieldMappings []DataSource_DataSourceToIndexFieldMapping `json:"FieldMappings,omitempty"`

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
func (r *DataSource_SalesforceStandardObjectAttachmentConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.SalesforceStandardObjectAttachmentConfiguration"
}
