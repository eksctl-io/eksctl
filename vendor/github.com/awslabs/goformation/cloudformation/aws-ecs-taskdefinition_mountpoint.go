package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_MountPoint AWS CloudFormation Resource (AWS::ECS::TaskDefinition.MountPoint)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-mountpoints.html
type AWSECSTaskDefinition_MountPoint struct {

	// ContainerPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-mountpoints.html#cfn-ecs-taskdefinition-containerdefinition-mountpoints-containerpath
	ContainerPath *Value `json:"ContainerPath,omitempty"`

	// ReadOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-mountpoints.html#cfn-ecs-taskdefinition-containerdefinition-mountpoints-readonly
	ReadOnly *Value `json:"ReadOnly,omitempty"`

	// SourceVolume AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-mountpoints.html#cfn-ecs-taskdefinition-containerdefinition-mountpoints-sourcevolume
	SourceVolume *Value `json:"SourceVolume,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_MountPoint) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.MountPoint"
}

func (r *AWSECSTaskDefinition_MountPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
