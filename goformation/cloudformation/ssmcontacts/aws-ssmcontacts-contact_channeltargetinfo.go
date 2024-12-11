package ssmcontacts

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Contact_ChannelTargetInfo AWS CloudFormation Resource (AWS::SSMContacts::Contact.ChannelTargetInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-channeltargetinfo.html
type Contact_ChannelTargetInfo struct {

	// ChannelId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-channeltargetinfo.html#cfn-ssmcontacts-contact-channeltargetinfo-channelid
	ChannelId *types.Value `json:"ChannelId,omitempty"`

	// RetryIntervalInMinutes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-channeltargetinfo.html#cfn-ssmcontacts-contact-channeltargetinfo-retryintervalinminutes
	RetryIntervalInMinutes *types.Value `json:"RetryIntervalInMinutes"`

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
func (r *Contact_ChannelTargetInfo) AWSCloudFormationType() string {
	return "AWS::SSMContacts::Contact.ChannelTargetInfo"
}
