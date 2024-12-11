package wafv2

import (
	"goformation/v4/cloudformation/policies"
)

// LoggingConfiguration_FieldToMatch AWS CloudFormation Resource (AWS::WAFv2::LoggingConfiguration.FieldToMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-loggingconfiguration-fieldtomatch.html
type LoggingConfiguration_FieldToMatch struct {

	// JsonBody AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-loggingconfiguration-fieldtomatch.html#cfn-wafv2-loggingconfiguration-fieldtomatch-jsonbody
	JsonBody interface{} `json:"JsonBody,omitempty"`

	// Method AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-loggingconfiguration-fieldtomatch.html#cfn-wafv2-loggingconfiguration-fieldtomatch-method
	Method interface{} `json:"Method,omitempty"`

	// QueryString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-loggingconfiguration-fieldtomatch.html#cfn-wafv2-loggingconfiguration-fieldtomatch-querystring
	QueryString interface{} `json:"QueryString,omitempty"`

	// SingleHeader AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-loggingconfiguration-fieldtomatch.html#cfn-wafv2-loggingconfiguration-fieldtomatch-singleheader
	SingleHeader interface{} `json:"SingleHeader,omitempty"`

	// UriPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-loggingconfiguration-fieldtomatch.html#cfn-wafv2-loggingconfiguration-fieldtomatch-uripath
	UriPath interface{} `json:"UriPath,omitempty"`

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
func (r *LoggingConfiguration_FieldToMatch) AWSCloudFormationType() string {
	return "AWS::WAFv2::LoggingConfiguration.FieldToMatch"
}
