package gamelift

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Fleet_RuntimeConfiguration AWS CloudFormation Resource (AWS::GameLift::Fleet.RuntimeConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-runtimeconfiguration.html
type Fleet_RuntimeConfiguration struct {

	// GameSessionActivationTimeoutSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-runtimeconfiguration.html#cfn-gamelift-fleet-runtimeconfiguration-gamesessionactivationtimeoutseconds
	GameSessionActivationTimeoutSeconds *types.Value `json:"GameSessionActivationTimeoutSeconds,omitempty"`

	// MaxConcurrentGameSessionActivations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-runtimeconfiguration.html#cfn-gamelift-fleet-runtimeconfiguration-maxconcurrentgamesessionactivations
	MaxConcurrentGameSessionActivations *types.Value `json:"MaxConcurrentGameSessionActivations,omitempty"`

	// ServerProcesses AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-runtimeconfiguration.html#cfn-gamelift-fleet-runtimeconfiguration-serverprocesses
	ServerProcesses []Fleet_ServerProcess `json:"ServerProcesses,omitempty"`

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
func (r *Fleet_RuntimeConfiguration) AWSCloudFormationType() string {
	return "AWS::GameLift::Fleet.RuntimeConfiguration"
}
