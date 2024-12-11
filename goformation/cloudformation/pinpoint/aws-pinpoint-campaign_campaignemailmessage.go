package pinpoint

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Campaign_CampaignEmailMessage AWS CloudFormation Resource (AWS::Pinpoint::Campaign.CampaignEmailMessage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaignemailmessage.html
type Campaign_CampaignEmailMessage struct {

	// Body AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaignemailmessage.html#cfn-pinpoint-campaign-campaignemailmessage-body
	Body *types.Value `json:"Body,omitempty"`

	// FromAddress AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaignemailmessage.html#cfn-pinpoint-campaign-campaignemailmessage-fromaddress
	FromAddress *types.Value `json:"FromAddress,omitempty"`

	// HtmlBody AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaignemailmessage.html#cfn-pinpoint-campaign-campaignemailmessage-htmlbody
	HtmlBody *types.Value `json:"HtmlBody,omitempty"`

	// Title AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-pinpoint-campaign-campaignemailmessage.html#cfn-pinpoint-campaign-campaignemailmessage-title
	Title *types.Value `json:"Title,omitempty"`

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
func (r *Campaign_CampaignEmailMessage) AWSCloudFormationType() string {
	return "AWS::Pinpoint::Campaign.CampaignEmailMessage"
}
