package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_LogConfiguration AWS CloudFormation Resource (AWS::ECS::TaskDefinition.LogConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-logconfiguration.html
type AWSECSTaskDefinition_LogConfiguration struct {

	// LogDriver AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-logconfiguration.html#cfn-ecs-taskdefinition-containerdefinition-logconfiguration-logdriver
	LogDriver *Value `json:"LogDriver,omitempty"`

	// Options AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-logconfiguration.html#cfn-ecs-taskdefinition-containerdefinition-logconfiguration-options
	Options map[string]*Value `json:"Options,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_LogConfiguration) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.LogConfiguration"
}

func (r *AWSECSTaskDefinition_LogConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
