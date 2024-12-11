package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_WebCrawlerConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.WebCrawlerConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html
type DataSource_WebCrawlerConfiguration struct {

	// AuthenticationConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-authenticationconfiguration
	AuthenticationConfiguration *DataSource_WebCrawlerAuthenticationConfiguration `json:"AuthenticationConfiguration,omitempty"`

	// CrawlDepth AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-crawldepth
	CrawlDepth *types.Value `json:"CrawlDepth,omitempty"`

	// MaxContentSizePerPageInMegaBytes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-maxcontentsizeperpageinmegabytes
	MaxContentSizePerPageInMegaBytes *types.Value `json:"MaxContentSizePerPageInMegaBytes,omitempty"`

	// MaxLinksPerPage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-maxlinksperpage
	MaxLinksPerPage *types.Value `json:"MaxLinksPerPage,omitempty"`

	// MaxUrlsPerMinuteCrawlRate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-maxurlsperminutecrawlrate
	MaxUrlsPerMinuteCrawlRate *types.Value `json:"MaxUrlsPerMinuteCrawlRate,omitempty"`

	// ProxyConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-proxyconfiguration
	ProxyConfiguration *DataSource_ProxyConfiguration `json:"ProxyConfiguration,omitempty"`

	// UrlExclusionPatterns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-urlexclusionpatterns
	UrlExclusionPatterns *types.Value `json:"UrlExclusionPatterns,omitempty"`

	// UrlInclusionPatterns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-urlinclusionpatterns
	UrlInclusionPatterns *types.Value `json:"UrlInclusionPatterns,omitempty"`

	// Urls AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-webcrawlerconfiguration.html#cfn-kendra-datasource-webcrawlerconfiguration-urls
	Urls *DataSource_WebCrawlerUrls `json:"Urls,omitempty"`

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
func (r *DataSource_WebCrawlerConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.WebCrawlerConfiguration"
}
