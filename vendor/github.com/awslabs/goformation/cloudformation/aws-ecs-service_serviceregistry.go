package cloudformation

// AWSECSService_ServiceRegistry AWS CloudFormation Resource (AWS::ECS::Service.ServiceRegistry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html
type AWSECSService_ServiceRegistry struct {

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html#cfn-ecs-service-serviceregistry-port
	Port int `json:"Port,omitempty"`

	// RegistryArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-serviceregistry.html#cfn-ecs-service-serviceregistry-registryarn
	RegistryArn *StringIntrinsic `json:"RegistryArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSService_ServiceRegistry) AWSCloudFormationType() string {
	return "AWS::ECS::Service.ServiceRegistry"
}
