package cloudformation

import (
	"encoding/json"
)

// AWSECSTaskDefinition_VolumeFrom AWS CloudFormation Resource (AWS::ECS::TaskDefinition.VolumeFrom)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-volumesfrom.html
type AWSECSTaskDefinition_VolumeFrom struct {

	// ReadOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-volumesfrom.html#cfn-ecs-taskdefinition-containerdefinition-volumesfrom-readonly
	ReadOnly *Value `json:"ReadOnly,omitempty"`

	// SourceContainer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-containerdefinitions-volumesfrom.html#cfn-ecs-taskdefinition-containerdefinition-volumesfrom-sourcecontainer
	SourceContainer *Value `json:"SourceContainer,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_VolumeFrom) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.VolumeFrom"
}

func (r *AWSECSTaskDefinition_VolumeFrom) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
