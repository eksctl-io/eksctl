package networkfirewall

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LoggingConfiguration_LogDestinationConfig AWS CloudFormation Resource (AWS::NetworkFirewall::LoggingConfiguration.LogDestinationConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-loggingconfiguration-logdestinationconfig.html
type LoggingConfiguration_LogDestinationConfig struct {

	// LogDestination AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-loggingconfiguration-logdestinationconfig.html#cfn-networkfirewall-loggingconfiguration-logdestinationconfig-logdestination
	LogDestination map[string]*types.Value `json:"LogDestination,omitempty"`

	// LogDestinationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-loggingconfiguration-logdestinationconfig.html#cfn-networkfirewall-loggingconfiguration-logdestinationconfig-logdestinationtype
	LogDestinationType *types.Value `json:"LogDestinationType,omitempty"`

	// LogType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-loggingconfiguration-logdestinationconfig.html#cfn-networkfirewall-loggingconfiguration-logdestinationconfig-logtype
	LogType *types.Value `json:"LogType,omitempty"`

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
func (r *LoggingConfiguration_LogDestinationConfig) AWSCloudFormationType() string {
	return "AWS::NetworkFirewall::LoggingConfiguration.LogDestinationConfig"
}
