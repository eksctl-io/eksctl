package cloudformation

import (
	"encoding/json"
)

// AWSECSService_ServiceRegistry AWS CloudFormation Resource (AWS::ECS::Service.ServiceRegistry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html
type AWSECSService_ServiceRegistry struct {

	// ContainerName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html#cfn-ecs-service-serviceregistry-containername
	ContainerName *Value `json:"ContainerName,omitempty"`

	// ContainerPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html#cfn-ecs-service-serviceregistry-containerport
	ContainerPort *Value `json:"ContainerPort,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html#cfn-ecs-service-serviceregistry-port
	Port *Value `json:"Port,omitempty"`

	// RegistryArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html#cfn-ecs-service-serviceregistry-registryarn
	RegistryArn *Value `json:"RegistryArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSService_ServiceRegistry) AWSCloudFormationType() string {
	return "AWS::ECS::Service.ServiceRegistry"
}

func (r *AWSECSService_ServiceRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
