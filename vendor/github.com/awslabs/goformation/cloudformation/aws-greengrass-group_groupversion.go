package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassGroup_GroupVersion AWS CloudFormation Resource (AWS::Greengrass::Group.GroupVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html
type AWSGreengrassGroup_GroupVersion struct {

	// ConnectorDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html#cfn-greengrass-group-groupversion-connectordefinitionversionarn
	ConnectorDefinitionVersionArn *Value `json:"ConnectorDefinitionVersionArn,omitempty"`

	// CoreDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html#cfn-greengrass-group-groupversion-coredefinitionversionarn
	CoreDefinitionVersionArn *Value `json:"CoreDefinitionVersionArn,omitempty"`

	// DeviceDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html#cfn-greengrass-group-groupversion-devicedefinitionversionarn
	DeviceDefinitionVersionArn *Value `json:"DeviceDefinitionVersionArn,omitempty"`

	// FunctionDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html#cfn-greengrass-group-groupversion-functiondefinitionversionarn
	FunctionDefinitionVersionArn *Value `json:"FunctionDefinitionVersionArn,omitempty"`

	// LoggerDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html#cfn-greengrass-group-groupversion-loggerdefinitionversionarn
	LoggerDefinitionVersionArn *Value `json:"LoggerDefinitionVersionArn,omitempty"`

	// ResourceDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html#cfn-greengrass-group-groupversion-resourcedefinitionversionarn
	ResourceDefinitionVersionArn *Value `json:"ResourceDefinitionVersionArn,omitempty"`

	// SubscriptionDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-group-groupversion.html#cfn-greengrass-group-groupversion-subscriptiondefinitionversionarn
	SubscriptionDefinitionVersionArn *Value `json:"SubscriptionDefinitionVersionArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassGroup_GroupVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::Group.GroupVersion"
}

func (r *AWSGreengrassGroup_GroupVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
