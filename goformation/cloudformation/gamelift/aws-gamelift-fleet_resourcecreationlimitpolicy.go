package gamelift

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Fleet_ResourceCreationLimitPolicy AWS CloudFormation Resource (AWS::GameLift::Fleet.ResourceCreationLimitPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-resourcecreationlimitpolicy.html
type Fleet_ResourceCreationLimitPolicy struct {

	// NewGameSessionsPerCreator AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-resourcecreationlimitpolicy.html#cfn-gamelift-fleet-resourcecreationlimitpolicy-newgamesessionspercreator
	NewGameSessionsPerCreator *types.Value `json:"NewGameSessionsPerCreator,omitempty"`

	// PolicyPeriodInMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-resourcecreationlimitpolicy.html#cfn-gamelift-fleet-resourcecreationlimitpolicy-policyperiodinminutes
	PolicyPeriodInMinutes *types.Value `json:"PolicyPeriodInMinutes,omitempty"`

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
func (r *Fleet_ResourceCreationLimitPolicy) AWSCloudFormationType() string {
	return "AWS::GameLift::Fleet.ResourceCreationLimitPolicy"
}
