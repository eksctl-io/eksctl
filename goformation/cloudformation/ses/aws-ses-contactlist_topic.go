package ses

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ContactList_Topic AWS CloudFormation Resource (AWS::SES::ContactList.Topic)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-contactlist-topic.html
type ContactList_Topic struct {

	// DefaultSubscriptionStatus AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-contactlist-topic.html#cfn-ses-contactlist-topic-defaultsubscriptionstatus
	DefaultSubscriptionStatus *types.Value `json:"DefaultSubscriptionStatus,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-contactlist-topic.html#cfn-ses-contactlist-topic-description
	Description *types.Value `json:"Description,omitempty"`

	// DisplayName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-contactlist-topic.html#cfn-ses-contactlist-topic-displayname
	DisplayName *types.Value `json:"DisplayName,omitempty"`

	// TopicName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-contactlist-topic.html#cfn-ses-contactlist-topic-topicname
	TopicName *types.Value `json:"TopicName,omitempty"`

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
func (r *ContactList_Topic) AWSCloudFormationType() string {
	return "AWS::SES::ContactList.Topic"
}
