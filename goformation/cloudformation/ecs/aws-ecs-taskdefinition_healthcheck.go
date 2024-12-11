package ecs

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// TaskDefinition_HealthCheck AWS CloudFormation Resource (AWS::ECS::TaskDefinition.HealthCheck)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html
type TaskDefinition_HealthCheck struct {

	// Command AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-command
	Command *types.Value `json:"Command,omitempty"`

	// Interval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-interval
	Interval *types.Value `json:"Interval,omitempty"`

	// Retries AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-retries
	Retries *types.Value `json:"Retries,omitempty"`

	// StartPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-startperiod
	StartPeriod *types.Value `json:"StartPeriod,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-timeout
	Timeout *types.Value `json:"Timeout,omitempty"`

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
func (r *TaskDefinition_HealthCheck) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.HealthCheck"
}
