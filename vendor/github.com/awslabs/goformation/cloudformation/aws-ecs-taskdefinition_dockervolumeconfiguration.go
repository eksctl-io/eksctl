package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_DockerVolumeConfiguration AWS CloudFormation Resource (AWS::ECS::TaskDefinition.DockerVolumeConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-dockervolumeconfiguration.html
type AWSECSTaskDefinition_DockerVolumeConfiguration struct {

	// Autoprovision AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-dockervolumeconfiguration.html#cfn-ecs-taskdefinition-dockervolumeconfiguration-autoprovision
	Autoprovision *Value `json:"Autoprovision,omitempty"`

	// Driver AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-dockervolumeconfiguration.html#cfn-ecs-taskdefinition-dockervolumeconfiguration-driver
	Driver *Value `json:"Driver,omitempty"`

	// DriverOpts AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-dockervolumeconfiguration.html#cfn-ecs-taskdefinition-dockervolumeconfiguration-driveropts
	DriverOpts map[string]*Value `json:"DriverOpts,omitempty"`

	// Labels AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-dockervolumeconfiguration.html#cfn-ecs-taskdefinition-dockervolumeconfiguration-labels
	Labels map[string]*Value `json:"Labels,omitempty"`

	// Scope AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-dockervolumeconfiguration.html#cfn-ecs-taskdefinition-dockervolumeconfiguration-scope
	Scope *Value `json:"Scope,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_DockerVolumeConfiguration) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.DockerVolumeConfiguration"
}

func (r *AWSECSTaskDefinition_DockerVolumeConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
