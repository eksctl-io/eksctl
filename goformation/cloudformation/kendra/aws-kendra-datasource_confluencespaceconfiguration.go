package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_ConfluenceSpaceConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.ConfluenceSpaceConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespaceconfiguration.html
type DataSource_ConfluenceSpaceConfiguration struct {

	// CrawlArchivedSpaces AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespaceconfiguration.html#cfn-kendra-datasource-confluencespaceconfiguration-crawlarchivedspaces
	CrawlArchivedSpaces *types.Value `json:"CrawlArchivedSpaces,omitempty"`

	// CrawlPersonalSpaces AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespaceconfiguration.html#cfn-kendra-datasource-confluencespaceconfiguration-crawlpersonalspaces
	CrawlPersonalSpaces *types.Value `json:"CrawlPersonalSpaces,omitempty"`

	// ExcludeSpaces AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespaceconfiguration.html#cfn-kendra-datasource-confluencespaceconfiguration-excludespaces
	ExcludeSpaces *types.Value `json:"ExcludeSpaces,omitempty"`

	// IncludeSpaces AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespaceconfiguration.html#cfn-kendra-datasource-confluencespaceconfiguration-includespaces
	IncludeSpaces *types.Value `json:"IncludeSpaces,omitempty"`

	// SpaceFieldMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespaceconfiguration.html#cfn-kendra-datasource-confluencespaceconfiguration-spacefieldmappings
	SpaceFieldMappings []DataSource_ConfluenceSpaceToIndexFieldMapping `json:"SpaceFieldMappings,omitempty"`

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
func (r *DataSource_ConfluenceSpaceConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ConfluenceSpaceConfiguration"
}
