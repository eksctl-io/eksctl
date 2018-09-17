package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_PortMapping AWS CloudFormation Resource (AWS::ECS::TaskDefinition.PortMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-portmappings.html
type AWSECSTaskDefinition_PortMapping struct {

	// ContainerPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-portmappings.html#cfn-ecs-taskdefinition-containerdefinition-portmappings-containerport
	ContainerPort *Value `json:"ContainerPort,omitempty"`

	// HostPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-portmappings.html#cfn-ecs-taskdefinition-containerdefinition-portmappings-readonly
	HostPort *Value `json:"HostPort,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-portmappings.html#cfn-ecs-taskdefinition-containerdefinition-portmappings-sourcevolume
	Protocol *Value `json:"Protocol,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_PortMapping) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.PortMapping"
}

func (r *AWSECSTaskDefinition_PortMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
