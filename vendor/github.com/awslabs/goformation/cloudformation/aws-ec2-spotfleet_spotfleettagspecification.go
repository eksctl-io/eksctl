package cloudformation

// AWSEC2SpotFleet_SpotFleetTagSpecification AWS CloudFormation Resource (AWS::EC2::SpotFleet.SpotFleetTagSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotfleetrequestconfigdata-launchspecifications-tagspecifications.html
type AWSEC2SpotFleet_SpotFleetTagSpecification struct {

	// ResourceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotfleetrequestconfigdata-launchspecifications-tagspecifications.html#cfn-ec2-spotfleet-spotfleettagspecification-resourcetype
	ResourceType *StringIntrinsic `json:"ResourceType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_SpotFleetTagSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.SpotFleetTagSpecification"
}
