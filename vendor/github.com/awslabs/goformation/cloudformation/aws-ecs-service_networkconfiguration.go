package cloudformation

// AWSECSService_NetworkConfiguration AWS CloudFormation Resource (AWS::ECS::Service.NetworkConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-networkconfiguration.html
type AWSECSService_NetworkConfiguration struct {

	// AwsvpcConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-networkconfiguration.html#cfn-ecs-service-networkconfiguration-awsvpcconfiguration
	AwsvpcConfiguration *AWSECSService_AwsVpcConfiguration `json:"AwsvpcConfiguration,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSService_NetworkConfiguration) AWSCloudFormationType() string {
	return "AWS::ECS::Service.NetworkConfiguration"
}
