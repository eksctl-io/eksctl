package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassSubscriptionDefinition_SubscriptionDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::SubscriptionDefinition.SubscriptionDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinition-subscriptiondefinitionversion.html
type AWSGreengrassSubscriptionDefinition_SubscriptionDefinitionVersion struct {

	// Subscriptions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-subscriptiondefinition-subscriptiondefinitionversion.html#cfn-greengrass-subscriptiondefinition-subscriptiondefinitionversion-subscriptions
	Subscriptions []AWSGreengrassSubscriptionDefinition_Subscription `json:"Subscriptions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassSubscriptionDefinition_SubscriptionDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::SubscriptionDefinition.SubscriptionDefinitionVersion"
}

func (r *AWSGreengrassSubscriptionDefinition_SubscriptionDefinitionVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
