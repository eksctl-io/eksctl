package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// OriginEndpoint_DashPackage AWS CloudFormation Resource (AWS::MediaPackage::OriginEndpoint.DashPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html
type OriginEndpoint_DashPackage struct {

	// AdTriggers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-adtriggers
	AdTriggers *types.Value `json:"AdTriggers,omitempty"`

	// AdsOnDeliveryRestrictions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-adsondeliveryrestrictions
	AdsOnDeliveryRestrictions *types.Value `json:"AdsOnDeliveryRestrictions,omitempty"`

	// Encryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-encryption
	Encryption *OriginEndpoint_DashEncryption `json:"Encryption,omitempty"`

	// ManifestLayout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-manifestlayout
	ManifestLayout *types.Value `json:"ManifestLayout,omitempty"`

	// ManifestWindowSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-manifestwindowseconds
	ManifestWindowSeconds *types.Value `json:"ManifestWindowSeconds,omitempty"`

	// MinBufferTimeSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-minbuffertimeseconds
	MinBufferTimeSeconds *types.Value `json:"MinBufferTimeSeconds,omitempty"`

	// MinUpdatePeriodSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-minupdateperiodseconds
	MinUpdatePeriodSeconds *types.Value `json:"MinUpdatePeriodSeconds,omitempty"`

	// PeriodTriggers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-periodtriggers
	PeriodTriggers *types.Value `json:"PeriodTriggers,omitempty"`

	// Profile AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-profile
	Profile *types.Value `json:"Profile,omitempty"`

	// SegmentDurationSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-segmentdurationseconds
	SegmentDurationSeconds *types.Value `json:"SegmentDurationSeconds,omitempty"`

	// SegmentTemplateFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-segmenttemplateformat
	SegmentTemplateFormat *types.Value `json:"SegmentTemplateFormat,omitempty"`

	// StreamSelection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-streamselection
	StreamSelection *OriginEndpoint_StreamSelection `json:"StreamSelection,omitempty"`

	// SuggestedPresentationDelaySeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-suggestedpresentationdelayseconds
	SuggestedPresentationDelaySeconds *types.Value `json:"SuggestedPresentationDelaySeconds,omitempty"`

	// UtcTiming AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-utctiming
	UtcTiming *types.Value `json:"UtcTiming,omitempty"`

	// UtcTimingUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-dashpackage.html#cfn-mediapackage-originendpoint-dashpackage-utctiminguri
	UtcTimingUri *types.Value `json:"UtcTimingUri,omitempty"`

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
func (r *OriginEndpoint_DashPackage) AWSCloudFormationType() string {
	return "AWS::MediaPackage::OriginEndpoint.DashPackage"
}
