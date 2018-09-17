package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_KernelCapabilities AWS CloudFormation Resource (AWS::ECS::TaskDefinition.KernelCapabilities)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-kernelcapabilities.html
type AWSECSTaskDefinition_KernelCapabilities struct {

	// Add AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-kernelcapabilities.html#cfn-ecs-taskdefinition-kernelcapabilities-add
	Add []*Value `json:"Add,omitempty"`

	// Drop AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-kernelcapabilities.html#cfn-ecs-taskdefinition-kernelcapabilities-drop
	Drop []*Value `json:"Drop,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_KernelCapabilities) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.KernelCapabilities"
}

func (r *AWSECSTaskDefinition_KernelCapabilities) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
