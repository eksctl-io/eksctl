package cloudformation

// AWSECSService_AwsVpcConfiguration AWS CloudFormation Resource (AWS::ECS::Service.AwsVpcConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-awsvpcconfiguration.html
type AWSECSService_AwsVpcConfiguration struct {

	// AssignPublicIp AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-awsvpcconfiguration.html#cfn-ecs-service-awsvpcconfiguration-assignpublicip
	AssignPublicIp *StringIntrinsic `json:"AssignPublicIp,omitempty"`

	// SecurityGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-awsvpcconfiguration.html#cfn-ecs-service-awsvpcconfiguration-securitygroups
	SecurityGroups []*StringIntrinsic `json:"SecurityGroups,omitempty"`

	// Subnets AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-awsvpcconfiguration.html#cfn-ecs-service-awsvpcconfiguration-subnets
	Subnets []*StringIntrinsic `json:"Subnets,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECSService_AwsVpcConfiguration) AWSCloudFormationType() string {
	return "AWS::ECS::Service.AwsVpcConfiguration"
}
