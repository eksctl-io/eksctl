package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PackagingConfiguration_DashManifest AWS CloudFormation Resource (AWS::MediaPackage::PackagingConfiguration.DashManifest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashmanifest.html
type PackagingConfiguration_DashManifest struct {

	// ManifestLayout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashmanifest.html#cfn-mediapackage-packagingconfiguration-dashmanifest-manifestlayout
	ManifestLayout *types.Value `json:"ManifestLayout,omitempty"`

	// ManifestName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashmanifest.html#cfn-mediapackage-packagingconfiguration-dashmanifest-manifestname
	ManifestName *types.Value `json:"ManifestName,omitempty"`

	// MinBufferTimeSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashmanifest.html#cfn-mediapackage-packagingconfiguration-dashmanifest-minbuffertimeseconds
	MinBufferTimeSeconds *types.Value `json:"MinBufferTimeSeconds,omitempty"`

	// Profile AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashmanifest.html#cfn-mediapackage-packagingconfiguration-dashmanifest-profile
	Profile *types.Value `json:"Profile,omitempty"`

	// StreamSelection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-dashmanifest.html#cfn-mediapackage-packagingconfiguration-dashmanifest-streamselection
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
func (r *PackagingConfiguration_DashManifest) AWSCloudFormationType() string {
	return "AWS::MediaPackage::PackagingConfiguration.DashManifest"
}
