package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_HealthCheck AWS CloudFormation Resource (AWS::ECS::TaskDefinition.HealthCheck)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html
type AWSECSTaskDefinition_HealthCheck struct {

	// Command AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-command
	Command []*Value `json:"Command,omitempty"`

	// Interval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-interval
	Interval *Value `json:"Interval,omitempty"`

	// Retries AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-retries
	Retries *Value `json:"Retries,omitempty"`

	// StartPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-startperiod
	StartPeriod *Value `json:"StartPeriod,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-healthcheck.html#cfn-ecs-taskdefinition-healthcheck-timeout
	Timeout *Value `json:"Timeout,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_HealthCheck) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.HealthCheck"
}

func (r *AWSECSTaskDefinition_HealthCheck) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
