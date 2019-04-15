package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_Volume AWS CloudFormation Resource (AWS::ECS::TaskDefinition.Volume)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-volumes.html
type AWSECSTaskDefinition_Volume struct {

	// DockerVolumeConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-volumes.html#cfn-ecs-taskdefinition-volume-dockervolumeconfiguration
	DockerVolumeConfiguration *AWSECSTaskDefinition_DockerVolumeConfiguration `json:"DockerVolumeConfiguration,omitempty"`

	// Host AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-volumes.html#cfn-ecs-taskdefinition-volumes-host
	Host *AWSECSTaskDefinition_HostVolumeProperties `json:"Host,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-volumes.html#cfn-ecs-taskdefinition-volumes-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_Volume) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.Volume"
}

func (r *AWSECSTaskDefinition_Volume) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
