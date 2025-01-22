package gamelift

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// GameServerGroup_InstanceDefinitions AWS CloudFormation Resource (AWS::GameLift::GameServerGroup.InstanceDefinitions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-instancedefinitions.html
type GameServerGroup_InstanceDefinitions struct {

	// InstanceDefinitions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-gameservergroup-instancedefinitions.html#cfn-gamelift-gameservergroup-instancedefinitions-instancedefinitions
	InstanceDefinitions []GameServerGroup_InstanceDefinition `json:"InstanceDefinitions,omitempty"`

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
func (r *GameServerGroup_InstanceDefinitions) AWSCloudFormationType() string {
	return "AWS::GameLift::GameServerGroup.InstanceDefinitions"
}
