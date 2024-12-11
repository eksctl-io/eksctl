package pinpoint

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Campaign_CampaignEventFilter AWS CloudFormation Resource (AWS::Pinpoint::Campaign.CampaignEventFilter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaigneventfilter.html
type Campaign_CampaignEventFilter struct {

	// Dimensions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaigneventfilter.html#cfn-pinpoint-campaign-campaigneventfilter-dimensions
	Dimensions *Campaign_EventDimensions `json:"Dimensions,omitempty"`

	// FilterType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaigneventfilter.html#cfn-pinpoint-campaign-campaigneventfilter-filtertype
	FilterType *types.Value `json:"FilterType,omitempty"`

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
func (r *Campaign_CampaignEventFilter) AWSCloudFormationType() string {
	return "AWS::Pinpoint::Campaign.CampaignEventFilter"
}
