package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_Tmpfs AWS CloudFormation Resource (AWS::ECS::TaskDefinition.Tmpfs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-tmpfs.html
type AWSECSTaskDefinition_Tmpfs struct {

	// ContainerPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-tmpfs.html#cfn-ecs-taskdefinition-tmpfs-containerpath
	ContainerPath *Value `json:"ContainerPath,omitempty"`

	// MountOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-tmpfs.html#cfn-ecs-taskdefinition-tmpfs-mountoptions
	MountOptions []*Value `json:"MountOptions,omitempty"`

	// Size AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-tmpfs.html#cfn-ecs-taskdefinition-tmpfs-size
	Size *Value `json:"Size,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_Tmpfs) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.Tmpfs"
}

func (r *AWSECSTaskDefinition_Tmpfs) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
