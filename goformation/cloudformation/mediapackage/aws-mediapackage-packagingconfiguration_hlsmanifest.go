package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PackagingConfiguration_HlsManifest AWS CloudFormation Resource (AWS::MediaPackage::PackagingConfiguration.HlsManifest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsmanifest.html
type PackagingConfiguration_HlsManifest struct {

	// AdMarkers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsmanifest.html#cfn-mediapackage-packagingconfiguration-hlsmanifest-admarkers
	AdMarkers *types.Value `json:"AdMarkers,omitempty"`

	// IncludeIframeOnlyStream AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsmanifest.html#cfn-mediapackage-packagingconfiguration-hlsmanifest-includeiframeonlystream
	IncludeIframeOnlyStream *types.Value `json:"IncludeIframeOnlyStream,omitempty"`

	// ManifestName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsmanifest.html#cfn-mediapackage-packagingconfiguration-hlsmanifest-manifestname
	ManifestName *types.Value `json:"ManifestName,omitempty"`

	// ProgramDateTimeIntervalSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsmanifest.html#cfn-mediapackage-packagingconfiguration-hlsmanifest-programdatetimeintervalseconds
	ProgramDateTimeIntervalSeconds *types.Value `json:"ProgramDateTimeIntervalSeconds,omitempty"`

	// RepeatExtXKey AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsmanifest.html#cfn-mediapackage-packagingconfiguration-hlsmanifest-repeatextxkey
	RepeatExtXKey *types.Value `json:"RepeatExtXKey,omitempty"`

	// StreamSelection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsmanifest.html#cfn-mediapackage-packagingconfiguration-hlsmanifest-streamselection
	StreamSelection *PackagingConfiguration_StreamSelection `json:"StreamSelection,omitempty"`

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
func (r *PackagingConfiguration_HlsManifest) AWSCloudFormationType() string {
	return "AWS::MediaPackage::PackagingConfiguration.HlsManifest"
}
