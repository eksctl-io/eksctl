package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_KeyValuePair AWS CloudFormation Resource (AWS::ECS::TaskDefinition.KeyValuePair)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-environment.html
type AWSECSTaskDefinition_KeyValuePair struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-environment.html#cfn-ecs-taskdefinition-containerdefinition-environment-name
	Name *Value `json:"Name,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-environment.html#cfn-ecs-taskdefinition-containerdefinition-environment-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_KeyValuePair) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.KeyValuePair"
}

func (r *AWSECSTaskDefinition_KeyValuePair) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
