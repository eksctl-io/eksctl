package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_Device AWS CloudFormation Resource (AWS::ECS::TaskDefinition.Device)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html
type AWSECSTaskDefinition_Device struct {

	// ContainerPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html#cfn-ecs-taskdefinition-device-containerpath
	ContainerPath *Value `json:"ContainerPath,omitempty"`

	// HostPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html#cfn-ecs-taskdefinition-device-hostpath
	HostPath *Value `json:"HostPath,omitempty"`

	// Permissions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html#cfn-ecs-taskdefinition-device-permissions
	Permissions []*Value `json:"Permissions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_Device) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.Device"
}

func (r *AWSECSTaskDefinition_Device) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
