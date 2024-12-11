package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_ConfluenceAttachmentConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.ConfluenceAttachmentConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluenceattachmentconfiguration.html
type DataSource_ConfluenceAttachmentConfiguration struct {

	// AttachmentFieldMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluenceattachmentconfiguration.html#cfn-kendra-datasource-confluenceattachmentconfiguration-attachmentfieldmappings
	AttachmentFieldMappings []DataSource_ConfluenceAttachmentToIndexFieldMapping `json:"AttachmentFieldMappings,omitempty"`

	// CrawlAttachments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluenceattachmentconfiguration.html#cfn-kendra-datasource-confluenceattachmentconfiguration-crawlattachments
	CrawlAttachments *types.Value `json:"CrawlAttachments,omitempty"`

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
func (r *DataSource_ConfluenceAttachmentConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ConfluenceAttachmentConfiguration"
}
