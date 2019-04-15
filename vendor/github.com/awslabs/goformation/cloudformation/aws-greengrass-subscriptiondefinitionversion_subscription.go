package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassSubscriptionDefinitionVersion_Subscription AWS CloudFormation Resource (AWS::Greengrass::SubscriptionDefinitionVersion.Subscription)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinitionversion-subscription.html
type AWSGreengrassSubscriptionDefinitionVersion_Subscription struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinitionversion-subscription.html#cfn-greengrass-subscriptiondefinitionversion-subscription-id
	Id *Value `json:"Id,omitempty"`

	// Source AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinitionversion-subscription.html#cfn-greengrass-subscriptiondefinitionversion-subscription-source
	Source *Value `json:"Source,omitempty"`

	// Subject AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinitionversion-subscription.html#cfn-greengrass-subscriptiondefinitionversion-subscription-subject
	Subject *Value `json:"Subject,omitempty"`

	// Target AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinitionversion-subscription.html#cfn-greengrass-subscriptiondefinitionversion-subscription-target
	Target *Value `json:"Target,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassSubscriptionDefinitionVersion_Subscription) AWSCloudFormationType() string {
	return "AWS::Greengrass::SubscriptionDefinitionVersion.Subscription"
}

func (r *AWSGreengrassSubscriptionDefinitionVersion_Subscription) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
