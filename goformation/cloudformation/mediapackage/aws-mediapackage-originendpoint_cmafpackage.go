package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// OriginEndpoint_CmafPackage AWS CloudFormation Resource (AWS::MediaPackage::OriginEndpoint.CmafPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafpackage.html
type OriginEndpoint_CmafPackage struct {

	// Encryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafpackage.html#cfn-mediapackage-originendpoint-cmafpackage-encryption
	Encryption *OriginEndpoint_CmafEncryption `json:"Encryption,omitempty"`

	// HlsManifests AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafpackage.html#cfn-mediapackage-originendpoint-cmafpackage-hlsmanifests
	HlsManifests []OriginEndpoint_HlsManifest `json:"HlsManifests,omitempty"`

	// SegmentDurationSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafpackage.html#cfn-mediapackage-originendpoint-cmafpackage-segmentdurationseconds
	SegmentDurationSeconds *types.Value `json:"SegmentDurationSeconds,omitempty"`

	// SegmentPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafpackage.html#cfn-mediapackage-originendpoint-cmafpackage-segmentprefix
	SegmentPrefix *types.Value `json:"SegmentPrefix,omitempty"`

	// StreamSelection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafpackage.html#cfn-mediapackage-originendpoint-cmafpackage-streamselection
	StreamSelection *OriginEndpoint_StreamSelection `json:"StreamSelection,omitempty"`

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
func (r *OriginEndpoint_CmafPackage) AWSCloudFormationType() string {
	return "AWS::MediaPackage::OriginEndpoint.CmafPackage"
}
