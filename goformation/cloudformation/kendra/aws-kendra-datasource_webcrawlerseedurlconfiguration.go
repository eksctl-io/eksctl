package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_WebCrawlerSeedUrlConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.WebCrawlerSeedUrlConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerseedurlconfiguration.html
type DataSource_WebCrawlerSeedUrlConfiguration struct {

	// SeedUrls AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerseedurlconfiguration.html#cfn-kendra-datasource-webcrawlerseedurlconfiguration-seedurls
	SeedUrls *types.Value `json:"SeedUrls,omitempty"`

	// WebCrawlerMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerseedurlconfiguration.html#cfn-kendra-datasource-webcrawlerseedurlconfiguration-webcrawlermode
	WebCrawlerMode *types.Value `json:"WebCrawlerMode,omitempty"`

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
func (r *DataSource_WebCrawlerSeedUrlConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.WebCrawlerSeedUrlConfiguration"
}
