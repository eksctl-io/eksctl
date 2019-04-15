package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassSubscriptionDefinition_Subscription AWS CloudFormation Resource (AWS::Greengrass::SubscriptionDefinition.Subscription)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinition-subscription.html
type AWSGreengrassSubscriptionDefinition_Subscription struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinition-subscription.html#cfn-greengrass-subscriptiondefinition-subscription-id
	Id *Value `json:"Id,omitempty"`

	// Source AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinition-subscription.html#cfn-greengrass-subscriptiondefinition-subscription-source
	Source *Value `json:"Source,omitempty"`

	// Subject AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinition-subscription.html#cfn-greengrass-subscriptiondefinition-subscription-subject
	Subject *Value `json:"Subject,omitempty"`

	// Target AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinition-subscription.html#cfn-greengrass-subscriptiondefinition-subscription-target
	Target *Value `json:"Target,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassSubscriptionDefinition_Subscription) AWSCloudFormationType() string {
	return "AWS::Greengrass::SubscriptionDefinition.Subscription"
}

func (r *AWSGreengrassSubscriptionDefinition_Subscription) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
