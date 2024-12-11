package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PackagingConfiguration_DashPackage AWS CloudFormation Resource (AWS::MediaPackage::PackagingConfiguration.DashPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashpackage.html
type PackagingConfiguration_DashPackage struct {

	// DashManifests AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashpackage.html#cfn-mediapackage-packagingconfiguration-dashpackage-dashmanifests
	DashManifests []PackagingConfiguration_DashManifest `json:"DashManifests,omitempty"`

	// Encryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashpackage.html#cfn-mediapackage-packagingconfiguration-dashpackage-encryption
	Encryption *PackagingConfiguration_DashEncryption `json:"Encryption,omitempty"`

	// IncludeEncoderConfigurationInSegments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashpackage.html#cfn-mediapackage-packagingconfiguration-dashpackage-includeencoderconfigurationinsegments
	IncludeEncoderConfigurationInSegments *types.Value `json:"IncludeEncoderConfigurationInSegments,omitempty"`

	// PeriodTriggers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashpackage.html#cfn-mediapackage-packagingconfiguration-dashpackage-periodtriggers
	PeriodTriggers *types.Value `json:"PeriodTriggers,omitempty"`

	// SegmentDurationSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashpackage.html#cfn-mediapackage-packagingconfiguration-dashpackage-segmentdurationseconds
	SegmentDurationSeconds *types.Value `json:"SegmentDurationSeconds,omitempty"`

	// SegmentTemplateFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashpackage.html#cfn-mediapackage-packagingconfiguration-dashpackage-segmenttemplateformat
	SegmentTemplateFormat *types.Value `json:"SegmentTemplateFormat,omitempty"`

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
func (r *PackagingConfiguration_DashPackage) AWSCloudFormationType() string {
	return "AWS::MediaPackage::PackagingConfiguration.DashPackage"
}
