package apprunner

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Service_CodeConfigurationValues AWS CloudFormation Resource (AWS::AppRunner::Service.CodeConfigurationValues)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-codeconfigurationvalues.html
type Service_CodeConfigurationValues struct {

	// BuildCommand AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-codeconfigurationvalues.html#cfn-apprunner-service-codeconfigurationvalues-buildcommand
	BuildCommand *types.Value `json:"BuildCommand,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-codeconfigurationvalues.html#cfn-apprunner-service-codeconfigurationvalues-port
	Port *types.Value `json:"Port,omitempty"`

	// Runtime AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-codeconfigurationvalues.html#cfn-apprunner-service-codeconfigurationvalues-runtime
	Runtime *types.Value `json:"Runtime,omitempty"`

	// RuntimeEnvironmentVariables AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-codeconfigurationvalues.html#cfn-apprunner-service-codeconfigurationvalues-runtimeenvironmentvariables
	RuntimeEnvironmentVariables []Service_KeyValuePair `json:"RuntimeEnvironmentVariables,omitempty"`

	// StartCommand AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-codeconfigurationvalues.html#cfn-apprunner-service-codeconfigurationvalues-startcommand
	StartCommand *types.Value `json:"StartCommand,omitempty"`

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
func (r *Service_CodeConfigurationValues) AWSCloudFormationType() string {
	return "AWS::AppRunner::Service.CodeConfigurationValues"
}
