package mwaa

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Environment_ModuleLoggingConfiguration AWS CloudFormation Resource (AWS::MWAA::Environment.ModuleLoggingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-moduleloggingconfiguration.html
type Environment_ModuleLoggingConfiguration struct {

	// CloudWatchLogGroupArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-moduleloggingconfiguration.html#cfn-mwaa-environment-moduleloggingconfiguration-cloudwatchloggrouparn
	CloudWatchLogGroupArn *types.Value `json:"CloudWatchLogGroupArn,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-moduleloggingconfiguration.html#cfn-mwaa-environment-moduleloggingconfiguration-enabled
	Enabled *types.Value `json:"Enabled,omitempty"`

	// LogLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-moduleloggingconfiguration.html#cfn-mwaa-environment-moduleloggingconfiguration-loglevel
	LogLevel *types.Value `json:"LogLevel,omitempty"`

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
func (r *Environment_ModuleLoggingConfiguration) AWSCloudFormationType() string {
	return "AWS::MWAA::Environment.ModuleLoggingConfiguration"
}
