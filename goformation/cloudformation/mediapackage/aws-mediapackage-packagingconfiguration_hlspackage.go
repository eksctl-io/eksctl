package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PackagingConfiguration_HlsPackage AWS CloudFormation Resource (AWS::MediaPackage::PackagingConfiguration.HlsPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlspackage.html
type PackagingConfiguration_HlsPackage struct {

	// Encryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlspackage.html#cfn-mediapackage-packagingconfiguration-hlspackage-encryption
	Encryption *PackagingConfiguration_HlsEncryption `json:"Encryption,omitempty"`

	// HlsManifests AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlspackage.html#cfn-mediapackage-packagingconfiguration-hlspackage-hlsmanifests
	HlsManifests []PackagingConfiguration_HlsManifest `json:"HlsManifests,omitempty"`

	// SegmentDurationSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlspackage.html#cfn-mediapackage-packagingconfiguration-hlspackage-segmentdurationseconds
	SegmentDurationSeconds *types.Value `json:"SegmentDurationSeconds,omitempty"`

	// UseAudioRenditionGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlspackage.html#cfn-mediapackage-packagingconfiguration-hlspackage-useaudiorenditiongroup
	UseAudioRenditionGroup *types.Value `json:"UseAudioRenditionGroup,omitempty"`

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
func (r *PackagingConfiguration_HlsPackage) AWSCloudFormationType() string {
	return "AWS::MediaPackage::PackagingConfiguration.HlsPackage"
}
