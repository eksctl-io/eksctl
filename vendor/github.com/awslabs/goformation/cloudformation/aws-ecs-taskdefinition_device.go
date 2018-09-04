package cloudformation

// AWSECSTaskDefinition_Device AWS CloudFormation Resource (AWS::ECS::TaskDefinition.Device)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html
type AWSECSTaskDefinition_Device struct {

	// ContainerPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html#cfn-ecs-taskdefinition-device-containerpath
	ContainerPath *StringIntrinsic `json:"ContainerPath,omitempty"`

	// HostPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html#cfn-ecs-taskdefinition-device-hostpath
	HostPath *StringIntrinsic `json:"HostPath,omitempty"`

	// Permissions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-taskdefinition-device.html#cfn-ecs-taskdefinition-device-permissions
	Permissions []*StringIntrinsic `json:"Permissions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSTaskDefinition_Device) AWSCloudFormationType() string {
	return "AWS::ECS::TaskDefinition.Device"
}
