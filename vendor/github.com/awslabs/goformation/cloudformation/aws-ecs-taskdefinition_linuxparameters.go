package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_LinuxParameters AWS CloudFormation Resource (AWS::ECS::TaskDefinition.LinuxParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-linuxparameters.html
type AWSECSTaskDefinition_LinuxParameters struct {

	// Capabilities AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-linuxparameters.html#cfn-ecs-taskdefinition-linuxparameters-capabilities
	Capabilities *AWSECSTaskDefinition_KernelCapabilities `json:"Capabilities,omitempty"`

	// Devices AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-linuxparameters.html#cfn-ecs-taskdefinition-linuxparameters-devices
	Devices []AWSECSTaskDefinition_Device `json:"Devices,omitempty"`

	// InitProcessEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-linuxparameters.html#cfn-ecs-taskdefinition-linuxparameters-initprocessenabled
	InitProcessEnabled *Value `json:"InitProcessEnabled,omitempty"`

	// SharedMemorySize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-linuxparameters.html#cfn-ecs-taskdefinition-linuxparameters-sharedmemorysize
	SharedMemorySize *Value `json:"SharedMemorySize,omitempty"`

	// Tmpfs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-linuxparameters.html#cfn-ecs-taskdefinition-linuxparameters-tmpfs
	Tmpfs []AWSECSTaskDefinition_Tmpfs `json:"Tmpfs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_LinuxParameters) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.LinuxParameters"
}

func (r *AWSECSTaskDefinition_LinuxParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
