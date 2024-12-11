package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PackagingConfiguration_CmafPackage AWS CloudFormation Resource (AWS::MediaPackage::PackagingConfiguration.CmafPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-cmafpackage.html
type PackagingConfiguration_CmafPackage struct {

	// Encryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-cmafpackage.html#cfn-mediapackage-packagingconfiguration-cmafpackage-encryption
	Encryption *PackagingConfiguration_CmafEncryption `json:"Encryption,omitempty"`

	// HlsManifests AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-cmafpackage.html#cfn-mediapackage-packagingconfiguration-cmafpackage-hlsmanifests
	HlsManifests []PackagingConfiguration_HlsManifest `json:"HlsManifests,omitempty"`

	// IncludeEncoderConfigurationInSegments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-cmafpackage.html#cfn-mediapackage-packagingconfiguration-cmafpackage-includeencoderconfigurationinsegments
	IncludeEncoderConfigurationInSegments *types.Value `json:"IncludeEncoderConfigurationInSegments,omitempty"`

	// SegmentDurationSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-cmafpackage.html#cfn-mediapackage-packagingconfiguration-cmafpackage-segmentdurationseconds
	SegmentDurationSeconds *types.Value `json:"SegmentDurationSeconds,omitempty"`

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
func (r *PackagingConfiguration_CmafPackage) AWSCloudFormationType() string {
	return "AWS::MediaPackage::PackagingConfiguration.CmafPackage"
}
